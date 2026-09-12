package clients

import (
	"strconv"

	"github.com/osuTitanic/titanic/internal/constants"
)

type Version struct {
	Date     int
	Name     string
	Revision *int
	Stream   string
}

func (version Version) Identifier() string {
	if version.Stream != "" {
		return version.Stream
	}
	if version.Name != "" {
		return version.Name
	}
	return "stable"
}

func ParseVersion(value string) (Version, bool) {
	if len(value) > 25 {
		return Version{}, false
	}

	match := constants.OsuVersion.FindStringSubmatch(value)
	if match == nil {
		return Version{}, false
	}

	date, err := strconv.Atoi(match[constants.OsuVersion.SubexpIndex("date")])
	if err != nil {
		return Version{}, false
	}

	version := Version{
		Date:   date,
		Name:   match[constants.OsuVersion.SubexpIndex("name")],
		Stream: match[constants.OsuVersion.SubexpIndex("stream")],
	}
	if rawRevision := match[constants.OsuVersion.SubexpIndex("revision")]; rawRevision != "" {
		revision, err := strconv.Atoi(rawRevision)
		if err != nil {
			return Version{}, false
		}
		version.Revision = new(revision)
	}
	return version, true
}
