package utils

import (
	"bufio"
	"elbe-prj/containers"
	"regexp"
	"strings"
)

var (
	BuildResultMap = map[string]containers.BuildResult{
		"build_failed":  containers.Build_Failed,
		"build_done":    containers.Build_Done,
		"empty_project": containers.Empty_Project,
		"busy":          containers.Busy,
		"needs_build":   containers.Needs_Build,
	}
)

func SplitLines(s string) []string {
	var lines []string
	sc := bufio.NewScanner(strings.NewReader(s))
	for sc.Scan() {
		lines = append(lines, sc.Text())
	}
	return lines
}

func ParseLine(s string) containers.Project {
	words := strings.Fields(s)
	var keys = []string{"build_failed", "build_done", "empty_project", "busy", "needs_build"}
	var c = containers.Needs_Build

	for _, v := range keys {
		matched, _ := regexp.MatchString(v, s)
		if matched {
			c, _ = BuildResultMap[strings.ToLower(v)]
		}
	}
	return containers.Project{Path: words[0], Name: words[1], Result: c}
}
