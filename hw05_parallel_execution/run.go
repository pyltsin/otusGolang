package hw05_parallel_execution //nolint:golint,stylecheck

import (
	"errors"
	"sync"
	"sync/atomic"
)

type Task func() error

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type TaskContext struct {
	errorCount *int32
	in         chan Task
	wg         sync.WaitGroup
	maxError   int32
}

// Run starts tasks in N goroutines and stops its work when receiving M errors from tasks.
func Run(tasks []Task, n int, m int) error {
	if m <= 0 {
		return ErrErrorsLimitExceeded
	}

	var errorPoint int32 = 0

	var in = make(chan Task, len(tasks))

	for _, task := range tasks {
		in <- task
	}

	close(in)

	context := TaskContext{
		errorCount: &errorPoint,
		in:         in,
		wg:         sync.WaitGroup{},
		maxError:   int32(m),
	}

	for i := 0; i < n; i++ {
		context.wg.Add(1)
		go worker(&context)
	}

	context.wg.Wait()

	localErrorCount := int(*context.errorCount)
	if localErrorCount >= m {
		return ErrErrorsLimitExceeded
	}

	return nil
}

func worker(context *TaskContext) {
	defer context.wg.Done()

	for task := range context.in {
		localErrorCount := atomic.LoadInt32(context.errorCount)
		if localErrorCount >= context.maxError {
			return
		}
		err := task()
		if err != nil {
			atomic.AddInt32(context.errorCount, 1)
		}
	}
}
