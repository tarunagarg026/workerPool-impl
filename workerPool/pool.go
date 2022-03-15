package workerPool

import (
	"fmt"
	"runtime"
	"sync"
	"time"
)

var WorkerQueue chan WorkRequest
var startPool sync.Once
var getPoolData sync.Once

type Pool struct {
	NumOfWorkers int
	WorkQueue    chan WorkRequest
	CPUs         int
	Duration     time.Duration
}

func GetPoolConfig(noOfWorkers, numCPU int, duration time.Duration) *Pool {
	var pool *Pool
	getPoolData.Do(
		func() {
			pool = &Pool{
				WorkQueue:    make(chan WorkRequest, 100),
				CPUs:         numCPU,
				NumOfWorkers: noOfWorkers,
				Duration:     duration,
			}
		})
	return pool
}

func (pool *Pool) collector(workRequest WorkRequest) {
	pool.WorkQueue <- workRequest
	fmt.Println("Work request queued.")
}

func InitializeWorkerPool(HTTPAddr string, pool *Pool) {
	startPool.Do(
		func() {
			runtime.GOMAXPROCS(pool.CPUs)
			startDispatcher(pool)
			StartServer(HTTPAddr, pool)
		})
}

func startDispatcher(pool *Pool) {

	fmt.Println("Starting the dispatcher")
	// First, initialize the channel we are going to but the workers' work channels into.
	WorkerQueue = make(chan WorkRequest, pool.NumOfWorkers)

	// Now, create all of our workers.
	for i := 0; i < pool.NumOfWorkers; i++ {
		fmt.Println("Starting worker", i+1)
		worker := NewWorker(i+1, pool.WorkQueue)
		worker.StartWorker()
	}

	go func() {
		for {
			select {
			case work := <-pool.WorkQueue:
				fmt.Println("Received work request")
				go func() {
					fmt.Println("Dispatching work request")
					pool.WorkQueue <- work
				}()
			}
		}
	}()
}
