package scheduler

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/coreos/etcd/client"
	metrics "github.com/rcrowley/go-metrics"
)

const (
	meticNameWatchFailure  = "cronjobWatchFailure"
	meticNameWatchUpdated  = "cronjobWatchUpdated"
	metricNameParseFailure = "cronjobParseFailure"
	metricNameParseOK      = "cronjobParseOK"
)

// Tasks store all the structural data needed to run the cronjobs
type Tasks struct {
	BasePath string
	Archive  map[string]TaskDefinition
	Keys     client.KeysAPI
}

// NewTasksWithKeys returns a new Task archive
func NewTasksWithKeys(keys client.KeysAPI) *Tasks {
	return &Tasks{
		BasePath: "/cron/Jobs",
		Keys:     keys,
		Archive:  make(map[string]TaskDefinition),
	}
}

// Watch watches the storage for changes
func (t *Tasks) Watch() *Tasks {
	watcherOpts := client.WatcherOptions{AfterIndex: 0, Recursive: true}
	go func() {
		w := t.Keys.Watcher(t.BasePath, &watcherOpts)
		for {
			if _, err := w.Next(context.Background()); err != nil {
				m := metrics.GetOrRegisterCounter(meticNameWatchFailure, metrics.DefaultRegistry)
				m.Inc(1)
				continue
			}
			m := metrics.GetOrRegisterCounter(meticNameWatchUpdated, metrics.DefaultRegistry)
			m.Inc(1)
			t.Load()
		}
	}()
	return t
}

// Load loads tasks from etcd
func (t *Tasks) Load() error {
	nodes, err := t.Keys.Get(context.Background(), t.BasePath, nil)
	if err != nil {
		return err
	}
	value := nodes.Node.Nodes
	return t.parseTasks(value)
}

// GetTask returns a single task definition
func (t *Tasks) GetTask(name string) (td *TaskDefinition, err error) {
	if item, ok := t.Archive[name]; ok {
		return &item, nil
	}
	err = fmt.Errorf("task %s was not found", name)
	return
}

func (t *Tasks) parseTasks(nodes client.Nodes) error {
	t.Archive = map[string]TaskDefinition{}
	for _, task := range nodes {
		var cfg TaskDefinition
		if err := json.Unmarshal([]byte(task.Value), &cfg); err != nil {
			m := metrics.GetOrRegisterCounter(metricNameParseFailure, metrics.DefaultRegistry)
			m.Inc(1)
			continue
		}
		cfg.Key = task.Key
		t.Archive[cfg.Name] = cfg
		m := metrics.GetOrRegisterCounter(metricNameParseOK, metrics.DefaultRegistry)
		m.Inc(1)
	}
	return nil
}
