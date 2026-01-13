package concurrency

import (
	"context"
	"log"
	"sync"
	"sync/atomic"
	"time"
)

type Job func(ctx context.Context) error

type WorkerPool struct {
	minWorker    int
	maxWorker    int
	jos          chan Job
	results      chan error
	rateLimiter  <-chan time.Time
	activeWorker int64
	totalJobs    int64
	failedJobs   int64
	wg           sync.WaitGroup
}

func NewWorkerPool(minWorker int, maxWorker int, rateLimite time.Duration) *WorkerPool {
	return &WorkerPool{
		minWorker:   minWorker,
		maxWorker:   maxWorker,
		jos:         make(chan Job),
		results:     make(chan error),
		rateLimiter: time.Tick(rateLimite),
	}
}

func (wp *WorkerPool) Start(ctx context.Context) {
	for i := 0; i < wp.minWorker; i++ {
		wp.startWorker(ctx)
	}
}

func (wp *WorkerPool) startWorker(ctx context.Context) {
	wp.wg.Add(1)
	atomic.AddInt64(&wp.activeWorker, 1)

	go func(workerID int) {
		defer wp.wg.Done()
		defer atomic.AddInt64(&wp.activeWorker, -1)

		// 🧯 PANIC RECOVERY
		defer func() {
			if r := recover(); r != nil {
				log.Printf("[Worker %d] PANIC recovered: %v", workerID, r)
			}
		}()

		for {
			select {
			case <-ctx.Done():
				return

			case <-wp.rateLimiter: // 🚦 Rate limiting

				job, ok := <-wp.jos
				if !ok {
					return
				}

				atomic.AddInt64(&wp.totalJobs, 1)

				err := wp.executeWithRetry(ctx, job, 3)
				if err != nil {
					atomic.AddInt64(&wp.failedJobs, 1)
				}

				wp.results <- err
			}
		}
	}(int(wp.activeWorker))
}

func (wp *WorkerPool) executeWithRetry(
	ctx context.Context,
	job Job,
	maxRetry int,
) error {

	var err error

	for attempt := 1; attempt <= maxRetry; attempt++ {
		err = job(ctx)
		if err == nil {
			return nil
		}

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(time.Duration(attempt) * 100 * time.Millisecond):
			// exponential backoff
		}
	}

	return err
}
