package scheduler_test

import (
	"path/filepath"
	"testing"
	"time"

	"github.com/nirnanaaa/asparagus/scheduler"
)

//
// func getConn() client.KeysAPI {
// 	cfg := client.Config{
// 		Endpoints:               []string{"http://localhost:4001"},
// 		HeaderTimeoutPerRequest: time.Second,
// 	}
// 	c, err := client.New(cfg)
// 	if err != nil {
// 		panic(err)
// 	}
// 	return client.NewKeysAPI(c)
// }

func TestLoadOK(t *testing.T) {
	config := scheduler.NewConfig()
	fp, _ := filepath.Abs("provider/crontab/crontab")
	config.CrontabSource.Crontab = fp
	cli := scheduler.NewSourceConfig(config)
	if err := cli.Load(); err != nil {
		t.Fatal(err.Error())
	}
	time.Sleep(1 * time.Second)
	if _, err := cli.GetTask("ExampleJob"); err != nil {
		t.Fatal(err.Error())
	}
	if _, err := cli.GetTask("ExampleJob1"); err == nil {
		t.Fatal("Task that shouldn't exist exists.")
	}
}
