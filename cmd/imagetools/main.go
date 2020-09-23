package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
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

func cleanup() error {
	files, err := glob(".github/workflows/*", "sync/Dockerfile.*")
	if err != nil {
		return err
	}

	for i := range files {
		if err := os.Remove(files[i]); err != nil {
			return err
		}
	}

	return nil
}

func main() {
	initMakefiles()

	if err := cleanup(); err != nil {
		panic(err)
	}

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
		for i := range p.Dockerfiles {
			name, tags := nameAndTagsFromDockerfile(p.Dockerfiles[i])

			w := &GithubWorkflow{}

			if p.Name == name {
				w.Name = p.Name
			} else {
				w.Name = p.Name + "-" + name
			}

			w.On = Values{
				"push": Values{
					"paths": []string{
						p.Dockerfiles[i],
						p.VersionFile,
					},
				},
			}

			w.Jobs = map[string]*WorkflowJob{}

			w.Jobs[name] = &WorkflowJob{
				RunsOn: runsOn(tags),
				Steps: []*WorkflowStep{
					Uses("actions/checkout@v2"),
					Uses("docker/setup-qemu-action@v1"),
					Uses("docker/setup-buildx-action@v1"),
					Uses("docker/login-action@v1").With(map[string]string{
						"username": "${{ secrets.DOCKER_USERNAME }}",
						"password": "${{ secrets.DOCKER_PASSWORD }}",
					}),
					Uses("").If("github.ref == 'refs/heads/master'").Named("Versioned Build").Do(fmt.Sprintf("cd %s/%s && make build HUB=%s NAME=%s", basePathForBuild, p.Name, hub, fullname(name, tags))),
					Uses("").If("github.ref != 'refs/heads/master'").Named("Temp Build").Do(fmt.Sprintf("cd %s/%s && make build HUB=%s NAME=%s TAG=temp-${{ github.sha }}", basePathForBuild, p.Name, hub, fullname(name, tags))),
				},
			}

			writeWorkflow(w)
		}
	})
}

func generateWorkflowsForSync(projects Projects) {
	projects.Range(func(p *Project) {
		for i := range p.Dockerfiles {
			name, _ := nameAndTagsFromDockerfile(p.Dockerfiles[i])
			dockerfile := fmt.Sprintf("sync/Dockerfile.%s,arm64", name)
			_ = generateFile(dockerfile, []byte(fmt.Sprintf("FROM "+hub+"/%s:%s", name, p.Version)))
		}
	})

	files, _ := filepath.Glob("sync/Dockerfile.*")

	for i := range files {
		name, tags := nameAndTagsFromDockerfile(files[i])
		writeWorkflow(githubWorkflowForSync(name, files[i], tags...))
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

	_ = generateFile(".github/dependabot.yml", buf.Bytes())
}

func writeWorkflow(w *GithubWorkflow) {
	if w == nil {
		return
	}
	data, _ := yaml.Marshal(w)
	_ = generateFile(workflowName(w.Name), data)
}

func workflowName(name string) string {
	return fmt.Sprintf(".github/workflows/%s.yml", name)
}

func githubWorkflowForSync(name string, dockerfile string, tags ...string) *GithubWorkflow {
	w := &GithubWorkflow{}
	w.Name = "zz-sync-" + name
	w.On = Values{
		"push": Values{
			"paths": []string{
				workflowName(w.Name),
				dockerfile,
			},
		},
	}

	w.Jobs = map[string]*WorkflowJob{
		"sync": {
			RunsOn: runsOn(tags),
			Steps: []*WorkflowStep{
				Uses("actions/checkout@v2"),
				Uses("docker/setup-qemu-action@v1"),
				Uses("docker/setup-buildx-action@v1"),
				Uses("docker/login-action@v1").With(map[string]string{
					"registry": "${{ secrets.DOCKER_MIRROR_REGISTRY }}",
					"username": "${{ secrets.DOCKER_MIRROR_USERNAME }}",
					"password": "${{ secrets.DOCKER_MIRROR_PASSWORD }}",
				}),
				Uses("").Do(`
DOCKERFILE=sync/Dockerfile.` + fullname(name, tags) + `
TAG=$(cat ${DOCKERFILE} | grep "^FROM " | sed -e "s/FROM //g" | head -1)
make -f ./build/Makefile buildx TAGS=${{ secrets.DOCKER_MIRROR_REGISTRY }}/${TAG} DOCKERFILE=${DOCKERFILE}
`),
			},
		},
	}

	return w
}

func runsOn(tags []string) []string {
	if len(tags) == 0 {
		return []string{"ubuntu-latest"}
	}
	return append([]string{"self-hosted"}, tags...)
}

func fullname(name string, flags []string) string {
	b := bytes.NewBufferString(name)

	for i := range flags {
		b.WriteByte(',')
		b.WriteString(flags[i])
	}

	return b.String()
}
