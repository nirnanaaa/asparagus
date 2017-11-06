package slack

import (
	"bytes"
	"encoding/json"
	"html/template"
	"io"
	"net/http"
	"time"

	"github.com/nirnanaaa/asparagus/metric/adapters"
)

// Reporter type
type Reporter struct {
	Config Config
}

// NewReporter returns a new reporter
func NewReporter(config Config) *Reporter {
	return &Reporter{
		config,
	}
}

// StartReporting starts the recording process
func (r *Reporter) StartReporting() error {
	if !r.Config.Enabled {
		return nil
	}
	return nil
}

// StopReporting stops the recording process
func (r *Reporter) StopReporting() error {
	return nil
}

// SupportsLogging returns if the adapter supports logging
func (r *Reporter) SupportsLogging() bool {
	return r.Config.DirectLoggingEnabled
}

func (r *Reporter) parseTemplate(evt *adapters.LogEvent, w io.Writer) error {
	funcMap := template.FuncMap{
		"TimeFormat": func(f time.Time) string { return f.Format("02.01.2006T15:04:05") },
	}
	tmpl, err := template.New("templ").Funcs(funcMap).Parse(r.Config.LogFormat)
	if err != nil {
		return err
	}

	// Run the template to verify the output.
	err = tmpl.Execute(w, evt)
	if err != nil {
		return err
	}
	return nil
}

// LogResult logs a log line if SupportsLogging is enabled
func (r *Reporter) LogResult(evt *adapters.LogEvent) error {
	var reader io.Reader
	buf := new(bytes.Buffer)
	if err := r.parseTemplate(evt, buf); err != nil {
		return err
	}
	data := map[string]interface{}{
		"text": buf.String(),
	}
	payload, err := json.Marshal(&data)
	if err != nil {
		return err
	}
	reader = bytes.NewBuffer(payload)
	req, err := http.NewRequest("POST", r.Config.WebookURL, reader)
	if err != nil {
		return err
	}
	client := &http.Client{}

	if _, err := client.Do(req); err != nil {
		return err
	}
	return nil
}
