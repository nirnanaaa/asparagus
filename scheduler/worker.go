package scheduler

import (
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/nirnanaaa/asparagus/metric/adapters"
	"github.com/nirnanaaa/asparagus/scheduler/provider"
)

// NewWorker creates, and returns a new Worker object. Its only argument
// is a channel that the worker can add itself to whenever it is done its
// work.
func NewWorker(id int, workerQueue chan chan *provider.Task, eps map[string]provider.ExecutionProvider, reporters []adapters.Reporter) Worker {
	// Create, and return the worker.
	worker := Worker{
		ID:                 id,
		Work:               make(chan *provider.Task),
		WorkerQueue:        workerQueue,
		QuitChan:           make(chan bool),
		ExecutionProviders: eps,
		Reporters:          reporters,
	}

	return worker
}

// Worker defines a working horse for tasks
type Worker struct {
	ID                 int
	Work               chan *provider.Task
	WorkerQueue        chan chan *provider.Task
	Reporters          []adapters.Reporter
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
				provider1, ok := w.ExecutionProviders[work.Executor]
				if !ok {
					// TODO: Set Error
					continue
				}
				if work.SourceProvider != nil {
					if err := work.SourceProvider.TaskStarted(work); err != nil {
						continue
					}
				}
				start := time.Now()
				var rVal provider.Response
				if err := provider1.Execute(work.ExecutionConfig, &rVal); err != nil {
					logrus.WithError(err).Error("Error inside execution provider.")
					w.LogReporters(err, work, start, &rVal)
					continue
				}
				if work.SourceProvider != nil {
					if err := work.SourceProvider.TaskDone(work); err != nil {
						continue
					}
				}
				w.LogReporters(nil, work, start, &rVal)
			case <-w.QuitChan:
				return
			}
		}
	}()
}

// LogReporters fans out to reporters.
func (w *Worker) LogReporters(err error, t *provider.Task, start time.Time, r *provider.Response) {
	end := time.Now()
	isErr := err != nil
	debugResult := ""
	if isErr {
		debugResult = err.Error()
	}
	took := time.Since(start)
	for _, reporter := range w.Reporters {
		if !reporter.SupportsLogging() {
			continue
		}
		line := adapters.LogEvent{
			Name:        t.Name,
			IsError:     isErr,
			StartDate:   start,
			EndDate:     end,
			ResultDebug: debugResult,
			StatusCode:  r.StatusCode,
			Took:        took.String(),
			Result:      string(r.Response),
		}
		go reporter.LogResult(&line)
	}
}

// Stop tells the worker to stop listening for work requests.
//
// Note that the worker will only stop *after* it has finished its work.
func (w *Worker) Stop() {
	go func() {
		w.QuitChan <- true
	}()
}
