package task

import (
	"context"
	"errors"
	"testing"
	"time"

	"mini.task.queue.io/utils"
)

func Test_WokerWithContext(t *testing.T) {
	in := make(chan *Job)
	done := make(chan byte)
	out := make(chan error)
	worker := NewWorker("test-"+utils.RandStringRunes(), in, done)

	go worker.Run()

	// context with timeout
	timeoutCtx, cancel1 := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel1()
	in <- NewJob(timeoutCtx, func(ctx context.Context) error {
		timer := time.NewTimer(time.Second * 5)
		defer timer.Stop()
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-timer.C:
			return errors.New("sleep 5 secondes")
		}
	}, out)

	select {
	case <-done:
		t.Errorf("want ctx timeout error, got woker done")
	case err := <-out:
		if !errors.Is(context.DeadlineExceeded, err) {
			t.Errorf("want ctx timeout error, got %v", err)
		}
	}

	// context with cancel
	cancelCtx, cancancel2 := context.WithCancel(context.Background())

	in <- NewJob(cancelCtx, func(ctx context.Context) error {
		timer := time.NewTimer(time.Second * 5)
		defer timer.Stop()
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-timer.C:
			return errors.New("sleep 5 secondes")
		}
	}, out)

	cancancel2()
	select {
	case <-done:
		t.Errorf("want ctx cancel error, got woker done")
	case err := <-out:
		if !errors.Is(context.Canceled, err) {
			t.Errorf("want ctx cancel error, got %v", err)
		}
	}

	// normal
	timeoutCtx1, cancel3 := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel3()
	in <- NewJob(timeoutCtx1, func(ctx context.Context) error {
		timer := time.NewTimer(time.Second * 1)
		defer timer.Stop()
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-timer.C:
			return nil
		}
	}, out)

	select {
	case <-done:
		t.Errorf("want nil error, got woker done")
	case err := <-out:
		if err != nil {
			t.Errorf("want nil error, got %v", err)
		}
	}

	close(done)
}
