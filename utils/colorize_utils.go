package utils

import (
	"elbe-prj/containers"

	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/termenv"
)

var (
	selectedItemStyle = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("170")) // green
	failedStyle       = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("9"))
	doneStyle         = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("#04B575"))
	unusedStyle       = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("#3C3C3C"))
	busyStyle         = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("#0000FF"))
	errorStyle        = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("##ff5100"))
	pbuilderStyle     = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("#ff00aa"))
	term              = termenv.EnvColorProfile()
)

func ColorizeBuildResult(p containers.Project) string {
	switch p.Result {
	case containers.Build_Done:
		return doneStyle.Render("[done]")
	case containers.Busy:
		return busyStyle.Render("[busy]")
	case containers.Build_Failed:
		return failedStyle.Render("[failed]")
	case containers.Needs_Build:
		return unusedStyle.Render("[needs build]")
	default:
		return ""
	}

}

func ColorizeErrorMessage(s string) string {
	return errorStyle.Render(s)
}

// Color a string's foreground with the given value.
func ColorFg(val, color string) string {
	return termenv.String(val).Foreground(term.Color(color)).String()
}

// Return a function that will colorize the foreground of a given string.
func MakeFgStyle(color string) func(string) string {
	return termenv.Style{}.Foreground(term.Color(color)).Styled
}
