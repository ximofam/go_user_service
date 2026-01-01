package worker

import (
	"context"
	"log"
	"sync"
)

type PoolV1 struct {
	jobs        chan Job
	wg          sync.WaitGroup
	concurrency int
}

func NewPoolV1(concurrency, queueSize int) *PoolV1 {
	return &PoolV1{
		jobs:        make(chan Job, queueSize),
		wg:          sync.WaitGroup{},
		concurrency: concurrency,
	}
}

func (p *PoolV1) Run() {
	for range p.concurrency {
		p.wg.Add(1)
		go p.worker()
	}
}

func (p *PoolV1) worker() {
	defer p.wg.Done()

	defer func() {
		if r := recover(); r != nil {
			log.Printf("Worker panic: %v", r)
		}
	}()

	for job := range p.jobs {
		if err := job(context.Background()); err != nil {
			log.Println(err)
		}
	}
}

func (p *PoolV1) Submit(job Job) error {
	p.jobs <- job

	return nil
}

func (p *PoolV1) Shutdown() {
	close(p.jobs)

	p.wg.Wait()
	log.Println("Worker pool shutdown gracefully")
}
