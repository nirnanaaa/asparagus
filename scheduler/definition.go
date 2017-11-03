package scheduler

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/coreos/etcd/client"
)

// TaskDefinition represents a task to be run
type TaskDefinition struct {
	Name       string    `json:"Name"`
	Service    string    `json:"Service,omitempty"`
	Running    bool      `json:"Running,omitempty"`
	Expression string    `json:"Expression,omitempty"`
	LastRunAt  time.Time `json:"LastRunAt,omitempty"`
	AfterTask  string    `json:"AfterTask,omitempty"`
	URI        string    `json:"URI"`
	Absolute   bool      `json:"URIIsAbsolute,omitempty"`
	Key        string    `json:"-"`
}

// Done marks a task as done
func (t *TaskDefinition) Done(kapi client.KeysAPI) {
	t.Running = false
	t.LastRunAt = time.Now()
	t.UpdateAttributes(kapi)
}

// SetRunning marks a task as "in progress"
func (t *TaskDefinition) SetRunning(kapi client.KeysAPI, running bool) {
	t.Running = running
	t.UpdateAttributes(kapi)
}

// UpdateAttributes changes the current attributes on a task
func (t *TaskDefinition) UpdateAttributes(kapi client.KeysAPI) {
	scfg, err := json.Marshal(&t)
	if err != nil {
		fmt.Printf(err.Error())
		return
	}
	if _, err := kapi.Set(context.Background(), t.Key, string(scfg), nil); err != nil {
		fmt.Printf(err.Error())
		return
	}
}
