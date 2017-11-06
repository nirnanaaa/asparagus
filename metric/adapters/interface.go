package adapters

import "time"

// Reporter starts a new reporting agent
type Reporter interface {
	StartReporting() error
	StopReporting() error
	SupportsLogging() bool

	LogResult(l *LogEvent) error
}

// LogEvent defines an event used for logging purposes
type LogEvent struct {
	Name        string    `json:"name"`
	StartDate   time.Time `json:"start_date"`
	EndDate     time.Time `json:"end_date"`
	Took        string    `json:"took"`
	Result      string    `json:"result"`
	StatusCode  int       `json:"status_code"`
	IsError     bool      `json:"is_error"`
	ResultDebug string    `json:"result_debug"`
}
