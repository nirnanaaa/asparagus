package provider_test

import (
	"fmt"
	"testing"

	"github.com/nirnanaaa/asparagus/scheduler/provider"
)

func TestProvider_Backoff_Min(t *testing.T) {
	backoff := provider.CalculateBackoffForAttempt(0)

	if backoff.String() != "100ms" {
		fmt.Printf("Backoff min should be 100ms, got: %s", backoff.String())
	}
}

func TestProvider_Backoff_ThreeTries(t *testing.T) {
	backoff := provider.CalculateBackoffForAttempt(3)

	if backoff.String() != "2.545323628s" {
		fmt.Printf("Backoff min should be 2.545323628s, got: %s", backoff.String())
	}
}

func TestProvider_Backoff_TenTries(t *testing.T) {
	backoff := provider.CalculateBackoffForAttempt(10)

	if backoff.String() != "1h5m24.194202244s" {
		fmt.Printf("Backoff min should be 1h5m24.194202244s, got: %s", backoff.String())
	}
}
func TestProvider_Backoff_Max(t *testing.T) {
	backoff := provider.CalculateBackoffForAttempt(100)

	if backoff.String() != "5h0m0s" {
		fmt.Printf("Backoff min should be 5h0m0s, got: %s", backoff.String())
	}
}
