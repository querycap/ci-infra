package workflow

import "strings"

type WorkflowStep struct {
	Uses string            `yaml:"uses,omitempty"`
	ID   string            `yaml:"id,omitempty"`
	Name string            `yaml:"name,omitempty"`
	If   string            `yaml:"if,omitempty"`
	With map[string]string `yaml:"with,omitempty"`
	Env  map[string]string `yaml:"env,omitempty"`
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

func StepEnv(env map[string]string, envExtra ...map[string]string) StepOptionFunc {
	return func(s *WorkflowStep) {
		f := map[string]string{}

		for _, e := range append([]map[string]string{env}, envExtra...) {
			for k, v := range e {
				f[k] = v
			}
		}

		s.Env = f
	}
}

func StepID(id string) StepOptionFunc {
	return func(s *WorkflowStep) {
		s.ID = id
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
