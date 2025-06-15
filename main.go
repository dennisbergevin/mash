package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
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

	helpText := fmt.Sprintf(
		"%s - A customizable command launcher\n\n"+
			"Usage:\n  mash [options]\n\n"+
			"Options:\n"+
			"  --help, -h         Show this help menu\n"+
			"  --tree             Display a tree view of all commands and tags\n"+
			"  --tag              Show only items that have tags (any tags)\n"+
			"  --tag=\"TAG1;TAG2\"  Show items with specified tags\n"+
			"  --skip-intro       Skip splash screen\n"+
			"  --global, -g       Load config from ~/.config/mash/config.json\n\n"+
			"Examples:\n"+
			"  mash --tag=dev\n"+
			"  mash --tree --tag=\"tools;infra\"\n"+
			"  mash --global --skip-intro\n",
		appName)

	fmt.Println(helpText)
}

func splashScreen() *huh.Note {
	return huh.NewNote().
		Title(fmt.Sprintf("\n%s", rootStyle.Render("mash"))).
		Description("A customizable command launcher\nfor storing and executing commands")
}
