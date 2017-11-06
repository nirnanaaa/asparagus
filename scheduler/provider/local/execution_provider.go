package local

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/nirnanaaa/asparagus/reflection"
	"github.com/nirnanaaa/asparagus/scheduler/provider"
)

// ExecutionProvider fills in the interface
type ExecutionProvider struct {
	Logger *logrus.Logger
	Config Config
}

// NewExecutionProvider creates a new ExecutionProvider
func NewExecutionProvider(conf Config, logger *logrus.Logger) *ExecutionProvider {
	return &ExecutionProvider{
		logger,
		conf,
	}
}

// ExecutionData is used to determine settings for the request
type ExecutionData struct {
	Command   string
	Arguments []string
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
	case string:
		mesg := msg.(string)
		parts := strings.Split(mesg, " ")
		if mesg == "" || len(parts) < 1 {
			return fmt.Errorf("command has no parts")
		}
		data.Command = parts[0]
		if len(parts) > 1 {
			data.Arguments = parts[1:]
		}
	default:
		return fmt.Errorf("unknown input type on executor: %T", v)
	}
	return nil
}

// Execute runs a task
func (p *ExecutionProvider) Execute(t interface{}, r *provider.Response) error {
	if !p.Config.Enabled {
		return fmt.Errorf("Local Executor is disabled. Please enable it in the configuration")
	}
	var task ExecutionData
	if err := p.extractData(&task, t); err != nil {
		return err
	}

	path, err := exec.LookPath(task.Command)
	if err != nil {
		return err
	}
	cmd := exec.Command(path, task.Arguments...)
	if p.Config.EnabledOutput {
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	}

	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
}
