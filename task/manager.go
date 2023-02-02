package task

import (
	"context"
	"fmt"

	"mini.task.queue.io/utils"
)

type Manager struct {
	id        string
	in        chan *Job
	done      chan byte
	workerNum int
}

func NewManager(id string, workerNum int) *Manager {
	return &Manager{
		id:        id,
		in:        make(chan *Job),
		done:      make(chan byte),
		workerNum: workerNum,
	}
}

func (m *Manager) Run(ctx context.Context) error {
	if m.workerNum <= 0 {
		return fmt.Errorf("task manager with %d workers", m.workerNum)
	}
	for i := 0; i < m.workerNum; i++ {
		worker := NewWorker(m.id+"-"+utils.RandStringRunes(), m.in, m.done)
		go worker.Run()
	}
	<-ctx.Done()
	close(m.done)
	return nil
}

func (m *Manager) Process(ctx context.Context, f func(ctx context.Context) error) error {
	out := make(chan error)
	job := NewJob(ctx, f, out)
	select {
	case <-job.ctx.Done():
		return job.ctx.Err()
	case m.in <- job:
	}
	return <-job.out
}
