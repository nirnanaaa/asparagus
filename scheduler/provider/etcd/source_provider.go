package etcd

import (
	"time"

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
	go func() {
		p.QuitChan <- true
	}()
	return nil
}

// Read reads from FS
func (p *SourceProvider) Read() error {
	if !p.Config.Enabled {
		return nil
	}
	if err := p.Connection.Connect(); err != nil {
		return err
	}
	var t []provider.Task
	if err := p.Connection.ReadAll(&t); err != nil {
		return err
	}
	for _, task := range t {
		p.TaskFlow <- task
	}
	go p.Connection.WatchTasks(func(t *provider.Task) {
		p.TaskFlow <- *t
	})
	return nil
}
