package bss

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/osuTitanic/titanic/internal/constants"
	"github.com/osuTitanic/titanic/internal/schemas"
	"github.com/osuTitanic/titanic/internal/state"
)

const (
	topicBumpCooldown       = 48 * time.Hour
	topicDescriptionDivider = "---------------"
)

// TopicContents contains the client-editable parts of a beatmap forum topic.
type TopicContents struct {
	Id      int
	Subject string
	Message string
}

// TopicPostUpdate describes a beatmap topic update.
type TopicPostUpdate struct {
	Subject   string
	Message   string
	Status    constants.BeatmapStatus
	Subscribe bool
	Bump      bool
}

func (post TopicPostUpdate) DestinationForumId() int {
	switch post.Status {
	case constants.BeatmapStatusPending:
		return constants.ForumBeatmapsPending
	default:
		return constants.ForumBeatmapsWIP
	}
}

func (post TopicPostUpdate) StatusText() string {
	switch post.Status {
	case constants.BeatmapStatusPending:
		return "Waiting for approval..."
	default:
		return "Needs modding"
	}
}

func (post TopicPostUpdate) Validate() error {
	switch post.Status {
	case constants.BeatmapStatusWIP, constants.BeatmapStatusPending:
		return nil
	default:
		return fmt.Errorf("%w: unsupported topic status %d", ErrInvalidSubmissionRequest, post.Status)
	}
}

// GetTopic returns the existing beatmap topic contents.
func GetTopic(repositories *state.Repositories, beatmapset *schemas.Beatmapset) (*TopicContents, error) {
	contents := &TopicContents{}
	if beatmapset.TopicId == nil {
		return contents, nil
	}
	if beatmapset.PostId == nil {
		return contents, nil
	}
	topicId := *beatmapset.TopicId

	topic, err := repositories.ForumTopics.ById(topicId)
	if err != nil {
		return nil, fmt.Errorf("bss: fetch topic %d for set %d: %w", topicId, beatmapset.Id, err)
	}
	if topic == nil {
		return contents, nil
	}

	firstPost, err := repositories.ForumPosts.ById(*beatmapset.PostId)
	if err != nil {
		return nil, fmt.Errorf("bss: fetch initial post for topic %d: %w", topic.Id, err)
	}
	if firstPost == nil {
		return contents, nil
	}

	contents.Id = topic.Id
	contents.Subject = topic.Title
	contents.Message = firstPost.Content
	return contents, nil
}

// PostTopic creates or updates a beatmap topic and applies its state to the set.
func PostTopic(transaction *state.Repositories, beatmapset *schemas.Beatmapset, post TopicPostUpdate) (topic *schemas.ForumTopic, err error) {
	if err := post.Validate(); err != nil {
		return nil, err
	}
	if beatmapset.CreatorId == nil {
		return nil, fmt.Errorf("bss: beatmapset creator id is nil")
	}

	if beatmapset.TopicId != nil {
		topic, err = transaction.ForumTopics.ById(*beatmapset.TopicId)
		if err != nil {
			return nil, err
		}
	}

	if topic == nil {
		// Create a new topic if the set has no topic or its reference is dangling.
		topic, err = createBeatmapTopic(transaction, beatmapset, post)
		if err != nil {
			return nil, err
		}
	} else if err := updateBeatmapTopic(transaction, beatmapset, topic, post); err != nil {
		return nil, err
	}

	// Update the beatmap(set) status, description, topic ID, ...
	if err := applyTopicStateToBeatmapset(transaction, beatmapset, topic.Id, post); err != nil {
		return nil, err
	}
	// Update user subscription to the topic
	if err := setTopicSubscription(transaction, topic.Id, topic.CreatorId, post.Subscribe); err != nil {
		return nil, err
	}
	return topic, nil
}

// TODO: Find a better place for PopBubble maybe...

