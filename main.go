package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
)

func main() {
	var includeTags, excludeTags []string
	var requireAnyTags bool

	showTree := false
	useGlobal := false
	skipIntro := false

	for i := 0; i < len(os.Args); i++ {
		arg := os.Args[i]
		switch {
		case arg == "--skip-intro":
			skipIntro = true
		case arg == "--tree":
			showTree = true
		case arg == "--global" || arg == "-g":
			useGlobal = true
		case arg == "--help" || arg == "-h":
			printHelp()
			return
		case strings.HasPrefix(arg, "--tag="):
			tags := strings.Split(strings.TrimPrefix(arg, "--tag="), ";")
			for _, tag := range tags {
				tag = strings.TrimSpace(tag)
				if tag != "" {
					includeTags = append(includeTags, tag)
				}
			}
		case arg == "--tag":
			if i+1 < len(os.Args) && !strings.HasPrefix(os.Args[i+1], "--") {
				i++
				tags := strings.Split(os.Args[i], ";")
				for _, tag := range tags {
					tag = strings.TrimSpace(tag)
					if tag != "" {
						includeTags = append(includeTags, tag)
					}
				}
			} else {
				includeTags = nil
				requireAnyTags = true
			}
		}
	}

	config, err := loadConfig(useGlobal)
	if err != nil {
		fmt.Println("Error loading config:", err)
		os.Exit(1)
	}

	if showTree {
		printTagTree(config, includeTags, excludeTags, requireAnyTags)
		return
	}

	var listItems []list.Item
	for _, ci := range config.Items {
		if shouldInclude(ci.Tags, includeTags, requireAnyTags) {
			listItems = append(listItems, item{
				title: ci.Title,
				desc:  ci.Desc,
				cmd:   ci.Cmd,
				tags:  ci.Tags,
			})
		}
	}

	if len(listItems) == 0 {
		fmt.Println("No items found.")
		os.Exit(1)
	}

	if !skipIntro && !config.SkipIntro {
		if err := splashScreen().Run(); err != nil {
			fmt.Println("Error displaying splash screen:", err)
			os.Exit(1)
		}
	}

	m := model{list: list.New(listItems, list.NewDefaultDelegate(), 0, 0)}
	m.list.AdditionalShortHelpKeys = func() []key.Binding { return []key.Binding{keyMap.Choose} }
	m.list.AdditionalFullHelpKeys = func() []key.Binding { return []key.Binding{keyMap.Choose} }
	m.list.Title = "mash"

	p := tea.NewProgram(m, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}

func printHelp() {
	appName := rootStyle.Render("mash")

	// Define styles
	sectionTitle := lipgloss.NewStyle().Foreground(lipgloss.Color("62")).Bold(true)
	description := lipgloss.NewStyle().Faint(true)

	borderBox := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("62")). // Orange border
		Padding(1, 2).
		Margin(1, 0)

	// Help content
	var b strings.Builder
	fmt.Fprintf(&b, "%s - A customizable command launcher\n\n", appName)

	fmt.Fprintln(&b, sectionTitle.Render("Usage"))
	fmt.Fprintln(&b, "  mash [options]")

	fmt.Fprintln(&b, "\n"+sectionTitle.Render("Options"))
	fmt.Fprintln(&b, "  --help, -h         "+description.Render("Show this help menu"))
	fmt.Fprintln(&b, "  --tree             "+description.Render("Display a tree view of all commands and tags"))
	fmt.Fprintln(&b, "  --tag              "+description.Render("Show only items that have tags (any tags)"))
	fmt.Fprintln(&b, "  --tag=\"TAG1;TAG2\"  "+description.Render("Show items with specified tags"))
	fmt.Fprintln(&b, "  --skip-intro       "+description.Render("Skip splash screen"))
	fmt.Fprintln(&b, "  --global, -g       "+description.Render("Load config from ~/.config/mash/config.json"))

	fmt.Fprintln(&b, "\n"+sectionTitle.Render("Examples"))
	fmt.Fprintln(&b, "  mash --tag=dev")
	fmt.Fprintln(&b, "  mash --tree --tag=\"tools;infra\"")
	fmt.Fprintln(&b, "  mash --global --skip-intro")

	// Print bordered help menu
	fmt.Println(borderBox.Render(b.String()))
}

func splashScreen() *huh.Note {
	return huh.NewNote().
		Title(fmt.Sprintf("\n%s", rootStyle.Render("mash"))).
		Description("A customizable command launcher\nfor storing and executing commands")
}
