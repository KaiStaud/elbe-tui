/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"elbe-prj/cmd"
	"elbe-prj/containers"
	"elbe-prj/utils"
	"fmt"
	"log" // TODO: use zap instead
	"os"
	"os/exec"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/termenv"
	"github.com/spf13/viper"
)

type KeyMap struct {
	Up       key.Binding
	Down     key.Binding
	Reset    key.Binding
	Delete   key.Binding
	GetFiles key.Binding
	Submit   key.Binding
	Package  key.Binding
}

var (
	selectedItemStyle = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("170"))
	failedStyle       = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("9"))
	doneStyle         = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("#04B575"))
	unusedStyle       = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("#3C3C3C"))
	busyStyle         = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("#0000FF"))
	errorStyle        = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("##ff5100"))
	keyword           = utils.MakeFgStyle("211")
	subtle            = utils.MakeFgStyle("241")
	dot               = utils.ColorFg(" • ", "236")
	term              = termenv.EnvColorProfile()
	pbuild_prj        = "/var/cache/elbe/a79a01ed-9091-4f8f-9f20-1ed6a7060634+"
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
			key.WithKeys("r", "reset"),
			key.WithHelp("r", "reset"),
		),
		Delete: key.NewBinding(
			key.WithKeys("t", "delete"),
			key.WithHelp("t", "delete"),
		),
		GetFiles: key.NewBinding(
			key.WithKeys("g", "get_files"),
			key.WithHelp("g", "get files"),
		),
		Package: key.NewBinding(
			key.WithKeys("p", "make_package"),
			key.WithHelp("p", "make deb package"),
		),
	}
	elbe_bin    = ""
	elbe_dl_dir = ""
)

const (
	kernel_release = iota
	architecture
	defconfig
	src_package
	src_path
)

const (
	hotPink  = lipgloss.Color("#FF06B7")
	darkGray = lipgloss.Color("#767676")
)

var (
	inputStyle    = lipgloss.NewStyle().Foreground(hotPink)
	continueStyle = lipgloss.NewStyle().Foreground(darkGray)
)

type elbe_hook struct {
	path   string
	dl_dir string
}

type model struct {
	textInput textinput.Model
	err       error
	choices   []string
	selected  map[int]struct{}
	cursor    int
	get_it    containers.DownloadState
	projects  []containers.Project
	debianize bool
	inputs    []textinput.Model
	focused   int
}

func initialModel(p []containers.Project) model {
	ti := textinput.New()
	ti.Placeholder = elbe_dl_dir
	ti.Focus()
	ti.CharLimit = 156
	ti.Width = 20
	var names []string
	for _, v := range p {
		names = append(names, v.Name)
	}

	var inputs []textinput.Model = make([]textinput.Model, 5)
	inputs[src_package] = textinput.New()
	inputs[src_package].Placeholder = "SrcPackageName"
	inputs[src_package].Focus()
	inputs[src_package].CharLimit = 30
	inputs[src_package].Width = 30
	inputs[src_package].Prompt = ""

	inputs[src_path] = textinput.New()
	inputs[src_path].Placeholder = "~/path/to/source"
	inputs[src_path].CharLimit = 30
	inputs[src_path].Width = 30
	inputs[src_path].Prompt = ""

	inputs[kernel_release] = textinput.New()
	inputs[kernel_release].Placeholder = "6.1.27"
	inputs[kernel_release].CharLimit = 30
	inputs[kernel_release].Width = 30
	inputs[kernel_release].Prompt = ""

	inputs[architecture] = textinput.New()
	inputs[architecture].Placeholder = "arm"
	inputs[architecture].CharLimit = 5
	inputs[architecture].Width = 5
	inputs[architecture].Prompt = ""

	inputs[defconfig] = textinput.New()
	inputs[defconfig].Placeholder = "stm32mp157a-dk1_defconfig"
	inputs[defconfig].CharLimit = 50
	inputs[defconfig].Width = 50
	inputs[defconfig].Prompt = ""
	return model{
		textInput: ti,
		err:       nil,
		choices:   names,
		selected:  make(map[int]struct{}),
		projects:  p,
		inputs:    inputs,
		focused:   0,
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
		// disable keybinds while textinput is active:
		if m.get_it == containers.DownloadStarted {
			switch msg.String() {
			case "esc":
				m.get_it = containers.DownloadFinished
			case "enter", " ":
				if m.get_it == containers.DownloadStarted {
					utils.GetFiles(m.projects[m.cursor].Path, m.textInput.Value())
					m.get_it = containers.DownloadFinished
					log.Printf("Downloaded to %s", m.textInput.Value())
				}
			}
		} else if m.debianize == true {
			switch msg.Type {
			case tea.KeyEnter, tea.KeySpace, tea.KeyEsc:
				m.debianize = false
			case tea.KeyShiftTab, tea.KeyCtrlP:
				m.prevInput()
			case tea.KeyTab, tea.KeyCtrlN:
				m.nextInput()
			}
			for i := range m.inputs {
				m.inputs[i].Blur()
			}
			m.inputs[m.focused].Focus()
			var cmds []tea.Cmd = make([]tea.Cmd, len(m.inputs))
			for i := range m.inputs {
				m.inputs[i], cmds[i] = m.inputs[i].Update(msg)
			}
			return m, tea.Batch(cmds...)
		} else {
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
				log.Printf("deleting %s@%s", m.choices[m.cursor], m.projects[m.cursor].Path)
				utils.ResetProject(m.projects[m.cursor].Path)
				// TODO: check for error
				m.projects[m.cursor].Result = containers.Needs_Build

			case "get_files", "g":
				log.Printf("getting %s@%s", m.choices[m.cursor], m.projects[m.cursor].Path)
				m.projects[m.cursor].Progress = containers.DownloadStarted
				m.get_it = containers.DownloadStarted
				// TODO: check for error
				// TODO: append a "-> <download-path>"
			case "delete", "t":
				log.Printf("deleting %s@%s", m.choices[m.cursor], m.projects[m.cursor].Path)
				utils.DeleteProject(m.projects[m.cursor].Path, false)
				delete(m.selected, m.cursor)

				i := m.cursor
				copy(m.choices[i:], m.choices[i+1:])     // Shift a[i+1:] left one index.
				m.choices[len(m.choices)-1] = ""         // Erase last element (write zero value).
				m.choices = m.choices[:len(m.choices)-1] // Truncate slice.

				copy(m.projects[i:], m.projects[i+1:])               // Shift a[i+1:] left one index.
				m.projects[len(m.projects)-1] = containers.Project{} // Erase last element (write zero value).
				m.projects = m.projects[:len(m.projects)-1]          // Truncate slice.

			case "make_package", "p":
				m.debianize = true
				// Download-Dir entered,switch back to list-view
			case "esc":
				m.get_it = containers.DownloadFinished
			case "enter", " ":
				if m.get_it == containers.DownloadStarted {
					m.get_it = containers.DownloadFinished
					utils.GetFiles(m.projects[m.cursor].Path, m.textInput.Value())
					log.Printf("Downloaded to %s", m.textInput.Value())
				}
			}
		}
	}
	m.textInput, cmd = m.textInput.Update(msg)
	return m, cmd
}

