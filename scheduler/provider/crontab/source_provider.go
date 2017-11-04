package crontab

import (
	"bufio"
	"fmt"
	"os"
	"regexp"

	"github.com/nirnanaaa/asparagus/reflection"
	"github.com/nirnanaaa/asparagus/scheduler/provider"
)

// SourceProvider creates a new src provider
type SourceProvider struct {
	Config Config
}

// NewSourceProvider creates a new src provider
func NewSourceProvider(config Config) *SourceProvider {
	return &SourceProvider{
		config,
	}
}

// Read reads from FS
func (p *SourceProvider) Read(tasks *map[string]*provider.Task) error {
	if !p.Config.Enabled {
		return nil
	}
	lines, err := p.openFile()
	if err != nil {
		return err
	}
	for _, line := range *lines {
		px := provider.Task{}
		if err := p.parseLine(&px, line); err != nil {
			return err
		}
		// TODO: This clearly sucks
		task := *tasks
		task[px.Name] = &px
		tasks = &task
	}
	return nil
}

func (p *SourceProvider) parseLine(t *provider.Task, line string) error {
	rex := regexp.MustCompile(`(?P<schedule>(.*?))\s+(?P<config>\[(.*?)\])\s+(?P<details>\[(.*?)\])`)
	if !rex.MatchString(line) {
		return fmt.Errorf("Regex error on line: \"%s\"", line)
	}
	t.Expression = rex.ReplaceAllString(line, "${schedule}")
	config := rex.ReplaceAllString(line, fmt.Sprintf("${config}"))
	executorConf := rex.ReplaceAllString(line, fmt.Sprintf("${details}"))
	confData := reflection.ExprToMap(config)
	if err := reflection.MapToStruct(t, confData); err != nil {
		return err
	}
	executorConfData := reflection.ExprToMap(executorConf)
	t.ExecutionConfig = executorConfData
	return nil
}

var commentRex = regexp.MustCompile(`^(;|#).*$`)

func (p *SourceProvider) openFile() (*[]string, error) {
	lines := []string{}
	file, err := os.Open(p.Config.Crontab)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		txt := scanner.Text()
		if !commentRex.MatchString(txt) {
			lines = append(lines, txt)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return &lines, nil
}
