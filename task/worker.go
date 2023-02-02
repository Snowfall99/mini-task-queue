package task

type Worker struct {
	id   string
	in   <-chan *Job
	done <-chan byte
}

func NewWorker(id string, in <-chan *Job, done <-chan byte) *Worker {
	return &Worker{
		id:   id,
		in:   in,
		done: done,
	}
}

func (w *Worker) Run() {
	for {
		select {
		case <-w.done:
			return
		case job := <-w.in:
			job.out <- job.f(job.ctx)
		}
	}
}
