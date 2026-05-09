# Repositories

This module contains repositories for the database schemas.
Repositories give services a common place to create, update, delete & fetch records without repeating GORM calls everywhere.

For the database model definitions, see [Schemas](../schemas/README.md).

## Usage

Services should use repositories through `state.NewState(...)`.
`State` creates the database session, wires every repository once & exposes them through `app.Repositories`.

```go
app, err := state.NewState()
if err != nil {
	return err
}
defer app.Close()

user, err := app.Repositories.Users.ById(userId, "Stats", "Badges")
if err != nil {
	return err
}
```

## Common Patterns

Most repositories follow a similar pattern of methods:

- `Create(model)` inserts a schema record
- `Delete(model)` deletes a schema record
- `Update(model, columns...)` saves a full model or updates only the selected columns
- `By...` methods fetch a single record
- `Many...` and `Fetch...` methods fetch collections or domain-specific results
- `Count...` methods return aggregate counts

Optional preload arguments are passed through to GORM `Preload`.

```go
beatmapset, err := app.Repositories.Beatmapsets.ById(setId, "Beatmaps", "CreatorUser")
if err != nil {
	return err
}
```

Column-limited updates are useful when only part of a record should be persisted.

```go
updates := &schemas.User{
	Id:             userId,
	LatestActivity: time.Now(),
}

rows, err := app.Repositories.Users.Update(updates, "latest_activity") // only "latest_activity" will be updated
if err != nil {
	return err
}
```
