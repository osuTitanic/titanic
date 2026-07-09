# Schemas

This module contains the database model definitions used by [GORM](https://gorm.io/).
The structs map Go fields to the existing PostgreSQL tables with explicit `gorm` column tags, table names, primary keys, defaults & relationships.

Most services should use the [repositories](../repositories/README.md) exposed by the [state](../state/README.md) system instead of querying schema structs directly.

## Usage

Use schemas when you need to create, update, preload, or inspect model data through a repository.

```go
user := &schemas.User{
	Name:     "peppy",
	SafeName: schemas.ResolveSafeName("peppy"),
	Email:    "pe@ppy.sh",
	Country:  "jp",
}

if err := app.Repositories.Users.Create(user); err != nil {
	return err
}
```

When querying with GORM directly, use the schema type as the model.

```go
var scores []*schemas.Score
err := db.
	Where("user_id = ? AND mode = ?", userId, constants.ModeOsu).
	Find(&scores).
	Error
```

## Relationships

Relationships are declared on the schema structs using explicit `foreignKey` and `references` tags.
Repositories accept optional preload names for queries that need those relationships.

```go
user, err := app.Repositories.Users.ById(userId, "Stats", "Badges", "Groups") // preloads the requested relationships
if err != nil {
	return err
}
```
