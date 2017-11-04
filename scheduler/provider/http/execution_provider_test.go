package http_test

import (
	"testing"

	"github.com/nirnanaaa/asparagus/scheduler/provider/http"
)

func initConf() http.Config {
	cfg := http.NewConfig()
	cfg.Enabled = true
	cfg.SignJWT = false
	cfg.DebugResponse = false
	cfg.LogHTTPStatus = true
	return cfg
}

func TestHTTPRequestGET(t *testing.T) {
	cfg := initConf()
	req := http.NewExecutionProvider(cfg)
	if err := req.Execute(map[string]string{
		"URL":    "https://httpbin.org/get",
		"Method": "GET",
	}); err != nil {
		t.Fatal(err)
	}
}
func TestHTTPRequestGETPostEndpoint(t *testing.T) {
	cfg := initConf()
	req := http.NewExecutionProvider(cfg)
	err := req.Execute(map[string]string{
		"URL":    "https://httpbin.org/post",
		"Method": "GET",
	})
	if err == nil {
		t.Fatal("Request should have failed for status code 405")
	}
	if err.Error() != "request failed with status code 405" {
		t.Fatalf("Expected error message \"request failed with status code 405\", but got: \"%s\"", err.Error())
	}
}
func TestHTTPRequestPOST(t *testing.T) {
	cfg := initConf()
	cfg.DebugResponse = false
	cfg.LogHTTPStatus = false
	req := http.NewExecutionProvider(cfg)
	if err := req.Execute(map[string]string{
		"URL":    "https://httpbin.org/post",
		"Method": "POST",
	}); err != nil {
		t.Fatal(err)
	}
}
func TestHTTPRequestPUT(t *testing.T) {
	cfg := initConf()
	cfg.DebugResponse = false
	cfg.LogHTTPStatus = false
	cfg.DebugResponse = false
	cfg.LogHTTPStatus = false
	req := http.NewExecutionProvider(cfg)
	if err := req.Execute(map[string]string{
		"URL":    "https://httpbin.org/put",
		"Method": "PUT",
	}); err != nil {
		t.Fatal(err)
	}
}
func TestHTTPRequestDELETE(t *testing.T) {
	cfg := initConf()
	cfg.DebugResponse = false
	cfg.LogHTTPStatus = false
	req := http.NewExecutionProvider(cfg)
	if err := req.Execute(map[string]string{
		"URL":    "https://httpbin.org/delete",
		"Method": "DELETE",
	}); err != nil {
		t.Fatal(err)
	}
}
func TestHTTPRequest409(t *testing.T) {
	cfg := initConf()
	cfg.DebugResponse = true
	req := http.NewExecutionProvider(cfg)
	err := req.Execute(map[string]string{
		"URL":    "https://httpbin.org/status/409",
		"Method": "GET",
	})
	if err == nil {
		t.Fatal("Request should have failed for status code 409")
	}
	if err.Error() != "request failed with status code 409" {
		t.Fatalf("Expected error message \"request failed with status code 409\", but got: \"%s\"", err.Error())
	}
}

func TestHTTPNoMap(t *testing.T) {
	cfg := initConf()
	req := http.NewExecutionProvider(cfg)
	err := req.Execute("https://httpbin.org/get")
	if err == nil {
		t.Fatal("should throw an error, because a map was required.")
	}
	if err.Error() != "unknown input type on executor: string" {
		t.Fatal("Error message isn't correct: %s", err.Error())
	}
}

func TestHTTPNotEnabled(t *testing.T) {
	cfg := initConf()
	cfg.Enabled = false
	req := http.NewExecutionProvider(cfg)
	err := req.Execute(map[string]string{
		"URL":    "https://httpbin.org/get",
		"Method": "GET",
	})
	if err == nil {
		t.Fatal("Should not have been enabled")
	}
	if err.Error() != "HTTP Executor is disabled. Please enable it in the configuration" {
		t.Fatal("Error message isn't correct")
	}
}
