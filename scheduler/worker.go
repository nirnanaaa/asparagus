package scheduler

import (
	"fmt"
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

// SubmitErr delivers an error to a source provider
func (w *Worker) SubmitErr(work *provider.Task, err error) {
	if work.SourceProvider != nil {
		work.SourceProvider.TaskError(work, err)
	}
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
				desiredProvider, ok := w.ExecutionProviders[work.Executor]
				if !ok {
					w.SubmitErr(work, fmt.Errorf("provider %s not found inside list of available providers", work.Executor))
					continue
				}
				if work.SourceProvider != nil {
					if err := work.SourceProvider.TaskStarted(work); err != nil {
						w.SubmitErr(work, err)
						continue
					}
				}
				start := time.Now()
				var rVal provider.Response
				if err := desiredProvider.Execute(work.ExecutionConfig, &rVal); err != nil {
					logrus.WithError(err).Error("Error inside execution provider.")
					w.LogReporters(err, work, start, &rVal)
					w.SubmitErr(work, err)
					continue
				}
				if work.SourceProvider != nil {
					if err := work.SourceProvider.TaskDone(work); err != nil {
						w.SubmitErr(work, err)
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
