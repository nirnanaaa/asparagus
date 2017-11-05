package scheduler_test

import (
	"testing"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/nirnanaaa/asparagus/scheduler"
)

func TestService_Initialization(t *testing.T) {
	cfg := scheduler.NewConfig()
	cfg.CrontabSource.Enabled = false
	srv := scheduler.NewService(cfg, logrus.New())
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
	srv.RegisterExecutionProvider("dummy", dummy)
	go srv.Start()
	time.Sleep(1 * time.Second)
	if err := srv.Stop(); err != nil {
		t.Fatal("Service didn't stop correctly")
	}
}
