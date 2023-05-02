package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type BuildResult int
type DownloadState int

const (
	Build_Failed BuildResult = iota // 0
	Build_Done
	Empty_Project
	Busy
	Needs_Build                          // 4
	DownloadStarted DownloadState = iota // 5
	DownloadPathEntered
	DownloadFinished
)

type KeyMap struct {
	Up       key.Binding
	Down     key.Binding
	Reset    key.Binding
	Delete   key.Binding
	GetFiles key.Binding
}

var (
	BuildResultMap = map[string]BuildResult{
		"build_failed":  Build_Failed,
		"build_done":    Build_Done,
		"empty_project": Empty_Project,
		"busy":          Busy,
		"needs_build":   Needs_Build,
	}
	selectedItemStyle = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("170"))
	failedStyle       = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("9"))
	doneStyle         = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("#04B575"))
	unusedStyle       = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("#3C3C3C"))
	busyStyle         = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("#0000FF"))
	DefaultKeyMap     = KeyMap{
		Up: key.NewBinding(
			key.WithKeys("k", "up"),        // actual keybindings
			key.WithHelp("↑/k", "move up"), // corresponding help text
		),
		Down: key.NewBinding(
			key.WithKeys("j", "down"),
			key.WithHelp("↓/j", "move down"),
		),
		Reset: key.NewBinding(
			key.WithKeys("r", "reset"), // actual keybindings
			key.WithHelp("r", "reset"), // corresponding help text
		),
		Delete: key.NewBinding(
			key.WithKeys("t", "delete"),
			key.WithHelp("t", "delete"),
		),
		GetFiles: key.NewBinding(
			key.WithKeys("g", "get_files"),
			key.WithHelp("g", "get files"),
		),
	}
	elbe_bin    = ""
	elbe_dl_dir = ""
)

type project struct {
	path     string
	name     string
	result   BuildResult
	progress DownloadState
	//builddate
}

func SplitLines(s string) []string {
	var lines []string
	sc := bufio.NewScanner(strings.NewReader(s))
	for sc.Scan() {
		lines = append(lines, sc.Text())
	}
	return lines
}

func ParseLine(s string) project {
	words := strings.Fields(s)
	c, _ := BuildResultMap[strings.ToLower(words[4])]
	return project{path: words[0], name: words[1], result: c} //,result:matched_result,builddate:ts}
}

type elbe_hook struct {
	path   string
	dl_dir string
}

func DeleteProject(path string, needs_reset bool) {
	app := "/hdd/elbe/elbe"
	arg0 := "control"
	arg1 := "del_project"
	arg2 := path

	log.Printf(" %s %s %s %s", app, arg0, arg1, arg2)
	cmd := exec.Command("/hdd/elbe/elbe", "control", "del_project", path)
	stdout, err := cmd.Output()
	log.Println(stdout)
	if err != nil {
		log.Println(err.Error())
		fmt.Println(err.Error())
		return
	}

}

func GetFiles(path string, target_dir string) {
	app := "/hdd/elbe/elbe"
	arg0 := "control"
	arg1 := "get_files"
	arg2 := path

	log.Printf(" %s %s %s %s", app, arg0, arg1, arg2)
	cmd := exec.Command("/hdd/elbe/elbe", "control", "get_files", "--output", "None", path)
	stdout, err := cmd.Output()
	log.Println(stdout)
	if err != nil {
		log.Println(err.Error())
		fmt.Println(err.Error())
		return
	}

}

func ResetProject(path string) {
	cmd := exec.Command("/hdd/elbe/elbe", "control", "reset_project", path)
	stdout, err := cmd.Output()
	log.Println(stdout)
	if err != nil {
		log.Println(err.Error())
		fmt.Println(err.Error())
		return
	}
}

type model struct {
	textInput textinput.Model
	err       error
	choices   []string
	selected  map[int]struct{}
	cursor    int
	get_it    DownloadState
	projects  []project
}

