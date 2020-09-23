package hw05_parallel_execution //nolint:golint,stylecheck

import (
	"fmt"
	"math/rand"
	"strconv"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.uber.org/goleak"
)

type WorkerCount struct {
	worker int
	task   int
	err    int
}

func TestRun(t *testing.T) {
	defer goleak.VerifyNone(t)

	t.Run("if were errors in first M tasks, than finished not more N+M tasks", func(t *testing.T) {
		for _, count := range [...]WorkerCount{
			{
				worker: 10,
				task:   50,
				err:    23,
			},
			{
				worker: 10,
				task:   2,
				err:    1,
			},
			{
				worker: 1,
				task:   10,
				err:    1,
			},
			{
				worker: 10,
				task:   10,
				err:    1,
			},
		} {

			tasksCount := count.task
			tasks := make([]Task, 0, tasksCount)

			var runTasksCount int32

			for i := 0; i < tasksCount; i++ {
				err := fmt.Errorf("error from task %d", i)
				tasks = append(tasks, func() error {
					time.Sleep(time.Millisecond * time.Duration(rand.Intn(100)))
					atomic.AddInt32(&runTasksCount, 1)
					return err
				})
			}

			workersCount := count.worker
			maxErrorsCount := count.err
			result := Run(tasks, workersCount, maxErrorsCount)

			require.Equal(t, ErrErrorsLimitExceeded, result)
			require.LessOrEqual(t, runTasksCount, int32(workersCount+maxErrorsCount), "extra tasks were started")
		}
	})

	t.Run("tasks without errors", func(t *testing.T) {
		for _, count := range [...]WorkerCount{
			{
				worker: 5,
				task:   50,
			},
			{
				worker: 10,
				task:   10,
			},
			{
				worker: 1,
				task:   10,
			},
			{
				worker: 10,
				task:   10,
			},
		} {

			tasksCount := count.task
			tasks := make([]Task, 0, tasksCount)

			var runTasksCount int32
			var sumTime time.Duration

			for i := 0; i < tasksCount; i++ {
				taskSleep := time.Millisecond * time.Duration(rand.Intn(100))
				sumTime += taskSleep

				tasks = append(tasks, func() error {
					time.Sleep(taskSleep)
					atomic.AddInt32(&runTasksCount, 1)
					return nil
				})
			}

			workersCount := count.worker
			maxErrorsCount := 1

			start := time.Now()
			result := Run(tasks, workersCount, maxErrorsCount)
			elapsedTime := time.Since(start)
			require.Nil(t, result)

			require.Equal(t, runTasksCount, int32(tasksCount), "not all tasks were completed")
			if count.worker > 2 {
				require.LessOrEqual(t, int64(elapsedTime), int64(sumTime/2), "tasks were run sequentially? - worker: "+strconv.Itoa(count.worker))
			}
		}
	})
}

func TestAdditional(t *testing.T) {
	defer goleak.VerifyNone(t)

	t.Run("if error count =0, than ErrErrorsLimitExceeded", func(t *testing.T) {
		tasksCount := 40
		tasks := make([]Task, 0, tasksCount)

		workersCount := 10
		maxErrorsCount := 0
		result := Run(tasks, workersCount, maxErrorsCount)

		require.Equal(t, ErrErrorsLimitExceeded, result)

	})
	t.Run("if error count <0, than ErrErrorsLimitExceeded", func(t *testing.T) {
		tasksCount := 40
		tasks := make([]Task, 0, tasksCount)

		workersCount := 10
		maxErrorsCount := -1
		result := Run(tasks, workersCount, maxErrorsCount)

		require.Equal(t, ErrErrorsLimitExceeded, result)

	})
}
