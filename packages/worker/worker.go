package worker

import (
	"sync"
)

type JobWorker interface {
	Execute() interface{}
}

type WorkerPool struct {
	numberWorkers int
	results       chan interface{}
	jobs          chan JobWorker
	wg            sync.WaitGroup
}

func NewWorkerPool(numberWorkers int) *WorkerPool {
	worker := &WorkerPool{
		numberWorkers: numberWorkers,
		results:       make(chan interface{}),
		jobs:          make(chan JobWorker),
		wg:            sync.WaitGroup{},
	}

	worker.start()

	return worker
}

func (wp *WorkerPool) Add(job JobWorker) {
	wp.jobs <- job
}

func (wp *WorkerPool) Result() <-chan interface{} {
	return wp.results
}

func (wp *WorkerPool) Stop() {
	close(wp.jobs)
	wp.wg.Wait()
	close(wp.results)
}

func (wp *WorkerPool) start() {
	for range wp.numberWorkers {
		wp.wg.Add(1)
		go func() {
			defer wp.wg.Done()
			wp.work()
		}()
	}
}

func (wp *WorkerPool) work() {
	for job := range wp.jobs {
		result := job.Execute()
		wp.results <- result
	}
}
