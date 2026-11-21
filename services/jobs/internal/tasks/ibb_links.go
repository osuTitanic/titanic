package tasks

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/osuTitanic/titanic/internal/schemas"
	"github.com/osuTitanic/titanic/internal/state"
)

var ibbLinkRegex = regexp.MustCompile(`https://i\.ibb\.co/([A-Za-z0-9]+)/([A-Za-z0-9._-]+)`)

func MoveIbbLinksToTitanic(app *state.State, logger *slog.Logger) error {
	var posts []*schemas.ForumPost
	err := app.Database.Where("content LIKE ?", "%ibb.co%").
		Find(&posts).Error
	if err != nil {
		return err
	}

	logger.Info("Moving ibb.co links to titanic", "count", len(posts))

	for _, post := range posts {
		err := processPost(post, app)
		if err != nil {
			logger.Error("Failed to process post", "post", post.Id, "error", err)
			continue
		}
	}
	return nil
}

func processPost(post *schemas.ForumPost, app *state.State) error {
	matchingLinks := ibbLinkRegex.FindAllStringSubmatch(post.Content, -1)
	if len(matchingLinks) == 0 {
		return nil
	}

	for _, match := range matchingLinks {
		if len(match) < 3 {
			continue
		}
		imageUrl := match[0]
		app.Logger.Debug("Moving ibb.co link", "post", post.Id, "link", imageUrl)

		data, err := downloadImage(imageUrl)
		if err != nil {
			return err
		}
		if data == nil {
			return fmt.Errorf("failed to download image from %s", imageUrl)
		}
		defer data.Close()

		screenshot := &schemas.Screenshot{UserId: post.UserId}
		if err := app.Screenshots.Create(screenshot); err != nil {
			return fmt.Errorf("create screenshot: %v", err)
		}

		key := strconv.Itoa(screenshot.Id)
		err = app.Storage.SaveStream(context.Background(), key, "screenshots", data)
		if err != nil {
			return fmt.Errorf("store screenshot: %v", err)
		}

		replacementPath := fmt.Sprintf("/ss/%d/%s", screenshot.Id, screenshot.Checksum())
		replacementLink := app.Config.OsuBaseUrl() + replacementPath

		post.Content = strings.ReplaceAll(post.Content, imageUrl, replacementLink)
		app.ForumPosts.Update(post, "content")

		app.Logger.Info(
			"Replaced ibb.co link",
			"post", post.Id, "original", imageUrl, "new", replacementLink,
		)
	}
	return nil
}

func downloadImage(url string) (io.ReadCloser, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return nil, err
	}

	return resp.Body, nil
}
