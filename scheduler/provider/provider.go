package provider

import (
	"time"
)

// SourceProvider is used to get tasks
type SourceProvider interface {
	Read(*map[string]*Task) error
	// Stop() error
}

// ExecutionProvider is used for executing cronjobs
type ExecutionProvider interface {
	Execute(interface{}) error
}

// Task represents a task to be run
type Task struct {
	Name            string      `json:"Name"`
	Service         string      `json:"Service,omitempty"`
	Running         bool        `json:"Running,omitempty"`
	Expression      string      `json:"Expression,omitempty"`
	LastRunAt       time.Time   `json:"LastRunAt,omitempty"`
	AfterTask       string      `json:"AfterTask,omitempty"`
	Key             string      `json:"-"`
	Executor        string      `json:"executor"`
	ExecutionConfig interface{} `json:"ProviderConfig"`
}
