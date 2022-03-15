package workerPool

import "time"

type WorkRequest struct {
	Name       string
	Delay      time.Duration
	ExpiryTime time.Time
}
