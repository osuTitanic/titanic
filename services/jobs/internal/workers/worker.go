package workers

import (
	"context"
	"sync"

	"github.com/osuTitanic/titanic/internal/state"
)

// RunWorkerPool processes items with up to a specified amount of concurrent workers
func RunWorkerPool[T any](items []T, workerCount int, work func(T) error) error {
	if workerCount <= 0 {
		return nil
	}
	if len(items) == 0 {
		return nil
	}
	// TODO: Add callback system

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var wg sync.WaitGroup
	jobs := make(chan T)
	errs := make(chan error, 1)

	worker := func() {
		defer wg.Done()

		for {
			select {
			case <-ctx.Done():
				return
			case item, ok := <-jobs:
				if !ok {
					return
				}
				if err := work(item); err != nil {
					select {
					case errs <- err:
						cancel()
					default:
					}
					return
				}
			}
		}
	}

	for range workerCount {
		wg.Add(1)
		go worker()
	}

	for _, item := range items {
		select {
		case <-ctx.Done():
			close(jobs)
			wg.Wait()
			return <-errs
		case jobs <- item:
		}
	}

	close(jobs)
	wg.Wait()

	select {
	case err := <-errs:
		return err
	default:
		return nil
	}
}

// TaskWorkerCount returns a bounded worker count for a task based on app configuration & item count
func TaskWorkerCount(app *state.State, itemCount int, fallback int) int {
	if itemCount <= 0 {
		return 0
	}

	workerCount := fallback
	if app != nil && app.Config.PostgresPoolSize > 0 {
		workerCount = app.Config.PostgresPoolSize
	}

	if workerCount < 1 {
		return 1
	}
	if workerCount > itemCount {
		return itemCount
	}
	return workerCount
}
