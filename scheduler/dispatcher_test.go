package scheduler_test

import (
	"testing"

	"github.com/nirnanaaa/asparagus/scheduler"
	"github.com/nirnanaaa/asparagus/scheduler/provider"
)

type DummyProvider struct {
	Callback func(interface{})
}

func NewDummyProvider(callback func(interface{})) *DummyProvider {
	return &DummyProvider{
		callback,
	}
}

func (d *DummyProvider) Execute(iface interface{}) error {
	d.Callback(iface)
	return nil
}

func TestDispatcher_Initialization(t *testing.T) {
	ok := make(chan bool)
	defer close(ok)
	dummy := NewDummyProvider(func(arg0 interface{}) {
		vars, isOk := arg0.(string)
		if !isOk {
			t.Fatalf("Executor didn't received provided string input. Instead: %T", vars)
		}
		if vars != "ok" {
			t.Fatalf("Unexpected string received. Wants: ok, Got: %s", vars)
		}
		ok <- true
	})
	providers := map[string]provider.ExecutionProvider{"dummy": dummy}
	dispatcher := scheduler.StartDispatcher(1, providers)
	task := provider.Task{
		ExecutionConfig: "ok",
		Executor:        "dummy",
	}
	dispatcher.WorkQueue <- &task
	<-ok
	dispatcher.Stop()
}
