package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	docStyle  = lipgloss.NewStyle().Margin(1, 2)
	rootStyle = lipgloss.NewStyle().Bold(true).Background(lipgloss.Color("62")).Foreground(lipgloss.Color("230")).Padding(0, 1)
)

type item struct {
	title, desc, cmd string
	tags             []string
}

func isWindows() bool {
	return strings.Contains(strings.ToLower(os.Getenv("OS")), "windows") || os.PathSeparator == '\\'
}

func (i item) Title() string {
	base := i.title
	if len(i.tags) == 0 {
		return base
	}
	var styledTags []string
	for _, tag := range i.tags {
		style := tagStyleFor(tag)
		styledTags = append(styledTags, style.Render(tag))
	}
	return fmt.Sprintf("%s  %s", base, lipgloss.JoinHorizontal(lipgloss.Left, styledTags...))
}

func (i item) Description() string { return i.desc }
func (i item) FilterValue() string { return i.title }

type model struct {
	list     list.Model
	choice   string
	quitting bool
}

type keymap struct {
	Choose key.Binding
}

var keyMap = keymap{
	Choose: key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "choose")),
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "enter":
			if m.list.FilterState() != list.Filtering {
				i, ok := m.list.SelectedItem().(item)
				if ok {
					m.choice = i.cmd
					m.quitting = true
					return m, openCommand(i.cmd)
				}
			}
		}
	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)
	}

	if m.View() == "" {
		return m, tea.Quit
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m model) View() string {
	if m.quitting {
		return ""
	}
	return docStyle.Render(m.list.View())
}

func openCommand(com string) tea.Cmd {
	if strings.HasPrefix(com, "cd:") {
		dir := os.ExpandEnv(strings.TrimPrefix(com, "cd:"))
		if isWindows() {
			return tea.ExecProcess(exec.Command("cmd", "/K", "cd /d "+dir), nil)
		}
		return tea.ExecProcess(exec.Command("sh", "-c", fmt.Sprintf("cd %s && exec $SHELL", dir)), nil)
	}

	if isWindows() {
		return tea.ExecProcess(exec.Command("cmd", "/C", com), nil)
	}
	return tea.ExecProcess(exec.Command("sh", "-c", com), nil)
}
