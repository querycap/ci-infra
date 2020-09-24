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

type JobOptionFunc func(s *WorkflowJob)

func Job(optFuncs ...JobOptionFunc) *WorkflowJob {
	s := &WorkflowJob{}

	for i := range optFuncs {
		optFuncs[i](s)
	}

	return s
}

func JobIf(cond string) JobOptionFunc {
	return func(s *WorkflowJob) {
		s.If = cond
	}
}

func JobNeeds(needs ...string) JobOptionFunc {
	return func(s *WorkflowJob) {
		if len(needs) > 0 {
			s.Needs = needs
		}
	}
}

func JobRunsOn(tags ...string) JobOptionFunc {
	return func(s *WorkflowJob) {
		if len(tags) == 0 {
			s.RunsOn = []string{"ubuntu-latest"}
			return
		}
		s.RunsOn = append([]string{"self-hosted"}, tags...)
	}
}

func JobSteps(steps ...*WorkflowStep) JobOptionFunc {
	return func(s *WorkflowJob) {
		s.Steps = steps
	}
}

func JobDefaultsWorkingDirectory(workingDirectory string) JobOptionFunc {
	return func(s *WorkflowJob) {
		if s.Defaults == nil {
			s.Defaults = &WorkflowJobDefaults{}
		}
		s.Defaults.Run.WorkingDirectory = workingDirectory
	}
}

type WorkflowJobDefaults struct {
	Run struct {
		Shell            string `yaml:"shell,omitempty"`
		WorkingDirectory string `yaml:"working-directory,omitempty"`
	} `yaml:"run,omitempty"`
}

type WorkflowJob struct {
	RunsOn   []string             `yaml:"runs-on,omitempty"`
	Needs    []string             `yaml:"needs,omitempty"`
	If       string               `yaml:"if,omitempty"`
	Defaults *WorkflowJobDefaults `yaml:"defaults,omitempty"`
	Steps    []*WorkflowStep      `steps:"steps,omitempty"`
}

type WorkflowStep struct {
	Uses string            `yaml:"uses,omitempty"`
	ID   string            `yaml:"id,omitempty"`
	Name string            `yaml:"name,omitempty"`
	If   string            `yaml:"if,omitempty"`
	With map[string]string `yaml:"with,omitempty"`
	Run  string            `yaml:"run,omitempty"`
}

type StepOptionFunc func(s *WorkflowStep)

func Step(optFuncs ...StepOptionFunc) *WorkflowStep {
	s := &WorkflowStep{}

	for i := range optFuncs {
		optFuncs[i](s)
	}

	return s
}

func StepUses(uses string) StepOptionFunc {
	return func(s *WorkflowStep) {
		s.Uses = uses
	}
}

func StepName(name string) StepOptionFunc {
	return func(s *WorkflowStep) {
		s.Name = name
	}
}

func StepIf(cond string) StepOptionFunc {
	return func(s *WorkflowStep) {
		s.If = cond
	}
}

func StepWith(values map[string]string) StepOptionFunc {
	return func(s *WorkflowStep) {
		s.With = values
	}
}

func StepRun(script string) StepOptionFunc {
	return func(s *WorkflowStep) {
		s.Run = strings.TrimSpace(script)
	}
}

func nameAndTagsFromDockerfile(dockerfile string) (string, []string) {
	parts := strings.Split(strings.Split(dockerfile, ".")[1], ",")
	if len(parts) > 1 {
		return parts[0], parts[1:]
	}
	return parts[0], []string{}
}
