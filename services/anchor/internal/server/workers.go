package server

type WorkerPool[C any] struct {
	handler func(C)
	tasks   chan C
}

func (p *WorkerPool[C]) worker() {
	for task := range p.tasks {
		p.handler(task)
	}
}

func NewWorkerPool[C any](size int, handler func(C)) *WorkerPool[C] {
	if size <= 0 {
		panic("worker pool size must be positive")
	}
	if handler == nil {
		panic("handler function must not be nil")
	}

	pool := &WorkerPool[C]{
		tasks:   make(chan C, 100),
		handler: handler,
	}
	for range size {
		go pool.worker()
	}
	return pool
}
