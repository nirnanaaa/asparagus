package scheduler

import (
	"fmt"

	"github.com/Sirupsen/logrus"
	"github.com/nirnanaaa/asparagus/scheduler/provider"
	"github.com/nirnanaaa/asparagus/scheduler/provider/crontab"
)

const (
	meticNameWatchFailure  = "cronjobWatchFailure"
	meticNameWatchUpdated  = "cronjobWatchUpdated"
	metricNameParseFailure = "cronjobParseFailure"
	metricNameParseOK      = "cronjobParseOK"
)

// Tasks store all the structural data needed to run the cronjobs
type Tasks struct {
	BasePath        string
	Tasks           map[string]*provider.Task
	SourceProviders []provider.SourceProvider
}

// NewSourceConfig returns a new Task archive
func NewSourceConfig(config *Config) *Tasks {
	source := []provider.SourceProvider{
		crontab.NewSourceProvider(config.CrontabSource),
	}

	return &Tasks{
		Tasks:           make(map[string]*provider.Task),
		SourceProviders: source,
	}
}

// Load loads tasks from etcd
func (t *Tasks) Load() error {
	for _, pv := range t.SourceProviders {
		pv.OnTaskUpdate(func(tx *provider.Task) error {
			t.Tasks[tx.Name] = tx
			return nil
		})
		if err := pv.Read(); err != nil {
			return err
		}
	}
	return nil
}

// GetTask returns a single task definition
func (t *Tasks) GetTask(name string) (td *provider.Task, err error) {
	if item, ok := t.Tasks[name]; ok {
		return item, nil
	}
	err = fmt.Errorf("task %s was not found", name)
	return
}

// DebugTasks dumps all tasks inside the task archive
func (t *Tasks) DebugTasks(log *logrus.Logger) {
	for _, task := range t.Tasks {
		log.Debugf("Task %s registered. After Task: %s. Has Execution Schedule: %s", task.Name, task.AfterTask, task.Expression)
	}
}
