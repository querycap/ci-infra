package main

import (
	"io/ioutil"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

func glob(patterns ...string) (files []string, err error) {
	for i := range patterns {
		matched, err := filepath.Glob(patterns[i])
		if err != nil {
			return nil, err
		}
		files = append(files, matched...)
	}
	return
}

func resolveProjects() (Projects, error) {
	dockerfiles, err := glob("./build/*/Dockerfile.*", "./build/*/.version", "./build/*/Makefile")
	if err != nil {
		return nil, err
	}

	projects := Projects{}

	for i := range dockerfiles {
		parts := strings.Split(dockerfiles[i], "/")
		projectName, dockerfileName := parts[1], parts[2]

		p, ok := projects[projectName]
		if !ok {
			p = &Project{
				Name: projectName,
			}
			projects[projectName] = p
		}

		switch dockerfileName {
		case "Makefile":
			p.Makefile = dockerfiles[i]
		case ".version", "Dockerfile.version":
			data, _ := ioutil.ReadFile(dockerfiles[i])
			p.Version = getVersionFromDockerfileVersionOrDotVersion(data)
			p.VersionFile = dockerfiles[i]
		default:
			p.Dockerfiles = append(p.Dockerfiles, dockerfiles[i])
		}

	}

	return projects, nil
}

type Projects map[string]*Project

func (projects Projects) Range(fn func(p *Project)) {
	names := make([]string, 0)

	for i := range projects {
		names = append(names, projects[i].Name)
	}

	sort.Strings(names)

	for _, n := range names {
		fn(projects[n])
	}
}

type Project struct {
	Name        string
	Version     string
	VersionFile string
	Makefile    string
	Dockerfiles []string
}

var reVersionPrefix = regexp.MustCompile("@opt:prefix +(.+)")
var reVersionInDockerFrom = regexp.MustCompile("FROM.+:(.+)")

func getVersionFromDockerfileVersionOrDotVersion(data []byte) string {
	prefix := "v"

	v := string(data)

	if strings.Contains(v, "FROM") {
		for _, l := range strings.Split(v, "\n") {
			if matched := reVersionPrefix.FindAllStringSubmatch(l, 1); matched != nil {
				prefix = matched[0][1]
			}
			if matched := reVersionInDockerFrom.FindAllStringSubmatch(l, 1); matched != nil {
				v = matched[0][1]
				break
			}
		}

	}
	return strings.TrimLeft(v, prefix)
}
