package main

import (
	"strings"
	"testing"
)

func contains(s, sub string) bool {
	return len(s) >= len(sub) && strings.Contains(s, sub)
}

func TestItem_TitleRenderingWithTags(t *testing.T) {
	it := item{
		title: "TestCmd",
		tags:  []string{"dev", "infra"},
	}

	rendered := it.Title()

	if !contains(rendered, "TestCmd") || !contains(rendered, "dev") || !contains(rendered, "infra") {
		t.Errorf("Expected title to include rendered tags and base title, got: %s", rendered)
	}
}

func TestItem_TitleRenderingWithoutTags(t *testing.T) {
	it := item{
		title: "SimpleCmd",
		tags:  []string{},
	}

	rendered := it.Title()

	if rendered != "SimpleCmd" {
		t.Errorf("Expected plain title without tags, got: %s", rendered)
	}
}
