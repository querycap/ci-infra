package workflow

type WorkflowJob struct {
	RunsOn   []string             `yaml:"runs-on,omitempty"`
	Needs    []string             `yaml:"needs,omitempty"`
	Defaults *WorkflowJobDefaults `yaml:"defaults,omitempty"`
	Strategy *WorkflowJobStrategy `yaml:"strategy,omitempty"`
	Steps    []*WorkflowStep      `steps:"steps,omitempty"`
}

type WorkflowJobDefaults struct {
	Run struct {
		Shell            string `yaml:"shell,omitempty"`
		WorkingDirectory string `yaml:"working-directory,omitempty"`
	} `yaml:"run,omitempty"`
}

type WorkflowJobStrategy struct {
	Matrix map[string][]string `yaml:"matrix,omitempty"`
}

type JobOptionFunc func(s *WorkflowJob)

func Job(optFuncs ...JobOptionFunc) *WorkflowJob {
	s := &WorkflowJob{}

	for i := range optFuncs {
		optFuncs[i](s)
	}

	return s
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

func JobStrategyMatrix(matrix map[string][]string) JobOptionFunc {
	return func(s *WorkflowJob) {
		if len(matrix) == 0 {
			return
		}

		if s.Strategy == nil {
			s.Strategy = &WorkflowJobStrategy{}
		}
		s.Strategy.Matrix = matrix
	}
}
