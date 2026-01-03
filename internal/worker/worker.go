package worker

import (
	"errors"
	"sync"
)

var (
	instance Pool
	once     sync.Once
)

func Init(concurrency, queueSize int) {
	once.Do(func() {
		instance = NewPoolV1(concurrency, queueSize)
		instance.Run()
	})
}

func Submit(job Job) error {
	if instance == nil {
		return errors.New("The Worker Pool has not been initialized")
	}

	return instance.Submit(job)
}

func Shutdown() {
	if instance != nil {
		instance.Shutdown()
	}
}
