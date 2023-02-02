package task

import (
	"context"
	"errors"
	"testing"
	"time"

	"mini.task.queue.io/utils"
)

func Test_ManagerWithContext(t *testing.T) {
	manager := NewManager("convert-"+utils.RandStringRunes(), 3)
	ctx1, cancel1 := context.WithCancel(context.Background())
	defer cancel1()
	go manager.Run(ctx1)

	ctx2, cancel := context.WithCancel(context.Background())

	go func() {
		time.Sleep(2 * time.Second)
		cancel()
	}()

	f := func(ctx context.Context, sleep int) error {
		timer := time.NewTimer(time.Second * time.Duration(sleep))
		timer.Stop()
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-timer.C:
			return errors.New("sleep done")
		}
	}

	err := manager.Process(ctx2, func(ctx context.Context) error {
		return f(ctx, 3)
	})

	if !errors.Is(context.Canceled, err) {
		t.Errorf("want ctx cancel error, got %v", err)
	}
}

func Test_ManagerWithMultiJobs(t *testing.T) {
	manager := NewManager("convert-"+utils.RandStringRunes(), 2)
	ctx1, cancel1 := context.WithCancel(context.Background())
	defer cancel1()
	go manager.Run(ctx1)

	f := func(ctx context.Context, sleep int) error {
		timer := time.NewTimer(time.Second * time.Duration(sleep))
		timer.Stop()
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-timer.C:
			return nil
		}
	}

	for i := 0; i < 2; i++ {
		go func() {
			err := manager.Process(context.Background(), func(ctx context.Context) error {
				return f(ctx, 4)
			})
			if err != nil {
				t.Errorf("want ctx nil error, got %v", err)
			}
		}()
	}

	timeoutCtx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	err := manager.Process(timeoutCtx, func(ctx context.Context) error {
		return f(ctx, 1)
	})

	if !errors.Is(context.DeadlineExceeded, err) {
		t.Errorf("want ctx timeout error, got %v", err)
	}
}
