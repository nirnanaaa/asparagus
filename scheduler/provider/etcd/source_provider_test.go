package etcd_test

import (
	"testing"

	"github.com/nirnanaaa/asparagus/scheduler/provider"
	"github.com/nirnanaaa/asparagus/scheduler/provider/etcd"
)

func prepareTasks(t *testing.T, conn *etcd.Connection) {
	t1 := provider.Task{
		Key:        "HTTPTEST1",
		Name:       "HTTP Test",
		Expression: "* * * * *",
		Executor:   "http",
		Running:    true,
		ExecutionConfig: map[string]string{
			"URL": "https://httpbin.org/get",
		},
	}
	if err := conn.WriteTask(&t1); err != nil {
		t.Fatal(err)
	}
	t2 := provider.Task{
		Key:        "HTTPTEST2",
		Name:       "HTTP Test 1",
		Expression: "* * * * *",
		Executor:   "http",
		AfterTask:  "HTTP Test",
		ExecutionConfig: map[string]string{
			"URL": "https://httpbin.org/get",
		},
	}
	if err := conn.WriteTask(&t2); err != nil {
		t.Fatal(err)
	}
	t3 := provider.Task{
		Key:        "HTTPTEST3",
		Name:       "ExampleJob",
		Expression: "* * * * *",
		AfterTask:  "HTTP Test",
		Executor:   "http",
		ExecutionConfig: map[string]string{
			"URL": "https://httpbin.org/get",
		},
	}
	if err := conn.WriteTask(&t3); err != nil {
		t.Fatal(err)
	}
}
func readCrontabs(t *testing.T) map[string]*provider.Task {
	cfg := etcd.NewConfig()
	cfg.Enabled = true
	cfg.SourceFolder = "/test1"
	cts := etcd.NewSourceProvider(cfg)
	cts.Connection.Connect()
	prepareTasks(t, cts.Connection)
	tabs := map[string]*provider.Task{}
	done := make(chan bool)
	cts.OnTaskUpdate(func(arg0 *provider.Task) error {
		tabs[arg0.Name] = arg0
		if len(tabs) >= 3 {
			done <- true
		}
		return nil
	})
	if err := cts.Read(); err != nil {
		t.Fatal(err)
	}
	<-done
	return tabs
}

func TestETCDRead(t *testing.T) {
	tabs := readCrontabs(t)
	if len(tabs) < 1 {
		t.Fatalf("undefined amount of cronjobs found: %d expected 1", len(tabs))
	}
}
func TestETCDExecutionSchedule(t *testing.T) {
	tabs := readCrontabs(t)
	tab := tabs["HTTP Test"]
	if tab.Expression != "* * * * *" {
		t.Fatal("Expression wasn't parsed correctly")
	}
}

func TestETCDParseExecutor(t *testing.T) {
	tabs := readCrontabs(t)
	tab := tabs["HTTP Test"]
	if tab.Executor != "http" {
		t.Fatal("Executor wasn't parsed correctly")
	}
}

func TestETCDParseRunning(t *testing.T) {
	tabs := readCrontabs(t)
	tab := tabs["HTTP Test"]
	if tab.Running != true {
		t.Fatal("Running wasn't parsed correctly")
	}
}

func TestETCDParseName(t *testing.T) {
	tabs := readCrontabs(t)
	tab := tabs["HTTP Test"]
	if tab.Name != "HTTP Test" {
		t.Fatal("Name wasn't parsed correctly")
	}
}

func TestETCDMapStr(t *testing.T) {
	tabs := readCrontabs(t)
	tab := tabs["HTTP Test"]
	mapping, ok := tab.ExecutionConfig.(map[string]interface{})
	if !ok {
		t.Fatal("mapping wasn't a map")
	}
	if _, ok := mapping["URL"]; !ok {
		t.Fatal("URL found")
	}
	val, _ := mapping["URL"]
	if val != "https://httpbin.org/get" {
		t.Fatal("URL wasn't parsed correctly")
	}
}
