package scheduler

import (
	"github.com/Sirupsen/logrus"
	"github.com/nirnanaaa/asparagus/scheduler/provider"
)

// Dispatcher for spinning up workers and distributing work
type Dispatcher struct {
	WorkerQueue chan chan *provider.Task
	WorkQueue   chan *provider.Task
	Workers     []*Worker
}

// StartDispatcher starts the dispatching provider
func StartDispatcher(nworkers int, eps map[string]provider.ExecutionProvider) *Dispatcher {
	d := &Dispatcher{
		WorkQueue:   make(chan *provider.Task, 100),
		WorkerQueue: make(chan chan *provider.Task, nworkers),
		// Workers:     make([]*Worker, nworkers),
	}
	logrus.Debugf("Starting %d workers", nworkers)
	for i := 0; i < nworkers; i++ {
		worker := NewWorker(i+1, d.WorkerQueue, eps)
		worker.Start()
		d.Workers = append(d.Workers, &worker)
	}

	go func() {
		for {
			select {
			case work := <-d.WorkQueue:
				go func() {
					worker := <-d.WorkerQueue
					worker <- work
				}()
			}
		}
	}()
	return d
}

// Stop stops all workers
func (d *Dispatcher) Stop() error {
	for _, w := range d.Workers {
		w.Stop()
	}
	return nil
}
