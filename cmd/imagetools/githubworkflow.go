package main

import (
	"strings"
)

type Values = map[string]interface{}

type GithubWorkflow struct {
	Name string                  `yaml:"name"`
	On   Values                  `yaml:"on"`
	Jobs map[string]*WorkflowJob `yaml:"jobs"`
}

type WorkflowJob struct {
	RunsOn []string        `yaml:"runs-on,omitempty"`
	Needs  []string        `yaml:"needs,omitempty"`
	Steps  []*WorkflowStep `steps:"steps,omitempty"`
}

func Uses(uses string) *WorkflowStep {
	return &WorkflowStep{Uses: uses}
}

type WorkflowStep struct {
	Uses   string            `yaml:"uses,omitempty"`
	StepID string            `yaml:"id,omitempty"`
	Name   string            `yaml:"name,omitempty"`
	Cond   string            `yaml:"if,omitempty"`
	Args   map[string]string `yaml:"with,omitempty"`
	Run    string            `yaml:"run,omitempty"`
}

func (s WorkflowStep) Named(name string) *WorkflowStep {
	s.Name = name
	return &s
}

func (s WorkflowStep) ID(id string) *WorkflowStep {
	s.StepID = id
	return &s
}

func (s WorkflowStep) If(cond string) *WorkflowStep {
	s.Cond = cond
	return &s
}

func (s WorkflowStep) With(args map[string]string) *WorkflowStep {
	s.Args = args
	return &s
}

func (s WorkflowStep) Do(run string) *WorkflowStep {
	s.Run = strings.TrimSpace(run)
	return &s
}

func nameFromDockerfile(dockerfile string) string {
	parts := strings.Split(dockerfile, ".")
	return parts[1]
}
