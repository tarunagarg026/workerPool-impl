package workerPool

import (
	"fmt"
	"time"
)

type Worker struct {
	ID          int
	WorkerQueue chan WorkRequest
}

// NewWorker creates, and returns a new Worker object. Its only argument
// is a channel that the worker can add itself to whenever it is done its
// work.

func NewWorker(id int, workerQueue chan WorkRequest) Worker {
	// Create, and return the worker.
	worker := Worker{
		ID:          id,
		WorkerQueue: workerQueue}

	return worker
}

// This function "starts" the worker by starting a goroutine, that is
// an infinite "for-select" loop.
func (w *Worker) StartWorker() {
	go func() {
		for {
			// Add ourselves into the worker queue.
			select {
			case work := <-w.WorkerQueue:
				if time.Now().Before(work.ExpiryTime) {
					// Receive a work request.
					fmt.Printf("worker%d: Received work request, delaying for %f seconds\n", w.ID, work.Delay.Seconds())

					time.Sleep(work.Delay)
					fmt.Printf("worker%d: Hello, %s!\n", w.ID, work.Name)
				} else {
					fmt.Printf("worker%d stopping\n", w.ID)
				}
			}
		}
	}()
}
