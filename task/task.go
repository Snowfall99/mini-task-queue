package task

import (
	"context"
	"fmt"

	"mini.task.queue.io/utils"
)

var DefaultManager *Manager

func RunDefault(ctx context.Context, workerNum int) error {
	if DefaultManager != nil {
		return fmt.Errorf("task demon already running")
	}

	DefaultManager = NewManager("default-"+utils.RandStringRunes(), workerNum)
	return DefaultManager.Run(ctx)
}

func Process(ctx context.Context, f func(context.Context) error) error {
	if DefaultManager == nil {
		return fmt.Errorf("task demon not running")
	}
	return DefaultManager.Process(ctx, f)
}
