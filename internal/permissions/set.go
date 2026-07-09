package permissions

import (
	"slices"
	"strings"

	"github.com/osuTitanic/titanic/internal/constants"
)

// Set is a snapshot of a user's resolved permission context, i.e.
// their group memberships and the permissions granted to & rejected
// from them (via both user- and group-level permissions).
type Set struct {
	GroupIds []int
	Granted  []string
	Rejected []string
}

// add appends a permission string to the granted / rejected list
func (s *Set) add(permission string, rejected bool) {
	if permission == "" {
		return
	}

	permission = strings.ToLower(permission)
	if rejected {
		s.Rejected = append(s.Rejected, permission)
	} else {
		s.Granted = append(s.Granted, permission)
	}
}

// Has reports whether the user is allowed to perform the given permission.
func (s *Set) Has(permission string) bool {
	if s == nil {
		return false
	}

	permission = normalize(permission)
	return includes(permission, s.Granted) && !includes(permission, s.Rejected)
}

// InGroup reports whether the user is a member of any of the given groups.
func (s *Set) InGroup(ids ...int) bool {
	if s == nil {
		return false
	}

	for _, id := range ids {
		if slices.Contains(s.GroupIds, id) {
			return true
		}
	}
	return false
}

func (s *Set) IsAdmin() bool {
	return s.InGroup(constants.GroupAdmin, constants.GroupDeveloper)
}

func (s *Set) IsBat() bool {
	return s.InGroup(constants.GroupAdmin, constants.GroupDeveloper, constants.GroupBAT)
}

func (s *Set) IsModerator() bool {
	return s.InGroup(constants.GroupAdmin, constants.GroupDeveloper, constants.GroupGMT)
}
