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
	Name            string         `json:"name"`
	Running         bool           `json:"is_running,omitempty"`
	Expression      string         `json:"schedule,omitempty"`
	LastRunAt       time.Time      `json:"last_run_at,omitempty"`
	AfterTask       string         `json:"after_task,omitempty"`
	Key             string         `json:"-"`
	SourceProvider  SourceProvider `json:"-"`
	Executor        string         `json:"executor"`
	ExecutionConfig interface{}    `json:"execution_config"`
}

// Response contains data that executors can return
type Response struct {
	StatusCode int    `json:"status_code"`
	Response   []byte `json:"response"`
}
