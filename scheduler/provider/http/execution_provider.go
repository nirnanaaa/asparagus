package http

import (
	"fmt"

	"github.com/Sirupsen/logrus"
	"github.com/nirnanaaa/asparagus/reflection"
	"github.com/nirnanaaa/asparagus/scheduler/provider"
)

// ExecutionProvider fills in the interface
type ExecutionProvider struct {
	Logger   *logrus.Logger
	Executor Executor
	Config   Config
}

// NewExecutionProvider creates a new ExecutionProvider
func NewExecutionProvider(conf Config, logger *logrus.Logger) *ExecutionProvider {
	exec := NewExecutor(logger)
	exec.HTTPConfig = conf
	return &ExecutionProvider{
		logger,
		exec,
		conf,
	}
}

// ExecutionData is used to determine settings for the request
type ExecutionData struct {
	URL    string `json:"url"`
	Method string `json:"method"`
}

// ParseExecutionContext parses the message
func (p *ExecutionProvider) extractData(data *ExecutionData, msg interface{}) error {
	switch v := msg.(type) {
	case map[string]interface{}:
		iface, ok := msg.(map[string]interface{})
		if !ok {
			return fmt.Errorf("interface conversion failed")
		}
		if err := reflection.MapToStruct(data, iface); err != nil {
			return err
		}
	default:
		return fmt.Errorf("unknown input type on executor: %T", v)
	}
	return nil
}

// Execute runs a task
func (p *ExecutionProvider) Execute(t interface{}, r *provider.Response) error {
	if !p.Config.Enabled {
		return fmt.Errorf("HTTP Executor is disabled. Please enable it in the configuration")
	}
	var task ExecutionData
	if err := p.extractData(&task, t); err != nil {
		return err
	}
	return p.Executor.FromTask(task, r)
}
