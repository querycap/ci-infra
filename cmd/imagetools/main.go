package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

var (
	hub              = os.Getenv("HUB")
	basePathForBuild = "build"
)

func init() {
	if hub == "" {
		panic(errors.New("missing ${HUB}"))
	}
}

func main() {
	projects, err := resolveProjects()
	if err != nil {
		panic(err)
	}

	data, _ := yaml.Marshal(projects)
	fmt.Println(string(data))

	generateWorkflows(projects)
	generateWorkflowsForSync(projects)
	generateDependabot(projects)
}

func generateWorkflows(projects Projects) {
	projects.Range(func(p *Project) {
		w := githubWorkflowFromProject(p)
		writeWorkflow(w)
	})
}

func githubWorkflowFromProject(p *Project) *GithubWorkflow {
	w := &GithubWorkflow{}
	w.Name = p.Name

	w.On = Values{
		"push": Values{
			"paths": []string{
				filepath.Join(basePathForBuild, p.Name, "**"),
			},
		},
	}

	w.Jobs = map[string]*WorkflowJob{}

	for i := range p.Dockerfiles {
		name := nameFromDockerfile(p.Dockerfiles[i])
		w.Jobs[name] = jobDockerBuild(p.Name, name)
	}

	return w
}

func jobDockerBuild(projectName string, name string, needs ...string) *WorkflowJob {
	return &WorkflowJob{
		RunsOn: []string{"ubuntu-latest"},
		Needs:  needs,
		Steps: []*WorkflowStep{
			Uses("actions/checkout@v2"),
			Uses("docker/setup-qemu-action@v1"),
			Uses("docker/setup-buildx-action@v1"),
			Uses("docker/login-action@v1").With(map[string]string{
				"username": "${{ secrets.DOCKER_USERNAME }}",
				"password": "${{ secrets.DOCKER_PASSWORD }}",
			}),
			Uses("").If("github.ref == 'refs/heads/master'").Named("Versioned Build").Do(fmt.Sprintf("cd %s/%s && make build HUB=%s NAME=%s", basePathForBuild, projectName, hub, name)),
			Uses("").If("github.ref != 'refs/heads/master'").Named("Temp Build").Do(fmt.Sprintf("cd %s/%s && make build HUB=%s NAME=%s TAG=${{ github.sha }}", basePathForBuild, projectName, hub, name)),
		},
	}
}

func generateWorkflowsForSync(projects Projects) {
	projects.Range(func(p *Project) {
		for i := range p.Dockerfiles {
			name := nameFromDockerfile(p.Dockerfiles[i])
			dockerfile := fmt.Sprintf("sync/Dockerfile.zz_%s", name)
			_ = ioutil.WriteFile(dockerfile, []byte(fmt.Sprintf("FROM "+hub+"/%s:%s", name, p.Version)), os.ModePerm)
		}
	})

	files, _ := filepath.Glob("sync/Dockerfile.*")

	for i := range files {
		writeWorkflow(githubWorkflowForSync(nameFromDockerfile(files[i]), files[i]))
	}
}

func generateDependabot(projects Projects) {
	buf := bytes.NewBufferString(`
version: 2
updates:
  - package-ecosystem: "github-actions"
    directory: "/"
    schedule:
      interval: "daily"

  - package-ecosystem: "docker"
    directory: "/sync"
    schedule:
      interval: "daily"
`)

	projects.Range(func(p *Project) {
		_, _ = io.WriteString(buf, fmt.Sprintf(`
  - package-ecosystem: "docker"
    directory: "/build/%s"
    schedule:
      interval: "daily"
`, p.Name))
	})

	_ = ioutil.WriteFile(".github/dependabot.yml", buf.Bytes(), os.ModePerm)
}

func writeWorkflow(w *GithubWorkflow) {
	if w == nil {
		return
	}

	data, _ := yaml.Marshal(w)
	_ = ioutil.WriteFile(fmt.Sprintf(".github/workflows/%s.yml", w.Name), data, os.ModePerm)
}

func githubWorkflowForSync(name string, dockerfile string) *GithubWorkflow {
	w := &GithubWorkflow{}
	w.Name = "sync-" + name
	w.On = Values{
		"push": Values{
			"paths": []string{dockerfile},
		},
	}
	w.Jobs = map[string]*WorkflowJob{
		"sync": {
			RunsOn: []string{"self-hosted", "ARM64"},
			Steps: []*WorkflowStep{
				Uses("actions/checkout@v2"),
				Uses("docker/setup-qemu-action@v1"),
				Uses("docker/setup-buildx-action@v1"),
				Uses("docker/login-action@v1").With(map[string]string{
					"registry": "${{ secrets.DOCKER_MIRROR_REGISTRY }}",
					"username": "${{ secrets.DOCKER_MIRROR_USERNAME }}",
					"password": "${{ secrets.DOCKER_MIRROR_PASSWORD }}",
				}),
				Uses("").Do(fmt.Sprintf(`cd sync && make sync HUB=${{ secrets.DOCKER_MIRROR_REGISTRY }} NAME=%s`, name)),
			},
		},
	}

	return w
}