func initialModel(p []project) model {
	ti := textinput.New()
	ti.Placeholder = "Pikachu"
	ti.Focus()
	ti.CharLimit = 156
	ti.Width = 20
	var names []string
	for _, v := range p {
		names = append(names, v.name)
	}
	return model{
		textInput: ti,
		err:       nil,
		choices:   names,
		selected:  make(map[int]struct{}),
		projects:  p,
	}
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	// Is it a key press?
	case tea.KeyMsg:

		// Cool, what was the actual key pressed?
		switch msg.String() {

		// These keys should exit the program.
		case "ctrl+c", "q":
			return m, tea.Quit

		// The "up" and "k" keys move the cursor up
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}

		// The "down" and "j" keys move the cursor down
		case "down", "j":
			if m.cursor < len(m.choices)-1 {
				m.cursor++
			}
		case "reset", "r":
			log.Printf("deleting %s@%s", m.choices[m.cursor], m.projects[m.cursor].path)
			ResetProject(m.projects[m.cursor].path)
			// todo: check for error
			m.projects[m.cursor].result = Needs_Build

		case "get_files", "g":
			log.Printf("getting %s@%s", m.choices[m.cursor], m.projects[m.cursor].path)
			m.projects[m.cursor].progress = DownloadStarted
			m.get_it = DownloadStarted
			// todo: check for error
			// todo: append a "-> <download-path>"
		case "delete", "t":
			log.Printf("deleting %s@%s", m.choices[m.cursor], m.projects[m.cursor].path)
			DeleteProject(m.projects[m.cursor].path, false)
			delete(m.selected, m.cursor)

			i := m.cursor
			copy(m.choices[i:], m.choices[i+1:])     // Shift a[i+1:] left one index.
			m.choices[len(m.choices)-1] = ""         // Erase last element (write zero value).
			m.choices = m.choices[:len(m.choices)-1] // Truncate slice.

			copy(m.projects[i:], m.projects[i+1:])      // Shift a[i+1:] left one index.
			m.projects[len(m.projects)-1] = project{}   // Erase last element (write zero value).
			m.projects = m.projects[:len(m.projects)-1] // Truncate slice.

		// Download-Dir entered,switch back to list-view
		case "enter", " ":
			if m.get_it == DownloadStarted {
				m.get_it = DownloadFinished
				log.Printf("Downloaded to %s", m.textInput.Value())
			}
		}
	}
	m.textInput, cmd = m.textInput.Update(msg)
	return m, cmd
}

func colorizeBuildResult(p project) string {
	switch p.result {
	case Build_Done:
		return doneStyle.Render("[done]")
	case Busy:
		return busyStyle.Render("[busy]")
	case Build_Failed:
		return failedStyle.Render("[failed]")
	case Needs_Build:
		return unusedStyle.Render("[needs build]")
	default:
		return ""
	}

}

func (m model) View() string {
	s := ""
	if m.get_it == DownloadStarted {
		return fmt.Sprintf(
			"Enter download path:\n\n%s\n\n%s",
			m.textInput.View(),
			"(esc to quit)",
		) + "\n"
	}
	// Iterate over our choices
	for i, choice := range m.choices {

		// Is the cursor pointing at this choice?
		cursor := " " // no cursor
		if m.cursor == i {
			cursor = ">" // cursor!
		}

		// Is this choice selected?
		checked := " " // not selected
		if _, ok := m.selected[i]; ok {
			checked = "x" // selected!
		}
		var result = colorizeBuildResult(m.projects[i])
		// Render the row
		s += fmt.Sprintf("%s %s [%s] %s\n", cursor, result, checked, choice)
		s += fmt.Sprintf("Press q to quit, t to delete and r to reset project\n")
		s += fmt.Sprintf("Press g to download files from initvm\n")

	}
	return s

}

func main() {
	/*
		viper.AddConfigPath("/hdd/go/elbe-prj")
		viper.SetConfigFile("config.env")
		viper.ReadInConfig()

		elbe_bin = viper.Get("ELBE").(string)
		elbe_dl_dir = viper.Get("DEFAULT_DOWNLOAD_DIR").(string)
	*/
	f, err := os.OpenFile("elbe.go.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()

	log.SetOutput(f)
	log.Println("Starting new session")

	app := "/hdd/elbe/elbe"

	arg0 := "control"
	arg1 := "list_projects"
	arg2 := ""
	arg3 := ""

	cmd := exec.Command(app, arg0, arg1, arg2, arg3)
	stdout, err := cmd.Output()

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	var projects []project
	s := SplitLines(string(stdout))

	for i, v := range s {
		log.Println(i, v)
		p := ParseLine(v)
		projects = append(projects, p)
	}

	p := tea.NewProgram(initialModel(projects))
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}

type (
	errMsg error
)
