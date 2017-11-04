package scheduler

import (
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/nirnanaaa/asparagus/scheduler/provider"
)

// NewWorker creates, and returns a new Worker object. Its only argument
// is a channel that the worker can add itself to whenever it is done its
// work.
func NewWorker(id int, workerQueue chan chan *provider.Task, eps map[string]provider.ExecutionProvider) Worker {
	// Create, and return the worker.
	worker := Worker{
		ID:                 id,
		Work:               make(chan *provider.Task),
		WorkerQueue:        workerQueue,
		QuitChan:           make(chan bool),
		ExecutionProviders: eps,
	}

	return worker
}

// Worker defines a working horse for tasks
type Worker struct {
	ID                 int
	Work               chan *provider.Task
	WorkerQueue        chan chan *provider.Task
	QuitChan           chan bool
	ExecutionProviders map[string]provider.ExecutionProvider
	Callbacks          map[string]func(*provider.Task) error
}

// Start function "starts" the worker by starting a goroutine, that is
// an infinite "for-select" loop.
func (w *Worker) Start() {
	go func() {
		for {
			// Add ourselves into the worker queue.
			w.WorkerQueue <- w.Work

			select {
			case work := <-w.Work:
				provider, ok := w.ExecutionProviders[work.Executor]
				if !ok {
					// TODO: Set Error
					continue
				}
				work.Running = true
				if err := provider.Execute(work.ExecutionConfig); err != nil {
					logrus.WithError(err).Error("Error inside execution provider.")
					continue
				}
				work.Running = false
				work.LastRunAt = time.Now()
			case <-w.QuitChan:
				return
			}
		}
	}()
}

// Stop tells the worker to stop listening for work requests.
//
// Note that the worker will only stop *after* it has finished its work.
func (w *Worker) Stop() {
	go func() {
		w.QuitChan <- true
	}()
}
