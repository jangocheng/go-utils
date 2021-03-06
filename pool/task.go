package pool

import (
	"sync"
	"time"

	"github.com/gxxgle/go-utils/log"
)

type TaskPool struct {
	name    string
	stopped bool
	tasks   chan func()

	sync.WaitGroup
}

func NewTaskPool(name string, size int) *TaskPool {
	out := &TaskPool{
		name:    name,
		stopped: false,
		tasks:   make(chan func()),
	}

	for i := 0; i < size; i++ {
		go out.loopRun()
	}

	return out
}

func (tp *TaskPool) AddTask(f func()) {
	tp.tasks <- f
}

func (tp *TaskPool) IsRunning() bool {
	return !tp.stopped
}

func (tp *TaskPool) Stop() {
	if tp.stopped {
		return
	}

	start := time.Now()
	defer func() {
		log.Infow("stop taskpool", "taskpool_name", tp.name, "cost", time.Since(start))
	}()

	tp.stopped = true
	close(tp.tasks)
	tp.Wait()
}

func (tp *TaskPool) loopRun() {
	tp.Add(1)
	defer tp.Done()

	for !tp.stopped {
		tp.run()
	}
}

func (tp *TaskPool) run() {
	defer func() {
		if err := recover(); err != nil {
			log.Errorw("task crash", "err", err)
		}
	}()

	task, ok := <-tp.tasks
	if !ok {
		return
	}

	task()
}
