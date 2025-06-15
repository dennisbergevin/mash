package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
)

type ConfigItem struct {
	Title string   `json:"title"`
	Desc  string   `json:"desc"`
	Cmd   string   `json:"cmd"`
	Tags  []string `json:"tags,omitempty"`
}

type Config struct {
	SkipIntro  bool         `json:"skipIntro"`
	Items      []ConfigItem `json:"items"`
	TagColor   string       `json:"tagColor,omitempty"`
	TitleColor string       `json:"titleColor,omitempty"`
	DescColor  string       `json:"descColor,omitempty"`
}

func loadConfig(useGlobal bool) (Config, error) {
	if useGlobal {
		globalPath := filepath.Join(os.Getenv("HOME"), ".config", "mash", "config.json")
		return loadConfigFile(globalPath)
	}

	dir, err := os.Getwd()
	if err != nil {
		return Config{}, err
	}

	for {
		localPath := filepath.Join(dir, ".mash.json")
		if _, err := os.Stat(localPath); err == nil {
			return loadConfigFile(localPath)
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}

	globalPath := filepath.Join(os.Getenv("HOME"), ".config", "mash", "config.json")
	return loadConfigFile(globalPath)
}

func loadConfigFile(path string) (Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return Config{}, err
	}

	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		var items []ConfigItem
		if err := json.Unmarshal(data, &items); err != nil {
			return Config{}, err
		}
		config.Items = items
	}
	return config, nil
}

func shouldInclude(itemTags, includeTags []string, requireAnyTags bool) bool {
	if requireAnyTags && len(itemTags) == 0 {
		return false
	}

	if len(includeTags) == 0 {
		return true
	}

	tagSet := make(map[string]struct{})
	for _, t := range itemTags {
		tagSet[strings.ToLower(t)] = struct{}{}
	}

	for _, inc := range includeTags {
		if _, found := tagSet[strings.ToLower(inc)]; found {
			return true
		}
	}

	return false
}
