package scheduler

import (
	"fmt"

	"github.com/Sirupsen/logrus"
	"github.com/nirnanaaa/asparagus/scheduler/provider"
	"github.com/nirnanaaa/asparagus/scheduler/provider/crontab"
	"github.com/nirnanaaa/asparagus/scheduler/provider/etcd"
	"github.com/nirnanaaa/asparagus/scheduler/provider/sql"
	uuid "github.com/satori/go.uuid"
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
		etcd.NewSourceProvider(config.ETCDSource),
		sql.NewSourceProvider(config.SQLSource),
	}

	logrus.Debugf("Discovered source providers: %v", source)
	return &Tasks{
		Tasks:           make(map[string]*provider.Task),
		SourceProviders: source,
	}
}

// Load loads tasks from etcd
func (t *Tasks) Load() error {
	for idx, pv := range t.SourceProviders {
		pv.OnTaskUpdate(func(tx *provider.Task) error {
			lock.Lock()
			tx.SourceProvider = t.SourceProviders[idx]
			if tx.Name == "" {
				tx.Name = uuid.NewV4().String()
			}
			t.Tasks[tx.Name] = tx
			lock.Unlock()
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
	lock.RLock()
	defer lock.RUnlock()
	if item, ok := t.Tasks[name]; ok {
		return item, nil
	}
	err = fmt.Errorf("task %s was not found", name)
	return
}

// DebugTasks dumps all tasks inside the task archive
func (t *Tasks) DebugTasks(log *logrus.Logger) {
	lock.RLock()
	defer lock.RUnlock()
	for _, task := range t.Tasks {
		log.Debugf("Task %s registered. After Task: %s. Has Execution Schedule: %s", task.Name, task.AfterTask, task.Expression)
	}
}
