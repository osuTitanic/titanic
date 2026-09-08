package routes

import (
	"fmt"
	"io"
	"io/fs"
	"mime"
	"net/http"
	"net/url"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/osuTitanic/titanic/internal/schemas"
	"github.com/osuTitanic/titanic/services/deck/internal/server"
)

// TODO: Document in detail how the old update system worked
//		 Also actually implement these update endpoints
//		 I am thinking of making the version rotate each day
//		 so that its not just a static version that never changes

var releaseMirrorUrl = "https://m1.ppy.sh/release/"
var releaseHttpClient = &http.Client{Timeout: 10 * time.Second}

// /release/update -> Text update manifest
func ReleaseUpdate(ctx *server.Context) {
	ctx.RenderText(http.StatusOK, "\n")
}

// /release/update.php -> Executable update check
func ReleaseUpdateCheck(ctx *server.Context) {
	ctx.RenderText(http.StatusOK, "0")
}

// /release/update2.php -> pUpdater / osume manifest
func ReleaseUpdateV2(ctx *server.Context) {
	ctx.RenderText(http.StatusOK, "\n")
}

// /release/patches.php -> Binary patch list
func ReleasePatches(ctx *server.Context) {
	ctx.RenderText(http.StatusOK, "")
}

// /release/filter.txt -> Chat filter list, mirrored from ppy.sh
func ReleaseFilter(ctx *server.Context) {
	proxyReleaseAsset(ctx, releaseMirrorUrl+"filter.txt")
}

// /release/Localisation/{filename} -> Text localisation files, mirrored from ppy.sh
func ReleaseLocalisation(ctx *server.Context) {
	filename := ctx.PathValue("filename")
	if !validReleasePath(filename) {
		ctx.Response.WriteHeader(http.StatusNotFound)
		return
	}

	target, err := url.Parse(releaseMirrorUrl + "Localisation/" + url.PathEscape(filename))
	if err != nil {
		ctx.Response.WriteHeader(http.StatusInternalServerError)
		return
	}
	target.RawQuery = ctx.Request.URL.RawQuery
	proxyReleaseAsset(ctx, target.String())
}

// /release/{language}/{filename} -> Localisation dll files, the predecessor to the above endpoint
func ReleaseLegacyLocalisation(ctx *server.Context) {
	language := ctx.PathValue("language")
	filename := ctx.PathValue("filename")

	if !validReleasePath(language) || !validReleasePath(filename) {
		ctx.Response.WriteHeader(http.StatusNotFound)
		return
	}
	if !strings.HasSuffix(filename, ".dll") {
		ctx.Response.WriteHeader(http.StatusNotFound)
		return
	}

	if !serveLocalReleaseFile(ctx, language+"/"+filename, "") {
		ctx.Response.WriteHeader(http.StatusNotFound)
	}
}

// /release/{filename} -> Download a release file
func ReleaseFile(ctx *server.Context) {
	filename := ctx.PathValue("filename")
	if !validReleasePath(filename) {
		ctx.Response.WriteHeader(http.StatusNotFound)
		return
	}

	checksum := ctx.QueryValueOptional("v")
	if checksum != nil {
		// If a checksum is provided, we can attempt to fetch the correct file from our database
		file, err := ctx.State.ReleasesOfficial.FetchFileByChecksum(*checksum)
		if err != nil {
			ctx.Logger.Error("Failed to find release file by checksum", "checksum", *checksum, "error", err)
			ctx.Response.WriteHeader(http.StatusInternalServerError)
			return
		}
		if file != nil && serveRemoteReleaseFile(ctx, file, file.UrlFull, file.Filename) {
			return
		}
	}

	// If no checksum is provided, we can attempt to fetch the file by its patch filename
	// This would also return a consistent file, since the filename includes a checksum of the 2 executables
	file, err := ctx.State.ReleasesOfficial.FetchFileByPatchFilename(filename)
	if err != nil {
		ctx.Logger.Error("Failed to find release patch", "filename", filename, "error", err)
		ctx.Response.WriteHeader(http.StatusInternalServerError)
		return
	}
	if file != nil && file.UrlPatch != nil && serveRemoteReleaseFile(ctx, file, *file.UrlPatch, filename) {
		return
	}

	// We tried our best... maybe we still have the file available locally
	if !serveLocalReleaseFile(ctx, filename, filename) {
		ctx.Response.WriteHeader(http.StatusNotFound)
	}
}

