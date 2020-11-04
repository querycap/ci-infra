package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	. "github.com/querycap/ci-infra/pkg/workflow"
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
	files, err := Glob(".github/workflows/zz-*")
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
	if err := cleanup(); err != nil {
		panic(err)
	}

	projects, err := ResolveProjects()
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
			name := nameDockerfile(p.Dockerfiles[i])

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

			jobs := map[string]*WorkflowJob{}

			steps := []*WorkflowStep{
				Step(StepUses("actions/checkout@v2")),
			}

			runsOn := p.Workflow.RunsOn

			platforms := p.Workflow.Platforms

			if len(platforms) == 0 {
				platforms = []string{"linux/amd64", "linux/arm64"}
			}

			if p.Workflow.QEMU == nil {
				enabled := true
				p.Workflow.QEMU = &enabled
			}

			env := map[string]string{}

			matrixArch := map[string][]string{}

			if len(p.Workflow.Schedule) > 0 {
				w.On["schedule"] = p.Workflow.Schedule
			}

			workingDir := filepath.Join(basePathForBuild, p.Name)

			dockerSetupSteps, dockerPushStep := resolveDockerSteps(workingDir, name)

			if *p.Workflow.QEMU {
				steps = append(steps, Step(StepUses("docker/setup-qemu-action@v1")))

				if len(runsOn) == 0 {
					runsOn = []string{
						"ubuntu-latest",
					}
				}
			} else {
				matrixArch["arch"] = toArchs(platforms)

				if len(runsOn) == 0 {
					runsOn = []string{
						Ref(`matrix.arch != 'amd64' && fromJSON(format('["self-hosted","linux","{0}"]', matrix.arch)) || 'ubuntu-latest'`),
					}
				} else {
					runsOn = append(runsOn, "linux", "${{ matrix.arch }}")
				}

				combineSteps := append([]*WorkflowStep{}, dockerSetupSteps...)

				combineSteps = append(combineSteps, Step(
					StepName("Combine"),
					StepEnv(env, map[string]string{
						"IMAGE":       Ref("needs", name, "outputs", "image"),
						"HUB":         hub,
						"TARGET_ARCH": strings.Join(toArchs(platforms), " "),
					}),
					StepRun(`
for h in ${HUB}; do
  SOURCES=""
  for arch in ${TARGET_ARCH}; do
    SOURCES="${SOURCES} ${h}/${IMAGE}-${arch}"
  done

  docker buildx imagetools create -t ${h}/${IMAGE} ${SOURCES};
  docker buildx imagetools inspect ${h}/${IMAGE};
done
`)))

				jobs[name+"-combine"] = Job(
					JobIf("${{ github.event_name != 'pull_request' }}"),
					JobStrategyMatrix(p.Workflow.Matrix),
					JobNeeds(name),
					JobRunsOn("ubuntu-latest"),
					JobSteps(combineSteps...),
				)
			}

			if len(p.Workflow.Matrix) > 0 {
				for k := range p.Workflow.Matrix {
					env[k] = Ref("matrix", k)
				}
			}

			steps = append(steps, dockerSetupSteps...)

			steps = append(steps,
				Step(
					StepName("Prepare"),
					StepID("prepare"),
					StepEnv(env, map[string]string{
						"GITHUB_SHA": Ref("github", "ref"),
						"GITHUB_REF": Ref("github", "sha"),
						"NAME":       name,
					}),
					StepRun(setOutputTargetPlatforms(platforms, *p.Workflow.QEMU)+`
if [[ ${GITHUB_REF} != "refs/heads/master" ]]; then
  export TAG=sha-${GITHUB_SHA::7}
fi

make prepare
`),
				),
				dockerPushStep,
			)

			jobs[name] = Job(
				JobStrategyMatrix(mergeMatrix(p.Workflow.Matrix, matrixArch)),
				JobOutputs(map[string]string{"image": stepPrepareOutput("image")}),
				JobDefaultsWorkingDirectory(workingDir),
				JobRunsOn(runsOn...),
				JobSteps(steps...),
			)

			w.Jobs = jobs

			writeWorkflow(w)
		}
	})
}

func mergeMatrix(matrixes ...map[string][]string) map[string][]string {
	m := map[string][]string{}
	for _, mat := range matrixes {
		for k := range mat {
			m[k] = mat[k]
		}
	}
	return m
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
	return fmt.Sprintf(".github/workflows/zz-%s.yml", name)
}

func nameDockerfile(dockerfile string) string {
	return strings.Split(dockerfile, ".")[1]
}

func generateFile(filename string, data []byte) error {
	data = append(bytes.TrimSpace(data), '\n')
	if err := os.MkdirAll(filepath.Dir(filename), os.ModePerm); err != nil {
		return err
	}
	return ioutil.WriteFile(filename, data, os.ModePerm)
}

var stepPrepareOutput = OutputFromStep("prepare")

func resolveDockerSteps(workingDir string, name string) ([]*WorkflowStep, *WorkflowStep) {
	imageTags := make([]string, 0)
	steps := []*WorkflowStep{
		Step(StepUses("docker/setup-buildx-action@v1"), StepWith(map[string]string{
			"driver-opts": "network=host",
		})),
	}

	for _, h := range strings.Split(hub, " ") {
		if h != "" {
			imageTags = append(imageTags, fmt.Sprintf(`%s/%s%s`, h, stepPrepareOutput("image"), stepPrepareOutput("image_suffix")))

			registry := strings.Split(h, "/")[0]

			hubLogin := map[string]string{
				"registry": registry,
			}

			name := strings.ToUpper(strings.Split(registry, ".")[0])

			switch registry {
			case "ghcr.io":
				hubLogin["username"] = "${{ github.repository_owner }}"
				hubLogin["password"] = "${{ secrets.CR_PAT }}"
			default:
				hubLogin["username"] = "${{ secrets." + name + "_USERNAME }}"
				hubLogin["password"] = "${{ secrets." + name + "_PASSWORD }}"
			}

			steps = append(steps, Step(StepName("Login "+registry), StepUses("docker/login-action@v1"), StepWith(hubLogin)))
		}
	}

	return steps, Step(
		StepName("Build & May Push"),
		StepUses("docker/build-push-action@v2"),
		StepWith(map[string]string{
			"context":    workingDir,
			"file":       workingDir + "/Dockerfile." + name,
			"push":       "${{ github.event_name != 'pull_request' }}",
			"build-args": stepPrepareOutput("build_args"),
			"labels": strings.Join([]string{
				"org.opencontainers.image.source=https://github.com/${{ github.repository }}",
				"org.opencontainers.image.revision=${{ github.sha }}",
			}, "\n"),
			"platforms": stepPrepareOutput("target_platforms"),
			"tags":      strings.Join(imageTags, "\n"),
		}),
	)
}

func setOutputTargetPlatforms(platforms []string, qemu bool) string {
	if qemu {
		return fmt.Sprintf(`
echo ::set-output name=target_platforms::%s
`, strings.Join(platforms, ","))
	}
	return `
echo ::set-output name=image_suffix::-${{ matrix.arch }}
echo ::set-output name=target_platforms::linux/${{ matrix.arch }}
`
}

func toArchs(platforms []string) (archs []string) {
	for _, arch := range platforms {
		archs = append(archs, strings.Split(arch, "/")[1])
	}
	return
}
