package crontab_test

import (
	"path/filepath"
	"testing"
	"time"

	"github.com/nirnanaaa/asparagus/scheduler/provider"
	"github.com/nirnanaaa/asparagus/scheduler/provider/crontab"
)

func readCrontabs(t *testing.T) []provider.Task {
	cfg := crontab.NewConfig()
	cfg.Enabled = true
	ct, err := filepath.Abs("crontab")
	if err != nil {
		t.Fatal(err)
	}
	cfg.Crontab = ct
	cts := crontab.NewSourceProvider(cfg)
	tabs := map[string]*provider.Task{}
	cts.OnTaskUpdate(func(arg0 *provider.Task) error {
		tabs[arg0.Name] = arg0
		return nil
	})
	if err := cts.Read(); err != nil {
		t.Fatal(err)
	}
	time.Sleep(100 * time.Millisecond)
	tabx := []provider.Task{}
	for _, tab := range tabs {
		tabx = append(tabx, *tab)
	}
	return tabx
}

func TestCrontabRead(t *testing.T) {
	tabs := readCrontabs(t)
	if len(tabs) < 1 {
		t.Fatalf("undefined amount of cronjobs found: %d expected 1", len(tabs))
	}
}
func TestCrontabExecutionSchedule(t *testing.T) {
	tabs := readCrontabs(t)
	tab := tabs[0]
	if tab.Expression != "* * * * *" {
		t.Fatal("Expression wasn't parsed correctly")
	}
}

func TestCrontabParseExecutor(t *testing.T) {
	tabs := readCrontabs(t)
	tab := tabs[0]
	if tab.Executor != "http" {
		t.Fatal("Executor wasn't parsed correctly")
	}
}

func TestCrontabParseRunning(t *testing.T) {
	tabs := readCrontabs(t)
	tab := tabs[0]
	if tab.Running != true {
		t.Fatal("Running wasn't parsed correctly")
	}
}

func TestCrontabParseName(t *testing.T) {
	tabs := readCrontabs(t)
	tab := tabs[0]
	if tab.Name != "HTTP Test" {
		t.Fatal("Name wasn't parsed correctly")
	}
}

func TestCrontabMapStr(t *testing.T) {
	tabs := readCrontabs(t)
	tab := tabs[0]
	mapping, ok := tab.ExecutionConfig.(map[string]string)
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
