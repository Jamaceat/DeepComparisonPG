package concurrent

import (
	"context"
	"sync"
	"time"
)

// WorkerPool represents a pool of workers for concurrent processing
type WorkerPool struct {
	workers    int
	jobQueue   chan Job
	resultChan chan Result
	wg         *sync.WaitGroup
	ctx        context.Context
	cancel     context.CancelFunc
}

// Job represents a unit of work to be processed
type Job struct {
	ID       string
	TaskType string
	Data     interface{}
	Timeout  time.Duration
}

// Result represents the result of a job
type Result struct {
	JobID string
	Data  interface{}
	Error error
}

// NewWorkerPool creates a new worker pool
func NewWorkerPool(workers int, jobQueueSize int) *WorkerPool {
	ctx, cancel := context.WithCancel(context.Background())
	return &WorkerPool{
		workers:    workers,
		jobQueue:   make(chan Job, jobQueueSize),
		resultChan: make(chan Result, jobQueueSize),
		wg:         &sync.WaitGroup{},
		ctx:        ctx,
		cancel:     cancel,
	}
}

// Start initializes and starts the worker pool
func (wp *WorkerPool) Start() {
	for i := 0; i < wp.workers; i++ {
		wp.wg.Add(1)
		go wp.worker(i)
	}
}

// Stop gracefully stops the worker pool
func (wp *WorkerPool) Stop() {
	close(wp.jobQueue)
	wp.cancel()
	wp.wg.Wait()
	close(wp.resultChan)
}

// SubmitJob submits a job to the worker pool
func (wp *WorkerPool) SubmitJob(job Job) {
	select {
	case wp.jobQueue <- job:
	case <-wp.ctx.Done():
	}
}

// GetResults returns the results channel
func (wp *WorkerPool) GetResults() <-chan Result {
	return wp.resultChan
}

// worker processes jobs from the job queue
func (wp *WorkerPool) worker(id int) {
	defer wp.wg.Done()

	for {
		select {
		case job, ok := <-wp.jobQueue:
			if !ok {
				return
			}

			// Process job with timeout
			result := wp.processJobWithTimeout(job)

			select {
			case wp.resultChan <- result:
			case <-wp.ctx.Done():
				return
			}

		case <-wp.ctx.Done():
			return
		}
	}
}

// processJobWithTimeout processes a job with a timeout
func (wp *WorkerPool) processJobWithTimeout(job Job) Result {
	if job.Timeout == 0 {
		job.Timeout = 30 * time.Second // Default timeout
	}

	ctx, cancel := context.WithTimeout(wp.ctx, job.Timeout)
	defer cancel()

	resultChan := make(chan Result, 1)

	go func() {
		result := wp.processJob(job)
		select {
		case resultChan <- result:
		case <-ctx.Done():
		}
	}()

	select {
	case result := <-resultChan:
		return result
	case <-ctx.Done():
		return Result{
			JobID: job.ID,
			Error: ctx.Err(),
		}
	}
}

// processJob processes a single job (to be extended by specific implementations)
func (wp *WorkerPool) processJob(job Job) Result {
	// This is a placeholder - specific job processors will be implemented
	// in the comparator package
	return Result{
		JobID: job.ID,
		Data:  job.Data,
		Error: nil,
	}
}