func (m model) View() string {
	s := ""
	if m.get_it == containers.DownloadStarted {
		return fmt.Sprintf(
			"Enter download path:\n\n%s\n\n%s",
			m.textInput.View(),
			"(esc to quit)",
		) + "\n"
	}

	if m.debianize == true {
		return fmt.Sprintf(
			` Enter build parameters:
	 %s %s  %s %s
	 %s
	 %s
	 %s
	 %s
	 %s
	 %s

	 %s
	`,
			inputStyle.Width(30).Render("SrcPackage"),
			m.inputs[src_package].View(),
			inputStyle.Width(30).Render("SourcePath"),
			m.inputs[src_path].View(),
			inputStyle.Width(30).Render("Kernel Release"),
			m.inputs[kernel_release].View(),
			inputStyle.Width(13).Render("Architecture"),
			m.inputs[architecture].View(),
			inputStyle.Width(6).Render("Config"),
			m.inputs[defconfig].View(),
			continueStyle.Render("Press Enter to submit"),
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
		var result = utils.ColorizeBuildResult(m.projects[i])
		// Render the row
		s += fmt.Sprintf("%s %s [%s] %s\n", cursor, result, checked, choice)

	}
	var tpl = subtle("q,ctrl+c: quit") + dot + subtle("j/k, up/down: select") + "\n" +
		subtle("r: reset_project") + dot + subtle("t:  delete_project") + dot + subtle("g: get_files") +
		dot + subtle("p: debianize source")
	s += fmt.Sprintf(tpl)
	return s

}

// nextInput focuses the next input field
func (m *model) nextInput() {
	m.focused = (m.focused + 1) % len(m.inputs)
}

// prevInput focuses the previous input field
func (m *model) prevInput() {
	m.focused--
	// Wrap around
	if m.focused < 0 {
		m.focused = len(m.inputs) - 1
	}
}
func main() {
	f, err := os.OpenFile("elbe.go.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()

	log.SetOutput(f)
	cmd.Execute()

	log.Println("Reading in config config.env")
	viper.SetConfigType("json")
	viper.AddConfigPath("/home/sta/projects/go/elbe-prj")
	viper.SetConfigFile("config.json")
	viper.ReadInConfig()

	elbe_bin = viper.GetString("elbe")
	elbe_dl_dir = viper.GetString("default_dl_dir")
	log.Println("elbe bin is located at" + elbe_bin + ", default-dl-dir is " + elbe_dl_dir)
	utils.LoadConfig(elbe_bin)

	app := elbe_bin

	arg0 := "control"
	arg1 := "list_projects"
	arg2 := ""
	arg3 := ""

	cmd := exec.Command(app, arg0, arg1, arg2, arg3)
	stdout, err := cmd.Output()

	if err != nil {
		log.Printf("Couldn't get initial project list from elbe-cmd %s %s %s, maybe your config.env isnt't handled correctly", app, arg0, arg1)
		log.Printf("get_projects backtrace:%s", err.Error())
		return
	}

	var projects []containers.Project
	s := utils.SplitLines(string(stdout))

	for i, v := range s {
		log.Println(i, v)
		p := utils.ParseLine(v)
		projects = append(projects, p)
	}

	p := tea.NewProgram(initialModel(projects))
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
