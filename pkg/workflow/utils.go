package workflow

import (
	"strings"
)

func OutputFromStep(stepName string) func(arg string) string {
	return func(arg string) string {
		return Ref("steps", stepName, "outputs", arg)
	}
}

func Ref(paths ...string) string {
	return "${{ " + strings.Join(paths, ".") + " }}"
}
