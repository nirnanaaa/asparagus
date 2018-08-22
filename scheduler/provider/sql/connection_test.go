package sql_test

import (
	"testing"

	"github.com/nirnanaaa/asparagus/scheduler/provider"
	"github.com/nirnanaaa/asparagus/scheduler/provider/sql"
	uuid "github.com/satori/go.uuid"
)

func getConnection(t *testing.T) *sql.Connection {
	conf := sql.NewConfig()
	conf.DatabaseDriver = "sqlite3"
	testuuid := uuid.NewV4().String()
	conf.DatabaseURI = "/tmp/" + testuuid
	conf.Enabled = true
	conn := sql.NewConnection(conf)
	if err := conn.Connect(); err != nil {
		t.Fatal(err)
	}
	if conn.SQL == nil {
		t.Fatal("client wasn't connected after handshake")
	}
	conn.SQL.Exec("insert into cronjobs values ('TestTask', 'TestTask', '0', '* * * * *', '2018-08-22 10:02:25.067050808+02:00', '', 'http', '{\"URL\": \"https://httpbin.org/get\", \"Method\": \"GET\"}', '', '10', '0');")
	return conn
}

func TestSQLConnection_New(t *testing.T) {
	getConnection(t)
}

func TestSQLConnection_WriteTask(t *testing.T) {
	conn := getConnection(t)
	task := provider.Task{
		Name:       "TestTask",
		Key:        "TestTask",
		Expression: "@daily",
		Executor:   "dummy",
	}
	if err := conn.WriteTask(&task); err != nil {
		t.Fatal(err)
	}
}

func TestSQLConnection_WatchTask(t *testing.T) {
	conn := getConnection(t)
	task := provider.Task{
		Name:       "TestTask",
		Key:        "TestTask",
		Expression: "@daily",
		Executor:   "dummy",
	}
	done := make(chan bool)
	go conn.WatchTasks(func(t1 provider.Task) {
		if t1.Name != task.Name {
			t.Fatalf("wrong item was saved, expected name %s got %s", task.Name, t1.Name)
		}
		done <- true
	})
	task.Running = true
	task.ExecutionConfig = map[string]string{
		"test": "true",
	}
	if err := conn.WriteTask(&task); err != nil {
		t.Fatal(err)
	}
	<-done
	close(done)
}
