package tasks

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"slices"
	"strings"

	"github.com/osuTitanic/titanic/internal/discord"
	"github.com/osuTitanic/titanic/internal/schemas"
	"github.com/osuTitanic/titanic/internal/state"
)

const updateUrl = "https://osu.ppy.sh/web/check-updates.php"
const windowsOSParameter = "10.0.0.26100.1.0"

var releaseStreams = []string{"stable40", "cuttingedge", "beta40"}

// ReleaseUpdates checks each release stream (stable, cuttingedge, ...) for updates
// and notifies on updates through a webhook.
func ReleaseUpdates(app *state.State, logger *slog.Logger) error {
	if !app.Config.ReleaseUpdatesEnabled {
		return errors.New("release updates are disabled in config")
	}

	for _, stream := range releaseStreams {
		updates, err := checkStream(app, logger, stream)
		if err != nil {
			logger.Error("Failed to check release stream", "stream", stream, "error", err)
			continue
		}

		for _, file := range updates {
			logger.Info(
				"New release file created",
				"stream", stream, "file_version", file.FileVersion, "filename", file.Filename,
			)
			postUpdateActions(app, logger, file, stream)
		}
	}
	return nil
}

func checkStream(app *state.State, logger *slog.Logger, stream string) ([]*schemas.ReleaseFiles, error) {
	logger.Info("Checking for updates on release stream", "stream", stream)

	resultsWindows, err := fetchStream(stream, windowsOSParameter)
	if err != nil {
		return nil, err
	}
	resultsLinux, err := fetchStream(stream, "")
	if err != nil {
		return nil, err
	}
	results := slices.Concat(resultsWindows, resultsLinux)
	resultsNew := make([]*schemas.ReleaseFiles, 0)

	for _, file := range results {
		existingFile, err := app.Repositories.ReleasesOfficial.FetchFileByVersion(file.FileVersion)
		if err != nil {
			return nil, err
		}

		if existingFile != nil {
			logger.Debug(
				"Release file already exists, skipping...",
				"stream", stream, "file_id", file.Id, "filename", file.Filename,
			)
			continue
		}

		if err := app.Repositories.ReleasesOfficial.CreateFile(file); err != nil {
			return nil, err
		}
		resultsNew = append(resultsNew, file)
	}
	return resultsNew, nil
}

func fetchStream(stream string, os string) ([]*schemas.ReleaseFiles, error) {
	params := url.Values{}
	params.Set("action", "check")
	params.Set("stream", stream)

	if os != "" {
		// Different versions of osu!auth.dll exist for linux & windows
		// The "os" query parameter determines which variant to retrieve
		params.Set("os", os)
	}

	targetUrl := fmt.Sprintf("%s?%s", updateUrl, params.Encode())
	response, err := http.Get(targetUrl)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	var files []*schemas.ReleaseFiles
	if err := json.NewDecoder(response.Body).Decode(&files); err != nil {
		return nil, err
	}

	return files, nil
}

func postUpdateActions(app *state.State, logger *slog.Logger, file *schemas.ReleaseFiles, stream string) error {
	notifyWebhook(app, logger, file, stream)
	uploadFile(file.UrlFull, app, logger)
	if file.UrlPatch != nil {
		uploadFile(*file.UrlPatch, app, logger)
	}
	return nil
}

func notifyWebhook(app *state.State, logger *slog.Logger, file *schemas.ReleaseFiles, stream string) error {
	if app.Config.ReleaseUpdateNotifyWebhook == "" {
		return nil
	}

	webhook := discord.Webhook{
		URL:       app.Config.ReleaseUpdateNotifyWebhook,
		Username:  ptr("Release Updates"),
		AvatarURL: ptr("https://osu.ppy.sh/images/layout/osu-logo.png"),
	}
	embed := discord.Embed{
		Title:       ptr("New release file"),
		Description: ptr(fmt.Sprintf("A new file has been added to the `%s` release stream.", stream)),
		Color:       ptr(0xFF66AB),
	}
	embed.AddField("Filename", fmt.Sprintf("`%s`", file.Filename), true)
	embed.AddField("Version", fmt.Sprintf("`%d`", file.FileVersion), true)
	embed.AddField("File Hash", fmt.Sprintf("`%s`", file.FileHash), true)
	embed.AddField("Filesize", fmt.Sprintf("`%d bytes`", file.Filesize), true)

	if file.UrlFull != "" {
		embed.AddField("Download", fmt.Sprintf("[Full](%s)", file.UrlFull), true)
	}
	if file.UrlPatch != nil {
		embed.AddField("Download", fmt.Sprintf("[Patch](%s)", *file.UrlPatch), true)
	}
	webhook.AddEmbed(embed)

	if err := webhook.Post(); err != nil {
		logger.Error("Failed to post webhook notification for new release file", "filename", file.Filename, "error", err)
	}
	return nil
}

func uploadFile(url string, app *state.State, logger *slog.Logger) error {
	if !app.Config.ReleaseUpdatesEnabled {
		return nil
	}
	if url == "" {
		return nil
	}
	logger.Info("Uploading release file to storage", "url", url)

	urlParts := strings.Split(url, "/")
	checksum := urlParts[len(urlParts)-1]
	filename := urlParts[len(urlParts)-2]
	location := fmt.Sprintf("%s/%s", app.Config.ReleaseUpdateLocation, filename)
	app.Storage.SaveUrl(checksum, location, url)
	return nil
}

func ptr[T any](v T) *T {
	return &v
}
