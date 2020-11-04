package workflow

type WorkflowJob struct {
	If       string                     `yaml:"if,omitempty"`
	Needs    []string                   `yaml:"needs,omitempty"`
	Strategy *WorkflowJobStrategy       `yaml:"strategy,omitempty"`
	RunsOn   interface{}                `yaml:"runs-on,omitempty"`
	Services map[string]WorkflowService `yaml:"services,omitempty"`
	Defaults *WorkflowJobDefaults       `yaml:"defaults,omitempty"`
	Outputs  map[string]string          `yaml:"outputs,omitempty"`
	Steps    []*WorkflowStep            `yaml:"steps,omitempty"`
}

type WorkflowService struct {
	Image string   `yaml:"image,omitempty"`
	Ports []string `yaml:"ports,omitempty"`
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

func JobRunsOn(runsOn ...string) JobOptionFunc {
	return func(s *WorkflowJob) {
		if len(runsOn) == 0 {
			s.RunsOn = "ubuntu-latest"
			return
		}
		if len(runsOn) == 1 {
			s.RunsOn = runsOn[0]
			return
		}
		s.RunsOn = runsOn
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

func JobService(name string, svc WorkflowService) JobOptionFunc {
	return func(s *WorkflowJob) {
		if s.Services == nil {
			s.Services = map[string]WorkflowService{}
		}
		s.Services[name] = svc
	}
}

func JobOutputs(outputs map[string]string) JobOptionFunc {
	return func(s *WorkflowJob) {
		s.Outputs = outputs
	}
}

func JobNeeds(needs ...string) JobOptionFunc {
	return func(s *WorkflowJob) {
		s.Needs = needs
	}
}

func JobIf(cond string) JobOptionFunc {
	return func(s *WorkflowJob) {
		s.If = cond
	}
}
