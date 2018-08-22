package local_test

import (
	"os"
	"testing"

	"github.com/Sirupsen/logrus"
	"github.com/nirnanaaa/asparagus/scheduler/provider"
	"github.com/nirnanaaa/asparagus/scheduler/provider/local"
)

func initConf() local.Config {
	cfg := local.NewConfig()
	cfg.Enabled = true
	cfg.EnabledOutput = false
	return cfg
}

func TestLocal_RequestTouch(t *testing.T) {
	cfg := initConf()
	cfg.EnabledOutput = true
	req := local.NewExecutionProvider(cfg, logrus.New())
	if err := req.Execute(map[string]interface{}{
		"Command": "touch",
		"Arguments": []string{
			"/tmp/test.file",
		},
	}, &provider.Response{}); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat("/tmp/test.file"); os.IsNotExist(err) {
		t.Fatal(err)
	}
	os.Remove("/tmp/test.file")
}

func TestLocal_RequestString(t *testing.T) {
	cfg := initConf()
	req := local.NewExecutionProvider(cfg, logrus.New())
	if err := req.Execute("touch /tmp/test.file", &provider.Response{}); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat("/tmp/test.file"); os.IsNotExist(err) {
		t.Fatal(err)
	}
	os.Remove("/tmp/test.file")
}

func TestLocal_RequestNotEnabled(t *testing.T) {
	cfg := initConf()
	cfg.Enabled = false
	req := local.NewExecutionProvider(cfg, logrus.New())
	if err := req.Execute(map[string]interface{}{
		"Command": "touch",
		"Arguments": []string{
			"/tmp/test.file",
		},
	}, &provider.Response{}); err == nil {
		t.Fatal("no error was thrown, but executor was disabled")
	}
}

func TestLocal_RequestNonExistingPath(t *testing.T) {
	cfg := initConf()
	req := local.NewExecutionProvider(cfg, logrus.New())
	err := req.Execute(map[string]interface{}{
		"Command": "txbawquasiudasojdi",
	}, &provider.Response{})
	if err == nil {
		t.Fatal("no error was thrown, but executor was disabled")
	}
	if err.Error() != "exec: \"txbawquasiudasojdi\": executable file not found in $PATH" {
		t.Fatalf("unexpected error message. Want \"exec: \"txbawquasiudasojdi\": executable file not found in $PATH\". Got \"%s\"", err.Error())
	}
}

func TestLocal_WrongInputType(t *testing.T) {
	cfg := initConf()
	req := local.NewExecutionProvider(cfg, logrus.New())
	err := req.Execute([]string{}, &provider.Response{})
	if err == nil {
		t.Fatal("no error was thrown, but wrong input type used")
	}
	if err.Error() != "unknown input type on executor: []string" {
		t.Fatalf("unexpected error thrown. expected unknown input type on executor: []string got %s", err.Error())
	}
}

func TestLocal_WrongInputTypeString(t *testing.T) {
	cfg := initConf()
	req := local.NewExecutionProvider(cfg, logrus.New())
	err := req.Execute("", &provider.Response{})
	if err == nil {
		t.Fatal("no error was thrown, but wrong input type used")
	}
	if err.Error() != "command has no parts" {
		t.Fatalf("unexpected error thrown. expected command has no parts got %s", err.Error())
	}
}
