package hw05_parallel_execution //nolint:golint,stylecheck

import (
	"errors"
	"sync"
	"sync/atomic"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

type TaskContext struct {
	errorCount *int32
	taskCount  *int32
	wg         sync.WaitGroup
	tasks      *[]Task
	maxError   int32
}

// Run starts tasks in N goroutines and stops its work when receiving M errors from tasks.
func Run(tasks []Task, n int, m int) error {
	if m <= 0 {
		return ErrErrorsLimitExceeded
	}

	var errorPoint int32 = 0
	var startPoint int32 = -1

	context := TaskContext{
		errorCount: &errorPoint,
		taskCount:  &startPoint,
		wg:         sync.WaitGroup{},
		tasks:      &tasks,
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

	for {
		localErrorCount := atomic.LoadInt32(context.errorCount)
		if localErrorCount >= context.maxError {
			return
		}

		currentTask := atomic.AddInt32(context.taskCount, 1)

		if int(currentTask) >= len(*context.tasks) {
			return
		}

		task := (*context.tasks)[currentTask]

		err := task()
		if err != nil {
			atomic.AddInt32(context.errorCount, 1)
		}
	}
}
