package background_test

import (
	"context"
	"sync"
	"sync/atomic"
	"testing"

	"github.com/microsoft/typescript-go/internal/project/background"
	"github.com/wow-look-at-my/testify/require"
	"github.com/wow-look-at-my/testify/assert"
)

func TestQueue(t *testing.T) {
	t.Parallel()
	t.Run("BasicEnqueue", func(t *testing.T) {
		t.Parallel()
		q := background.NewQueue()
		defer q.Close()

		executed := false
		q.Enqueue(context.Background(), func(ctx context.Context) {
			executed = true
		})

		q.Wait()

		assert.True(t, executed)
	})

	t.Run("MultipleTasksExecution", func(t *testing.T) {
		t.Parallel()
		q := background.NewQueue()
		defer q.Close()

		var counter int64
		numTasks := 10

		for range numTasks {
			q.Enqueue(context.Background(), func(ctx context.Context) {
				atomic.AddInt64(&counter, 1)
			})
		}

		q.Wait()

		require.Equal(t, atomic.LoadInt64(&counter), int64(numTasks))
	})

	t.Run("NestedEnqueue", func(t *testing.T) {
		t.Parallel()
		q := background.NewQueue()
		defer q.Close()

		var executed []string
		var mu sync.Mutex

		q.Enqueue(context.Background(), func(ctx context.Context) {
			mu.Lock()
			executed = append(executed, "parent")
			mu.Unlock()

			q.Enqueue(ctx, func(childCtx context.Context) {
				mu.Lock()
				executed = append(executed, "child")
				mu.Unlock()
			})
		})

		q.Wait()

		mu.Lock()
		defer mu.Unlock()

		require.Equal(t, len(executed), 2)
	})

	t.Run("ClosedQueueRejectsNewTasks", func(t *testing.T) {
		t.Parallel()
		q := background.NewQueue()
		q.Close()

		executed := false
		q.Enqueue(context.Background(), func(ctx context.Context) {
			executed = true
		})

		q.Wait()

		assert.True(t, !executed, "Task should not execute after queue is closed")
	})
}
