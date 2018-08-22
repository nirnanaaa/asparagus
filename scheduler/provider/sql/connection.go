package sql

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/mattn/go-sqlite3"
	"github.com/nirnanaaa/asparagus/scheduler/provider"
)

// Connection defines a wrapper for the actual etcd connection
type Connection struct {
	Config    Config
	SQL       *sql.DB
	WatchChan *time.Ticker
	Cancel    chan struct{}
}

// NewConnection creates a new etcd connection
func NewConnection(config Config) *Connection {
	return &Connection{
		Config: config,
		Cancel: make(chan struct{}),
	}
}
func (c *Connection) Close() error {
	if c.WatchChan != nil {
		c.WatchChan.Stop()
	}
	fmt.Printf("Before closing")
	close(c.Cancel)
	fmt.Printf("after closing")
	if c.SQL != nil {
		return c.SQL.Close()
	}
	return nil
}

const createStatementSQlite = `
CREATE TABLE IF NOT EXISTS cronjobs (
	name	TEXT NOT NULL,
	key	TEXT,
	is_running	INTEGER NOT NULL DEFAULT 0,
	schedule	TEXT NOT NULL,
	last_run_at	datetime,
	after_task	TEXT,
	executor	TEXT NOT NULL,
	executor_parameters	TEXT,
	last_output	TEXT,
	max_retries	INTEGER DEFAULT 10,
	retry_count	INTEGER DEFAULT 0,
	PRIMARY KEY(name)
) WITHOUT ROWID;`

const createStatementMysql = "CREATE TABLE IF NOT EXISTS `cronjobs` (`key` varchar(100) NOT NULL,`name` varchar(100) NOT NULL,`is_running` tinyint(1) NOT NULL,`schedule` varchar(32) NOT NULL,`last_run_at` datetime DEFAULT '1970-01-01 00:00:00',`after_task` varchar(100) DEFAULT '',`executor` varchar(100) NOT NULL,`executor_parameters` json,`last_output` mediumtext,`max_retries` int(10) DEFAULT 10,`retry_count` int(10) DEFAULT 0,PRIMARY KEY (name))"

// Connect connects the etcd target
func (c *Connection) Connect() error {
	c.WatchChan = time.NewTicker(5 * time.Second)
	db, err := sql.Open(c.Config.DatabaseDriver, c.Config.DatabaseURI)
	if err != nil {
		return nil
	}
	statement := createStatementMysql
	if c.Config.DatabaseDriver == "sqlite3" {
		statement = createStatementSQlite
	}
	if _, err := db.Exec(statement); err != nil {
		return err
	}
	c.SQL = db
	return nil
}

// WriteTask persists a task inside etcd
func (c *Connection) WriteTask(t *provider.Task) error {
	_, err := c.SQL.Exec(`
	update cronjobs
	set is_running = ?,
	last_run_at = ?,
	last_output = ?,
	retry_count = ?
	where name = ?`, t.Running, t.LastRunAt, t.LastError, t.CurrentRetryCount, t.Name)
	return err
}

// ReadTask reads a task from etcd
func (c *Connection) ReadTask(key string, task *provider.Task) error {
	err := c.SQL.QueryRow("select `name`, `key`, `is_running`, `schedule`, `last_run_at`, `after_task`, `executor`, `executor_parameters`, `last_output`, `max_retries`, `retry_count` FROM cronjobs where name = ?", task.Name).Scan(
		&task.Name,
		&task.Key,
		&task.Running,
		&task.Expression,
		&task.LastRunAt,
		&task.AfterTask,
		&task.Executor,
		&task.ExecutionConfig,
		&task.LastError,
		&task.MaxRetries,
		&task.CurrentRetryCount,
	)
	if err != nil {
		return err
	}

	executorConfig, executionConfigErr := c.parseExecutionConfig(task.ExecutionConfig)
	if err != nil {
		err = executionConfigErr
		return err
	}
	task.ExecutionConfig = executorConfig
	return nil
}

func (c *Connection) parseExecutionConfig(input interface{}) (output map[string]interface{}, err error) {
	output = map[string]interface{}{}
	iface, ok := input.([]uint8)
	if !ok {
		err = fmt.Errorf("interface conversion failed")
		return
	}
	err = json.Unmarshal(iface, &output)
	return
}

// ReadAll reads all cronjobs inside a directory
func (c *Connection) ReadAll() (tasks []provider.Task, err error) {
	tasks = []provider.Task{}
	rows, err := c.SQL.Query("select `name`, `key`, `is_running`, `schedule`, `last_run_at`, `after_task`, `executor`, `executor_parameters`, `last_output`, `max_retries`, `retry_count` FROM cronjobs")
	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		var task provider.Task
		err = rows.Scan(
			&task.Name,
			&task.Key,
			&task.Running,
			&task.Expression,
			&task.LastRunAt,
			&task.AfterTask,
			&task.Executor,
			&task.ExecutionConfig,
			&task.LastError,
			&task.MaxRetries,
			&task.CurrentRetryCount,
		)
		if err != nil {
			return
		}
		executorConfig, executionConfigErr := c.parseExecutionConfig(task.ExecutionConfig)
		if err != nil {
			err = executionConfigErr
			return
		}
		task.ExecutionConfig = executorConfig
		tasks = append(tasks, task)
	}

	err = rows.Err()
	return
}

// WatchTasks looks for new keys inside etcd
func (c *Connection) WatchTasks(cb func(t provider.Task)) {
	for {
		select {
		case <-c.WatchChan.C:
			tasks, err := c.ReadAll()
			if err != nil {
				return
			}
			for _, task := range tasks {
				cb(task)
			}
		case <-c.Cancel:
			return
		}
	}
}
