package etcd_test

import (
	"testing"

	"github.com/nirnanaaa/asparagus/scheduler/provider"
	"github.com/nirnanaaa/asparagus/scheduler/provider/etcd"
)

func getConnection(t *testing.T) *etcd.Connection {
	conf := etcd.NewConfig()
	conn := etcd.NewConnection(conf)
	if err := conn.Connect(); err != nil {
		t.Fatal(err)
	}
	if conn.KAPI == nil {
		t.Fatal("client wasn't connected after handshake")
	}
	return conn
}

func TestETCDConnection_New(t *testing.T) {
	getConnection(t)
}

func TestETCDConnection_WriteTask(t *testing.T) {
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
	t1 := provider.Task{}
	if err := conn.ReadTask(task.Key, &t1); err != nil {
		t.Fatal(err)
	}
	if t1.Name != task.Name {
		t.Fatalf("wrong item was saved, expected name %s got %s", task.Name, t1.Name)
	} else if t1.Expression != task.Expression {
		t.Fatalf("wrong item was saved, expected expression %s got %s", task.Expression, t1.Expression)
	} else if t1.Executor != task.Executor {
		t.Fatalf("wrong item was saved, expected executor %s got %s", task.Executor, t1.Executor)
	}
}

func TestETCDConnection_WatchTask(t *testing.T) {
	conn := getConnection(t)
	task := provider.Task{
		Name:       "TestTask",
		Key:        "TestTask",
		Expression: "@daily",
		Executor:   "dummy",
	}
	done := make(chan bool)
	go conn.WatchTasks(func(t1 *provider.Task) {
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
