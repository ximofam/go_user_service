package worker

import "context"

type Job func(ctx context.Context) error

type Pool interface {
	Run()
	Submit(job Job) error
	Shutdown()
}
