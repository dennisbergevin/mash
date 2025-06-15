package main

import (
	"crypto/sha256"
	"fmt"
	"sort"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/tree"
)

func tagStyleFor(tag string) lipgloss.Style {
	hash := sha256.Sum256([]byte(tag))
	r, g, b := hash[0], hash[1], hash[2]
	bgColor := fmt.Sprintf("#%02x%02x%02x", r, g, b)

	fgColor := "#000000"
	if brightness := 0.299*float64(r) + 0.587*float64(g) + 0.114*float64(b); brightness < 128 {
		fgColor = "#ffffff"
	}

	return lipgloss.NewStyle().
		Foreground(lipgloss.Color(fgColor)).
		Background(lipgloss.Color(bgColor)).
		Padding(0, 1)
}

func printTagTree(config Config, includeTags, excludeTags []string, requireAnyTags bool) {
	items := config.Items
	tagMap := make(map[string][]ConfigItem)
	var untaggedItems []ConfigItem

	for _, item := range items {
		if !shouldInclude(item.Tags, includeTags, requireAnyTags) {
			continue
		}
		if len(item.Tags) == 0 {
			untaggedItems = append(untaggedItems, item)
			continue
		}
		if len(includeTags) > 0 && includeTags[0] != "__ANY__" {
			for _, tag := range item.Tags {
				for _, inc := range includeTags {
					if strings.EqualFold(tag, inc) {
						tagMap[tag] = append(tagMap[tag], item)
						break
					}
				}
			}
		} else {
			for _, tag := range item.Tags {
				tagMap[tag] = append(tagMap[tag], item)
			}
		}
	}

	root := tree.Root("mash").
		Enumerator(tree.RoundedEnumerator).
		EnumeratorStyle(lipgloss.NewStyle().Foreground(lipgloss.Color("240"))).
		RootStyle(rootStyle)

	tagStyle := lipgloss.NewStyle().Bold(true)
	titleStyle := lipgloss.NewStyle()
	descStyle := lipgloss.NewStyle().Faint(true)

	if config.TagColor != "" {
		tagStyle = tagStyle.Foreground(lipgloss.Color(config.TagColor))
	}
	if config.TitleColor != "" {
		titleStyle = titleStyle.Foreground(lipgloss.Color(config.TitleColor))
	}
	if config.DescColor != "" {
		descStyle = lipgloss.NewStyle().Foreground(lipgloss.Color(config.DescColor))
	}

	var tags []string
	for tag := range tagMap {
		tags = append(tags, tag)
	}
	sort.Strings(tags)

	for _, tag := range tags {
		tagNode := tree.Root(tagStyle.Render(tag))
		for _, item := range tagMap[tag] {
			itemNode := tree.Root(titleStyle.Render(item.Title))
			if item.Desc != "" {
				itemNode.Child(tree.Root(descStyle.Render(item.Desc)))
			}
			tagNode.Child(itemNode)
		}
		root.Child(tagNode)
	}

	if len(untaggedItems) > 0 {
		untaggedNode := tree.Root(tagStyle.Render("(untagged)"))
		for _, item := range untaggedItems {
			itemNode := tree.Root(titleStyle.Render(item.Title))
			if item.Desc != "" {
				itemNode.Child(tree.Root(descStyle.Render(item.Desc)))
			}
			untaggedNode.Child(itemNode)
		}
		root.Child(untaggedNode)
	}

	fmt.Println()
	fmt.Println(root.String())
	fmt.Println()
}