func proxyReleaseAsset(ctx *server.Context, target string) {
	request, err := newReleaseRequest(ctx, target)
	if err != nil {
		ctx.Logger.Error("Failed to create release proxy request", "url", target, "error", err)
		ctx.Response.WriteHeader(http.StatusBadGateway)
		return
	}

	response, err := releaseHttpClient.Do(request)
	if err != nil {
		ctx.Logger.Warn("Failed to fetch release asset", "url", target, "error", err)
		ctx.Response.WriteHeader(http.StatusBadGateway)
		return
	}
	defer response.Body.Close()

	if response.StatusCode < http.StatusOK || response.StatusCode >= http.StatusBadRequest {
		ctx.Response.WriteHeader(response.StatusCode)
		return
	}

	if contentType := response.Header.Get("Content-Type"); contentType != "" {
		ctx.Response.Header().Set("Content-Type", contentType)
	}
	if response.ContentLength >= 0 {
		ctx.Response.Header().Set("Content-Length", strconv.FormatInt(response.ContentLength, 10))
	}
	ctx.Response.WriteHeader(response.StatusCode)

	if _, err := io.Copy(ctx.Response, response.Body); err != nil {
		ctx.Logger.Warn("Failed to stream release asset", "url", target, "error", err)
	}
}

func serveLocalReleaseFile(ctx *server.Context, key string, downloadFilename string) bool {
	stream, size, err := ctx.State.Storage.ReadStreamAt(key, "release")
	if err != nil {
		return false
	}
	defer stream.Close()

	if downloadFilename != "" {
		ctx.Response.Header().Set(
			"Content-Disposition",
			mime.FormatMediaType("attachment", map[string]string{"filename": downloadFilename}),
		)
	}

	ctx.Response.Header().Set(
		"Content-Type",
		"application/octet-stream",
	)
	http.ServeContent(
		ctx.Response,
		ctx.Request,
		downloadFilename,
		time.Time{},
		io.NewSectionReader(stream, 0, size),
	)
	return true
}

func serveRemoteReleaseFile(ctx *server.Context, file *schemas.ReleaseFiles, target string, downloadFilename string) bool {
	request, err := newReleaseRequest(ctx, target)
	if err != nil {
		ctx.Logger.Warn("Failed to create release download request", "url", target, "error", err)
		return false
	}

	response, err := releaseHttpClient.Do(request)
	if err != nil {
		ctx.Logger.Warn("Failed to fetch release file", "url", target, "error", err)
		return false
	}
	defer response.Body.Close()

	if response.StatusCode < http.StatusOK || response.StatusCode >= http.StatusBadRequest {
		ctx.Logger.Warn("Release file upstream returned an error", "url", target, "status", response.StatusCode)
		return false
	}

	ctx.Response.Header().Set(
		"Content-Type",
		"application/octet-stream",
	)
	ctx.Response.Header().Set(
		"Content-Disposition",
		mime.FormatMediaType("attachment", map[string]string{"filename": downloadFilename}),
	)
	if response.ContentLength >= 0 {
		ctx.Response.Header().Set("Content-Length", strconv.FormatInt(response.ContentLength, 10))
	}
	if !file.Timestamp.IsZero() {
		ctx.Response.Header().Set("Last-Modified", file.Timestamp.UTC().Format(http.TimeFormat))
	}
	ctx.Response.WriteHeader(response.StatusCode)

	if _, err := io.Copy(ctx.Response, response.Body); err != nil {
		ctx.Logger.Warn("Failed to stream release file", "url", target, "error", err)
	}
	return true
}

func newReleaseRequest(ctx *server.Context, target string) (*http.Request, error) {
	request, err := http.NewRequestWithContext(
		ctx.Request.Context(),
		http.MethodGet, target, nil,
	)
	if err != nil {
		return nil, err
	}
	request.Header.Set(
		"User-Agent",
		fmt.Sprintf("osuTitanic (%s)", ctx.State.Config.DomainName),
	)
	return request, nil
}

func validReleasePath(value string) bool {
	return value != "." &&
		// Prevent path traversal
		value != ".." &&
		// Prevent absolute paths
		!strings.ContainsAny(value, `/\\`) &&
		// Prevent invalid UTF-8
		fs.ValidPath(value) &&
		// Prevent non-local paths, e.g. symlinks
		filepath.IsLocal(filepath.FromSlash(value))
}
