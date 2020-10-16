package workflow

type Values = map[string]interface{}

type GithubWorkflow struct {
	Name string                  `yaml:"name"`
	On   Values                  `yaml:"on"`
	Jobs map[string]*WorkflowJob `yaml:"jobs"`
}
