package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

type BuildResult int

const (
	Build_Failed  BuildResult = iota // Read   = 0
	Build_Done                       // Create = 1
	Empty_Project                    // Update = 2
	Busy                             // Delete = 3
	Needs_Build                      // List   = 4
)

var (
	BuildResultMap = map[string]BuildResult{
		"build_failed":  Build_Failed,
		"build_done":    Build_Done,
		"empty_project": Empty_Project,
		"busy":          Busy,
		"needs_build":   Needs_Build,
	}
)

type project struct {
	path   string
	name   string
	result BuildResult
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
	fmt.Println(words, len(words))
	c, _ := BuildResultMap[strings.ToLower(words[3])]
	return project{path: words[0], name: words[1], result: c} //,result:matched_result,builddate:ts}
}

func DeleteProject(path string, needs_reset bool) {
	app := "/home/sta/projects/elbe/elbe"
	arg0 := "control"
	arg1 := "del_project"
	arg2 := path

	log.Printf(" %s %s %s %s", app, arg0, arg1, arg2)
	cmd := exec.Command("/home/sta/projects/elbe/elbe", "control", "del_project", path)
	stdout, err := cmd.Output()
	log.Println(stdout)
	if err != nil {
		log.Println(err.Error())
		fmt.Println(err.Error())
		return
	}

}

type model struct {
	choices    []string         // items on the to-do list
	cursor     int              // which to-do list item our cursor is pointing at
	selected   map[int]struct{} // which to-do items are selected
	projects   []project
	is_reseted bool
	is_deleted bool
}

func initialModel(p []project) model {

	var names []string
	for _, v := range p {
		names = append(names, v.name)
	}

	return model{
		// Our to-do list is a grocery list

		choices:  names,
		projects: p,
		// A map which indicates which choices are selected. We're using
		// the  map like a mathematical set. The keys refer to the indexes
		// of the `choices` slice, above.
		selected: make(map[int]struct{}),
	}
}

func (m model) Init() tea.Cmd {
	// Just return `nil`, which means "no I/O right now, please."
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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

		// The "enter" key and the spacebar (a literal space) toggle
		// the selected state for the item that the cursor is pointing at.
		case "enter", " ":
			_, ok := m.selected[m.cursor]
			if ok {
				log.Printf("deleting %s@%s", m.choices[m.cursor], m.projects[m.cursor].path)
				DeleteProject(m.projects[m.cursor].path, false)
				delete(m.selected, m.cursor)

				//todo: also delete project from slice choices and projects!
				i := m.cursor
				copy(m.choices[i:], m.choices[i+1:])     // Shift a[i+1:] left one index.
				m.choices[len(m.choices)-1] = ""         // Erase last element (write zero value).
				m.choices = m.choices[:len(m.choices)-1] // Truncate slice.

				copy(m.projects[i:], m.projects[i+1:])      // Shift a[i+1:] left one index.
				m.projects[len(m.projects)-1] = project{}   // Erase last element (write zero value).
				m.projects = m.projects[:len(m.projects)-1] // Truncate slice.

			} else {
				m.selected[m.cursor] = struct{}{}
			}
		}
	}

	// Return the updated model to the Bubble Tea runtime for processing.
	// Note that we're not returning a command.
	return m, nil
}

func (m model) View() string {
	// The header
	s := "Press enter to delete project\n\n"

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

		// Render the row
		s += fmt.Sprintf("%s [%s] %s\n", cursor, checked, choice)
	}

	// The footer
	s += "\nPress q to quit.\n"

	// Send the UI for rendering
	return s
}

func main() {
	f, err := os.OpenFile("elbe.go.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()

	log.SetOutput(f)
	log.Println("Starting new session")

	app := "/home/sta/projects/elbe/elbe"

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
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
