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

func initAll() {
	_ = generateFile("./build/Dockerfile.sync", []byte(`
ARG TAG
FROM ${TAG}
`))

	_ = generateFile("./build/Makefile", []byte(`
TAGS = latest    master
BUILD_ARGS = VERSION=0.0.0 GOPROXY=https://goproxy.com,direct
TARGET_PLATFORM = linux/arm64,linux/amd64
DOCKERFILE = Dockerfile
CONTEXT = .

TAG_FLAGS = $(foreach v,$(TAGS),$(shell echo "--tag $(v)"))
BUILD_ARG_FLAGS = $(foreach v,$(BUILD_ARGS),$(shell echo "--build-arg $(v)"))

buildx:
	docker buildx build \
		--push \
		--cache-from type=local,src=/tmp/cache \
		--cache-to type=local,dest=/tmp/cache \
		--platform $(TARGET_PLATFORM) \
		$(BUILD_ARG_FLAGS) \
		$(TAG_FLAGS) \
		--file $(DOCKERFILE) $(CONTEXT)
`))
}

func main() {
	initAll()

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
						workflowFilename(w.Name),
						p.Dockerfiles[i],
						p.VersionFile,
						p.Makefile,
					},
				},
			}

			dockerfileName := fullname(name, tags)
			workingDir := filepath.Join(basePathForBuild, p.Name)

			w.Jobs = map[string]*WorkflowJob{
				name: Job(
					JobDefaultsWorkingDirectory(workingDir),
					JobRunsOn(tags...),
					JobSteps(
						Step(StepUses("actions/checkout@v2")),
						Step(StepUses("actions/cache@v2"), StepWith(map[string]string{
							"key":  "${{ runner.os }}-" + name,
							"path": "/tmp/cache",
						})),
						Step(StepUses("docker/setup-qemu-action@v1")),
						Step(StepUses("docker/setup-buildx-action@v1")),
						Step(StepUses("docker/login-action@v1"), StepWith(map[string]string{
							"username": "${{ secrets.DOCKER_USERNAME }}",
							"password": "${{ secrets.DOCKER_PASSWORD }}",
						})),
						Step(
							StepName("Versioned Build"),
							StepIf("github.ref == 'refs/heads/master'"),
							StepRun(fmt.Sprintf(`
make build HUB=%s NAME=%s
`, hub, dockerfileName))),
						Step(
							StepName("Temp Build"),
							StepIf("github.ref != 'refs/heads/master'"),
							StepRun(fmt.Sprintf(`
make build TAG=temp-${{ github.sha }} HUB=%s NAME=%s
`, hub, dockerfileName)),
						),
					),
				),
				"sync-" + name: Job(
					JobDefaultsWorkingDirectory(workingDir),
					JobRunsOn("arm64"),
					JobNeeds(name),
					JobIf("github.ref == 'refs/heads/master'"),
					JobSteps(
						Step(StepUses("actions/checkout@v2")),
						Step(StepUses("docker/setup-qemu-action@v1")),
						Step(StepUses("docker/setup-buildx-action@v1")),
						Step(StepUses("docker/login-action@v1"), StepWith(map[string]string{
							"registry": "${{ secrets.DOCKER_MIRROR_REGISTRY }}",
							"username": "${{ secrets.DOCKER_MIRROR_USERNAME }}",
							"password": "${{ secrets.DOCKER_MIRROR_PASSWORD }}",
						})),
						Step(StepRun(fmt.Sprintf(`
export TAG=$(make image HUB=%s NAME=%s)

docker buildx build --push --platform linux/arm64,linux/amd64 --tag ${{ secrets.DOCKER_MIRROR_REGISTRY }}/${TAG} --build-arg TAG=${TAG} --file ../Dockerfile.sync .
`, hub, dockerfileName))),
					),
				),
			}

			writeWorkflow(w)
		}
	})
}

func generateDependabot(projects Projects) {
	buf := bytes.NewBufferString(`
version: 2
updates:
  - package-ecosystem: "github-actions"
    directory: "/"
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
	_ = generateFile(workflowFilename(w.Name), data)
}

func workflowFilename(name string) string {
	return fmt.Sprintf(".github/workflows/%s.yml", name)
}

func fullname(name string, flags []string) string {
	b := bytes.NewBufferString(name)

	for i := range flags {
		b.WriteByte(',')
		b.WriteString(flags[i])
	}

	return b.String()
}