func PopBubble(transaction *state.Repositories, beatmapset *schemas.Beatmapset) error {
	hadNominations, err := transaction.Nominations.ExistsForSet(beatmapset.Id)
	if err != nil {
		return fmt.Errorf("bss: check nominations: %w", err)
	}
	if err := transaction.Nominations.DeleteAll(beatmapset.Id); err != nil {
		return fmt.Errorf("bss: reset nominations: %w", err)
	}
	if hadNominations {
		beatmapset.StarPriority += 5
		transaction.Beatmapsets.Update(beatmapset, "star_priority")
	}
	if beatmapset.TopicId == nil {
		return nil
	}

	topic, err := transaction.ForumTopics.ById(*beatmapset.TopicId)
	if err != nil {
		return fmt.Errorf("bss: fetch beatmap topic: %w", err)
	}
	if topic == nil {
		// hm weird, surely nothing has gone wrong here...
		return nil
	}

	hasBubbleIcon := topic.IconId != nil && *topic.IconId == constants.ForumIconBubble
	shouldUpdate := hadNominations || hasBubbleIcon

	if !shouldUpdate {
		return nil
	}

	topic.IconId = new(constants.ForumIconBubblePop)
	_, err = transaction.ForumTopics.Update(topic, "icon")
	if err != nil {
		return fmt.Errorf("bss: update topic icon: %w", err)
	}

	return nil
}

func createBeatmapTopic(transaction *state.Repositories, beatmapset *schemas.Beatmapset, post TopicPostUpdate) (*schemas.ForumTopic, error) {
	statusText := post.StatusText()
	now := time.Now()

	topic := &schemas.ForumTopic{
		ForumId:       post.DestinationForumId(),
		CreatorId:     *beatmapset.CreatorId,
		Title:         post.Subject,
		StatusText:    &statusText,
		CreatedAt:     now,
		LastPostAt:    now,
		CanChangeIcon: true,
	}
	if err := transaction.ForumTopics.Create(topic); err != nil {
		return nil, fmt.Errorf("bss: create topic for set %d: %w", beatmapset.Id, err)
	}

	firstPost := &schemas.ForumPost{
		TopicId:    topic.Id,
		ForumId:    topic.ForumId,
		UserId:     *beatmapset.CreatorId,
		Content:    post.Message,
		CreatedAt:  now,
		EditTime:   now,
		EditLocked: true,
	}
	if err := transaction.ForumPosts.Create(firstPost); err != nil {
		return nil, fmt.Errorf("bss: create initial post for topic %d: %w", topic.Id, err)
	}
	return topic, nil
}

func updateBeatmapTopic(
	transaction *state.Repositories,
	beatmapset *schemas.Beatmapset,
	topic *schemas.ForumTopic,
	post TopicPostUpdate,
) (err error) {
	if topic == nil {
		return errors.New("bss: topic is nil")
	}
	var firstPost *schemas.ForumPost

	if beatmapset.PostId != nil {
		firstPost, err = transaction.ForumPosts.ById(*beatmapset.PostId)
		if err != nil {
			return fmt.Errorf("bss: fetch initial post for topic %d: %w", topic.Id, err)
		}
	}

	now := time.Now()
	topic.ForumId = post.DestinationForumId()
	topic.StatusText = new(post.StatusText())
	topic.Title = post.Subject
	if _, err := transaction.ForumTopics.Update(topic, "forum_id", "title", "status_text"); err != nil {
		return fmt.Errorf("bss: update topic %d: %w", topic.Id, err)
	}

	postMove := &schemas.ForumPost{
		TopicId: topic.Id,
		ForumId: topic.ForumId,
	}
	if _, err := transaction.ForumPosts.UpdateByTopic(postMove, "forum_id"); err != nil {
		return fmt.Errorf("bss: move posts for topic %d: %w", topic.Id, err)
	}

	if firstPost == nil {
		firstPost = &schemas.ForumPost{
			TopicId:    topic.Id,
			ForumId:    topic.ForumId,
			UserId:     *beatmapset.CreatorId,
			Content:    post.Message,
			CreatedAt:  now,
			EditTime:   now,
			EditLocked: true,
		}
		if err := transaction.ForumPosts.Create(firstPost); err != nil {
			return fmt.Errorf("bss: recreate initial post for topic %d: %w", topic.Id, err)
		}

		topic.LastPostAt = now
		if _, err := transaction.ForumTopics.Update(topic, "last_post_at"); err != nil {
			return fmt.Errorf("bss: update activity for topic %d: %w", topic.Id, err)
		}
		return nil
	}

	firstPost.ForumId = topic.ForumId
	firstPost.Content = post.Message
	firstPost.Deleted = false
	firstPost.EditTime = now
	firstPost.EditCount++
	if _, err := transaction.ForumPosts.Update(
		firstPost,
		"content",
		"deleted",
		"edit_time",
		"edit_count",
	); err != nil {
		return fmt.Errorf("bss: update initial post %d for topic %d: %w", firstPost.Id, topic.Id, err)
	}

	// Bumping a topic has a cooldown of 48 hours
	bumpAllowed := firstPost.CreatedAt.Before(now.Add(-topicBumpCooldown))

	if post.Bump && bumpAllowed {
		topic.LastPostAt = now
		if _, err := transaction.ForumTopics.Update(topic, "last_post_at"); err != nil {
			return fmt.Errorf("bss: bump topic %d: %w", topic.Id, err)
		}
	}
	return nil
}

