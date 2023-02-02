package task

import (
	"context"

	"mini.task.queue.io/utils"
)

type Job struct {
	id  string
	ctx context.Context
	f   func(ctx context.Context) error
	out chan error
}

func NewJob(ctx context.Context, f func(ctx context.Context) error, out chan error) *Job {
	return &Job{
		id:  "job-" + utils.RandStringRunes(),
		ctx: ctx,
		f:   f,
		out: out,
	}
}
