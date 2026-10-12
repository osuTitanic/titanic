package bss

import (
	"fmt"

	"github.com/osuTitanic/titanic/internal/constants"
	"github.com/osuTitanic/titanic/internal/schemas"
	"github.com/osuTitanic/titanic/internal/state"
)

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

// TODO: Forum post logic
