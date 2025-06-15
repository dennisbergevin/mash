package main

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"
)

func captureStdout(fn func()) string {
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	fn()

	w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	io.Copy(&buf, r)
	return buf.String()
}

func TestTagStyleFor_ProducesDeterministicColor(t *testing.T) {
	style1 := tagStyleFor("dev")
	style2 := tagStyleFor("dev")
	style3 := tagStyleFor("infra")

	if style1.String() != style2.String() {
		t.Errorf("Expected deterministic styles for same tag, got different: %v vs %v", style1, style2)
	}
	if style1.String() == style3.String() {
		t.Errorf("Expected different styles for different tags: 'dev' and 'infra' got same")
	}
}

func TestPrintTagTree_BasicRendering(t *testing.T) {
	config := Config{
		Items: []ConfigItem{
			{Title: "Build", Desc: "Run build", Cmd: "make build", Tags: []string{"dev"}},
			{Title: "Deploy", Desc: "Deploy app", Cmd: "make deploy", Tags: []string{"prod"}},
			{Title: "Untagged", Desc: "No tag", Cmd: "true"},
		},
	}

	output := captureStdout(func() {
		// No includeTags and requireAnyTags = false → include all
		printTagTree(config, []string{}, nil, false)
	})

	if !strings.Contains(output, "Build") || !strings.Contains(output, "Deploy") {
		t.Errorf("Expected tagged items to appear in tree output")
	}
	if !strings.Contains(output, "(untagged)") || !strings.Contains(output, "Untagged") {
		t.Errorf("Expected untagged section in tree output")
	}
}

func TestPrintTagTree_WithTagFilter(t *testing.T) {
	config := Config{
		Items: []ConfigItem{
			{Title: "Build", Desc: "Build it", Cmd: "make", Tags: []string{"dev"}},
			{Title: "SkipMe", Desc: "Not dev", Cmd: "skip", Tags: []string{"infra"}},
		},
	}

	output := captureStdout(func() {
		printTagTree(config, []string{"dev"}, nil, false)
	})

	if !strings.Contains(output, "Build") {
		t.Errorf("Expected 'Build' to be in tree output")
	}
	if strings.Contains(output, "SkipMe") {
		t.Errorf("Did not expect 'SkipMe' to appear with tag filter")
	}
}

func TestPrintTagTree_OnlyAnyTags(t *testing.T) {
	config := Config{
		Items: []ConfigItem{
			{Title: "TaggedItem", Desc: "Tagged", Cmd: "true", Tags: []string{"utils"}},
			{Title: "UntaggedItem", Desc: "No tags", Cmd: "echo"},
		},
	}

	output := captureStdout(func() {
		// Simulate the case where user passed --tag with no value (requireAnyTags = true)
		printTagTree(config, []string{}, nil, true)
	})

	if !strings.Contains(output, "TaggedItem") {
		t.Errorf("Expected tagged item to be shown when requireAnyTags is true")
	}
	if strings.Contains(output, "UntaggedItem") {
		t.Errorf("Did not expect untagged item to be shown when requireAnyTags is true")
	}
}
