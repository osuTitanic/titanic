package permissions

import (
	"github.com/osuTitanic/titanic-go/internal/repositories"
)

// Resolver resolves the authorization context for a user.
type Resolver interface {
	Resolve(userId int) (*Set, error)
}

type dbResolver struct {
	permissions *repositories.PermissionsRepository
	groups      *repositories.GroupRepository
}

// TODO: Add a caching layer to resolver

// New creates a database-backed permission resolver.
func New(permissions *repositories.PermissionsRepository, groups *repositories.GroupRepository) Resolver {
	return &dbResolver{permissions: permissions, groups: groups}
}

// Resolve fetches the user's own permissions and the permissions of every group
// they belong to, merging them into a single Set. This goes for both granted and
// rejected permissions, with rejected ones taking priority.
func (r *dbResolver) Resolve(userId int) (*Set, error) {
	userPermissions, err := r.permissions.ManyUserPermissionsByUserId(userId)
	if err != nil {
		return nil, err
	}

	entries, err := r.groups.ManyEntriesByUserId(userId, "Group", "Group.Permissions")
	if err != nil {
		return nil, err
	}

	set := &Set{
		GroupIds: make([]int, 0, len(entries)),
		Granted:  make([]string, 0),
		Rejected: make([]string, 0),
	}

	for _, entry := range entries {
		set.GroupIds = append(set.GroupIds, entry.GroupId)
		if entry.Group == nil {
			continue
		}
		for _, permission := range entry.Group.Permissions {
			set.add(permission.Permission, permission.Rejected)
		}
	}

	for _, permission := range userPermissions {
		set.add(permission.Permission, permission.Rejected)
	}

	return set, nil
}
