package authentication

import (
	"slices"

	"github.com/osuTitanic/titanic-go/internal/schemas"
)

const AuthenticatedScope = "users.authenticated"

func ResolveScopes(granted []string, rejected []string) []string {
	scopeSet := map[string]struct{}{
		AuthenticatedScope: {},
	}

	for _, scope := range granted {
		if scope == "" {
			continue
		}
		scopeSet[scope] = struct{}{}
	}

	for _, scope := range rejected {
		delete(scopeSet, scope)
	}

	scopes := make([]string, 0, len(scopeSet))
	for _, scope := range append([]string{AuthenticatedScope}, granted...) {
		if _, exists := scopeSet[scope]; !exists {
			continue
		}

		duplicate := slices.Contains(scopes, scope)
		if !duplicate {
			scopes = append(scopes, scope)
		}
	}

	if len(scopes) == 0 {
		return []string{AuthenticatedScope}
	}

	return scopes
}

func ResolveUserScopes(userPermissions []*schemas.UserPermission, groupPermissions []*schemas.GroupPermission) []string {
	granted := make([]string, 0, len(userPermissions)+len(groupPermissions))
	rejected := make([]string, 0)

	for _, permission := range groupPermissions {
		if permission == nil || permission.Permission == "" {
			continue
		}
		if permission.Rejected {
			rejected = append(rejected, permission.Permission)
			continue
		}
		granted = append(granted, permission.Permission)
	}

	for _, permission := range userPermissions {
		if permission == nil || permission.Permission == "" {
			continue
		}
		if permission.Rejected {
			rejected = append(rejected, permission.Permission)
			continue
		}
		granted = append(granted, permission.Permission)
	}

	return ResolveScopes(granted, rejected)
}
