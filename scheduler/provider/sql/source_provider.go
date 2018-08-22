package sql

import (
	"fmt"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/nirnanaaa/asparagus/scheduler/provider"
)

// SourceProvider creates a new src provider
type SourceProvider struct {
	Config     Config
	Connection *Connection
	TaskFlow   chan provider.Task
	QuitChan   chan bool
}

// NewSourceProvider creates a new src provider
func NewSourceProvider(config Config) *SourceProvider {
	conn := NewConnection(config)
	return &SourceProvider{
		Config:     config,
		Connection: conn,
		TaskFlow:   make(chan provider.Task),
		QuitChan:   make(chan bool),
	}
}

// String returns the Providers identity
func (p SourceProvider) String() string {
	return fmt.Sprintf("SQL Provider, enabled: %t, driver: %s, connection: %s", p.Config.Enabled, p.Config.DatabaseDriver, p.Config.DatabaseURI)
}

// TaskError will be called when an error occured
func (p SourceProvider) TaskError(t *provider.Task, err error) error {
	t.Running = false
	t.LastError = err.Error()
	t.CurrentRetryCount++
	nextDuration := provider.CalculateBackoffForAttempt(float64(t.CurrentRetryCount))
	t.LastRunAt = time.Now().Add(nextDuration)
	logrus.Debug("resetting running after task has resulted in an error")
	if err := p.Connection.WriteTask(t); err != nil {
		return err
	}
	return nil
}

// TaskStarted will be executed when a task is started
func (p SourceProvider) TaskStarted(t *provider.Task) error {
	t.Running = true
	if err := p.Connection.WriteTask(t); err != nil {
		return err
	}
	return nil
}

// TaskDone will be executed when a task is done
func (p SourceProvider) TaskDone(t *provider.Task) error {
	t.LastRunAt = time.Now()
	t.Running = false
	t.LastError = ""
	t.CurrentRetryCount = 0
	if err := p.Connection.WriteTask(t); err != nil {
		return err
	}
	return nil
}

// OnTaskUpdate runs when a task gets updated
func (p SourceProvider) OnTaskUpdate(fn func(*provider.Task) error) {
	go func() {
		for {
			select {
			case task := <-p.TaskFlow:
				fn(&task)
			case <-p.QuitChan:
				return
			}
		}
	}()
}

// Stop quits reading tasks
func (p SourceProvider) Stop() error {
	p.Connection.Close()
	go func() {
		p.QuitChan <- true
	}()
	return nil
}

// Read reads from FS
func (p *SourceProvider) Read() (err error) {
	if !p.Config.Enabled {
		return nil
	}
	if err := p.ensureConnection(); err != nil {
		return err
	}

	tasks, err := p.Connection.ReadAll()
	if err != nil {
		return
	}
	for _, task := range tasks {
		p.TaskFlow <- task
	}
	go p.Connection.WatchTasks(func(t provider.Task) {
		p.TaskFlow <- t
	})
	return nil
}

func (p *SourceProvider) ensureConnection() (err error) {
	if p.Connection.SQL == nil {
		err = p.Connection.Connect()
	}
	return
}
