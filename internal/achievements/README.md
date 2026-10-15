# Achievements

This module contains the definitions for all achievements, which are used to unlock them in score submission.

## Evaluating achievements

Use `LoadContext` to fetch the data required by achievement checks, then call `Evaluate` to get the newly matched achievements.

```go
context, err := achievements.LoadContext(
	app.Repositories,
	score,
	currentStats,
)
if err != nil {
	return err
}

matched := context.Evaluate()
```

This does not automatically unlock the achievements for the user. That has to be done manually, e.g. like [score submission](../../services/deck/internal/scoring/achievements.go) does.

## Definitions

`Definitions` contains every achievement checked by `Evaluate`. Each `AchievementDefinition` has a display name, category, filename & check function.
