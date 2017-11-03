package scheduler_test

import (
	"context"
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
	c := getConn()
	c.Set(context.Background(), "/cron/Config", `{"SchedulerTickDuration":1000,"Enabled":true,"SecretKey":"cf86bb5624de56b1d88c205ac261b1cc6cb2a529f0c408475f007d4d00c2f3aa75fb5c2a0c33c3eb7e3a628acd7fa02aa2f5774614201b780a36542f2fc67ead"}`, nil)
	c.Set(context.Background(), "/cron/Jobs/ExampleJob", `{"Name":"ExampleJob","Expression":"* * * * *","URIIsAbsolute":true,"URI":"http://stage-logx.miha-bodytec.com/api/v1/cron/ManagementRatioGenerate"}`, nil)
}

func TestLoadOK(t *testing.T) {
	seedEtcdData()
	c := getConn()
	cli := scheduler.NewTasksWithKeys(c)
	if err := cli.Load(); err != nil {
		t.Fatal(err.Error())
	}
	if _, err := cli.GetTask("ExampleJob"); err != nil {
		t.Fatal(err.Error())
	}
	if _, err := cli.GetTask("ExampleJob1"); err == nil {
		t.Fatal("Task that shouldn't exist exists.")
	}
	cli.Watch()
	time.Sleep(1 * time.Second)
	c.Set(context.Background(), "/cron/Jobs/ExampleJob", `{"Name":"ExampleJob1","Expression":"* * * * *","URIIsAbsolute":true,"URI":"http://stage-logx.miha-bodytec.com/api/v1/cron/ManagementRatioGenerate"}`, nil)
	time.Sleep(500 * time.Millisecond)
	if _, err := cli.GetTask("ExampleJob1"); err != nil {
		t.Fatal(err.Error())
	}
	if _, err := cli.GetTask("ExampleJob"); err == nil {
		t.Fatal("Task that shouldn't exist exists.")
	}
}
