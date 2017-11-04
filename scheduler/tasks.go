package scheduler

import (
	"fmt"

	"github.com/Sirupsen/logrus"
	"github.com/coreos/etcd/client"
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
	Archive         map[string]TaskDefinition
	Keys            client.KeysAPI
	Tasks           map[string]*provider.Task
	SourceProviders []provider.SourceProvider
}

// NewTasksWithKeys returns a new Task archive
func NewTasksWithKeys(config *Config, keys client.KeysAPI) *Tasks {
	source := []provider.SourceProvider{
		crontab.NewSourceProvider(config.CrontabSource),
	}

	return &Tasks{
		BasePath:        "/cron/Jobs",
		Keys:            keys,
		Archive:         make(map[string]TaskDefinition),
		Tasks:           make(map[string]*provider.Task),
		SourceProviders: source,
	}
}

// Load loads tasks from etcd
func (t *Tasks) Load() error {
	for _, provider := range t.SourceProviders {
		if err := provider.Read(&t.Tasks); err != nil {
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
