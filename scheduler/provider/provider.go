package provider

import (
	"time"
)

// SourceProvider is used to get tasks
type SourceProvider interface {
	Read() error
	OnTaskUpdate(func(*Task) error)
	Stop() error
	TaskStarted(*Task) error
	TaskDone(*Task) error
}

// ExecutionProvider is used for executing cronjobs
type ExecutionProvider interface {
	Execute(interface{}, *Response) error
}

// Task represents a task to be run
type Task struct {
	Name            string         `json:"Name"`
	Service         string         `json:"Service,omitempty"`
	Running         bool           `json:"Running,omitempty"`
	Expression      string         `json:"Expression,omitempty"`
	LastRunAt       time.Time      `json:"LastRunAt,omitempty"`
	AfterTask       string         `json:"AfterTask,omitempty"`
	Key             string         `json:"-"`
	SourceProvider  SourceProvider `json:"-"`
	Executor        string         `json:"executor"`
	ExecutionConfig interface{}    `json:"ProviderConfig"`
}

// Response contains data that executors can return
type Response struct {
	StatusCode int    `json:"status_code"`
	Response   []byte `json:"response"`
}
