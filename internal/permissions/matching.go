package permissions

import "strings"

// normalize lowercases a permission and removes a trailing ".*" wildcard, so
// that querying "beatmaps.*" is treated the same as querying "beatmaps".
func normalize(permission string) string {
	return strings.TrimSuffix(strings.ToLower(permission), ".*")
}

// includes reports whether the given permission is covered by any
// entry in the list. Criterias are as follows:
//
//   - an exact match grants the permission
//   - the global wildcard "*" grants every permission
//   - a namespace wildcard such as "beatmaps.*" grants any permission below it
func includes(permission string, entries []string) bool {
	for _, entry := range entries {
		if entry == permission || entry == "*" {
			return true
		}

		if !strings.HasSuffix(entry, ".*") {
			continue
		}

		if strings.HasPrefix(permission, entry[:len(entry)-2]) {
			return true
		}
	}
	return false
}
