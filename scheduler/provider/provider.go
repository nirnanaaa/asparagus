package provider

import (
	"math"
	"math/rand"
	"time"
)

// SourceProvider is used to get tasks
type SourceProvider interface {
	Read() error
	OnTaskUpdate(func(*Task) error)
	Stop() error
	TaskStarted(*Task) error
	TaskDone(*Task) error
	TaskError(*Task, error) error
}

// ExecutionProvider is used for executing cronjobs
type ExecutionProvider interface {
	Execute(interface{}, *Response) error
}

// Task represents a task to be run
type Task struct {
	Name              string         `json:"name"`
	Running           bool           `json:"is_running,omitempty"`
	Expression        string         `json:"schedule,omitempty"`
	LastRunAt         time.Time      `json:"last_run_at,omitempty"`
	AfterTask         string         `json:"after_task,omitempty"`
	Key               string         `json:"-"`
	SourceProvider    SourceProvider `json:"-"`
	Executor          string         `json:"executor"`
	ExecutionConfig   interface{}    `json:"execution_config"`
	LastError         string         `json:"last_error"`
	CurrentRetryCount int            `json:"current_retry_counter"`
	MaxRetries        int            `json:"max_retries"`
}

// Response contains data that executors can return
type Response struct {
	StatusCode int    `json:"status_code"`
	Response   []byte `json:"response"`
}

const maxInt64 = float64(math.MaxInt64 - 512)

// CalculateBackoffForAttempt returns the duration for a specific attempt. This is useful if
// you have a large number of independent Backoffs, but don't want use
// unnecessary memory storing the Backoff parameters per Backoff. The first
// attempt should be 0.
//
// ForAttempt is concurrent-safe.
func CalculateBackoffForAttempt(attempt float64) time.Duration {
	// Zero-values are nonsensical, so we use
	// them to apply defaults
	min := 100 * time.Millisecond
	max := 5 * time.Hour
	factor := 3.0
	//calculate this duration
	minf := float64(min)
	durc := minf * math.Pow(factor, attempt)
	durf := rand.Float64()*(durc-minf) + minf
	//ensure float64 wont overflow int64
	if durf > maxInt64 {
		return max
	}
	dur := time.Duration(durf)
	//keep within bounds
	if dur < min {
		return min
	} else if dur > max {
		return max
	}
	return dur
}
