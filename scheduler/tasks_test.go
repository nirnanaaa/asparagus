package scheduler_test

import (
	"path/filepath"
	"testing"
	"time"

	"github.com/nirnanaaa/asparagus/scheduler"

	"github.com/coreos/etcd/client"
)

func getConn() client.KeysAPI {
	cfg := client.Config{
		Endpoints:               []string{"http://localhost:4001"},
		HeaderTimeoutPerRequest: time.Second,
	}
	c, err := client.New(cfg)
	if err != nil {
		panic(err)
	}
	return client.NewKeysAPI(c)
}

func seedEtcdData() {
	// c := getConn()
}

func TestLoadOK(t *testing.T) {
	seedEtcdData()
	c := getConn()
	config := scheduler.NewConfig()
	fp, _ := filepath.Abs("provider/crontab/crontab")
	config.CrontabSource.Crontab = fp
	cli := scheduler.NewTasksWithKeys(config, c)
	if err := cli.Load(); err != nil {
		t.Fatal(err.Error())
	}
	time.Sleep(100 * time.Millisecond)
	if _, err := cli.GetTask("ExampleJob"); err != nil {
		t.Fatal(err.Error())
	}
	if _, err := cli.GetTask("ExampleJob1"); err == nil {
		t.Fatal("Task that shouldn't exist exists.")
	}
}