func applyTopicStateToBeatmapset(
	transaction *state.Repositories,
	beatmapset *schemas.Beatmapset,
	topicId int,
	post TopicPostUpdate,
) error {
	description := topicDescription(post.Message)

	beatmapset.TopicId = &topicId
	beatmapset.Status = post.Status
	beatmapset.Description = &description
	beatmapset.Available = true
	beatmapset.LastUpdate = time.Now()

	if _, err := transaction.Beatmapsets.Update(
		beatmapset,
		"topic_id",
		"submission_status",
		"description",
		"available",
		"last_updated",
	); err != nil {
		return fmt.Errorf("bss: update publication state for set %d: %w", beatmapset.Id, err)
	}

	// Postgres will automatically update the PostId for us
	// Lets reload the beatmapset and update our local copy
	updated, err := transaction.Beatmapsets.ById(beatmapset.Id)
	if err != nil {
		return fmt.Errorf("bss: reload publication state for set %d: %w", beatmapset.Id, err)
	}
	if updated == nil {
		return fmt.Errorf("bss: beatmapset %d disappeared during publication", beatmapset.Id)
	}
	beatmapset.PostId = updated.PostId

	beatmapUpdate := &schemas.Beatmap{
		SetId:  beatmapset.Id,
		Status: post.Status,
	}
	if _, err := transaction.Beatmaps.UpdateBySetId(beatmapUpdate, "status"); err != nil {
		return fmt.Errorf("bss: update difficulty status for set %d: %w", beatmapset.Id, err)
	}

	for _, beatmap := range beatmapset.Beatmaps {
		if beatmap != nil {
			beatmap.Status = post.Status
		}
	}
	return nil
}

func setTopicSubscription(transaction *state.Repositories, topicId, userId int, subscribe bool) error {
	subscriber := &schemas.ForumSubscriber{
		TopicId: topicId,
		UserId:  userId,
	}
	if !subscribe {
		if err := transaction.ForumSubscribers.Delete(subscriber); err != nil {
			return fmt.Errorf("bss: remove subscription for user %d from topic %d: %w", userId, topicId, err)
		}
		return nil
	}

	exists, err := transaction.ForumSubscribers.Exists(topicId, userId)
	if err != nil {
		return fmt.Errorf("bss: check subscription for user %d to topic %d: %w", userId, topicId, err)
	}
	if exists {
		return nil
	}
	if err := transaction.ForumSubscribers.Create(subscriber); err != nil {
		return fmt.Errorf("bss: add subscription for user %d to topic %d: %w", userId, topicId, err)
	}
	return nil
}

func topicDescription(message string) string {
	normalized := strings.ReplaceAll(message, "\r\n", "\n")
	lines := strings.Split(normalized, "\n")

	for i, line := range lines {
		if line == topicDescriptionDivider {
			return strings.Join(lines[i+1:], "\n")
		}
	}
	return message
}
