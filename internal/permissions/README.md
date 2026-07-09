# Permissions

This module resolves & checks user permissions. Permissions are namespaced strings such as `beatmaps.moderation.owner`, granted to (or rejected from) a user both directly and through the groups they belong to.

## Resolving

`state.NewState(...)` creates the repositories & exposes the resolver through `app.Permissions`.

```go
app, err := state.NewState()
if err != nil {
	return err
}
defer app.Close()

set, err := app.Permissions.Resolve(user.Id)
if err != nil {
	return err
}
```

`Resolve` merges a user's own permissions with the permissions of every group they belong to into a single permission `Set`. Both granted and rejected permissions are collected, with rejected ones taking priority during a check.

Create a resolver from existing repositories when tests or tooling need their own:

```go
resolver := permissions.New(permissionsRepository, groupsRepository)

set, err := resolver.Resolve(user.Id)
if err != nil {
	return err
}
```

## Checking

A `Set` is a snapshot of a user's resolved permission context. Resolve it once, and reuse it for every check on a request.

```go
if !set.Has("beatmaps.moderation.owner") {
	return errors.New("not allowed")
}
```

`Has` takes wildcards (`.*`) into account on both the granted & rejected side:

- an exact match grants the permission
- the global wildcard `*` grants every permission
- a namespace wildcard such as `beatmaps.*` grants any permission below it

Checks are case-insensitive, and a trailing `.*` on the queried permission is ignored, so `set.Has("beatmaps.*")` is equivalent to `set.Has("beatmaps")`.

A `nil` `*Set` denies everything, which makes it a good default for guests:

```go
var set *permissions.Set
set.Has("beatmaps.upload") // false
```

## Groups

A `Set` also carries the user's group memberships. Use `InGroup` for an arbitrary check, or the helpers for the staff groups:

```go
set.InGroup(constants.GroupBAT, constants.GroupGMT)

set.IsAdmin()     // Admin or Developer
set.IsBat()       // Admin, Developer or Beatmap Approval Team
set.IsModerator() // Admin, Developer or Global Moderation Team
```

Group IDs are defined in `constants` (e.g. `constants.GroupBAT`).

## Usage in stern

The website resolves permissions through the request context. Resolution happens once per request & is memoized.

```go
if ctx.HasPermission("beatmaps.moderation.owner") {
	// ...
}
```

`ctx.Permissions()` returns the resolved `Set` directly, e.g. to pass into a template view for repeated checks. For guests, it yields an empty set that denies everything.
