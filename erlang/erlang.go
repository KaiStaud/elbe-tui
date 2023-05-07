package erlang

import (
	"elbe-prj/containers"

	"github.com/charmbracelet/bubbles/textinput"
)

type Model struct {
	textInput textinput.Model
	err       error
	choices   []string
	selected  map[int]struct{}
	cursor    int
	get_it    containers.DownloadState
	projects  []containers.Project
	debianize bool
	Inputs    []textinput.Model
	Focused   int
}
