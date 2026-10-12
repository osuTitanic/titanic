package bss

import (
	"fmt"
	"time"

	"github.com/osuTitanic/titanic/internal/constants"
	"github.com/osuTitanic/titanic/internal/schemas"
	"github.com/osuTitanic/titanic/internal/state"
)

const (
	topicBumpCooldown       = 48 * time.Hour
	topicDescriptionDivider = "---------------"
)

// TopicContents contains the client-editable parts of a beatmap forum topic
type TopicContents struct {
	Id      int
	Subject string
	Message string
}

// GetTopic returns the existing beatmap topic contents.
func GetTopic(repositories *state.Repositories, beatmapset *schemas.Beatmapset) (*TopicContents, error) {
	if beatmapset.TopicId == nil {
		return nil, nil
	}
	contents := &TopicContents{}
	contents.Id = *beatmapset.TopicId

	topic, err := repositories.ForumTopics.ById(contents.Id)
	if err != nil {
		return nil, fmt.Errorf("bss: fetch topic %d for set %d: %w", contents.Id, beatmapset.Id, err)
	}
	if topic == nil {
		return nil, nil
	}
	contents.Subject = topic.Title

	firstPost, err := repositories.ForumPosts.FetchInitialByTopic(topic.Id)
	if err != nil {
		return nil, fmt.Errorf("bss: fetch initial post for topic %d: %w", topic.Id, err)
	}
	if firstPost != nil {
		contents.Message = firstPost.Content
	}
	return contents, nil
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
