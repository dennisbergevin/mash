<div align="center">
<h1>️  mash  </h1>
<h5 align="center">
A customizable command launcher for storing and executing commands.
</h5>
</div>
<br>
<div align="center">
  <img alt="Mash logo" src="./assets/mash-logo.png">
</div>

## Features

- 📓 A customizable interactive list view of executable commands
- 🌲 Tree view of commands
- 🏷 Filterable list tagging

![Demo](./assets/mash-demo.gif)

#### Table of Contents

- [Installation](#installation)
- [Configuration](#configuration)
- [Command line arguments](#command-line-arguments)
  - [Help mode](#help-mode)
  - [Global](#global)
  - [Tag](#tag)
  - [Tree](#tree)
  - [Skip intro](#skip-intro)

---

## Installation

### Homebrew

```console
brew tap dennisbergevin/tools
brew install mash
```

### Go

Install with Go:

```console
go install github.com/dennisbergevin/mash@latest
```

Or grab a binary from [the latest release](https://github.com/dennisbergevin/mash/releases/latest).

---

## Configuration

Create a `.mash.json` anywhere in the directory tree (at or above the current working directory). The config file closest to the current working directory will be preferred.

This enables you to have different configs for different parent directories, such as one for a specific repository, one for personal projects, one for work scripts, etc.

For global configurations you can create a `config.json` file in the `~/.config/mash/` directory.

The content should be in the following format:

```json
{
  "tagColor": "#FFA500",
  "titleColor": "#00CED1",
  "descColor": "#888888",
  "skipIntro": false,
  "items": [
    {
      "title": "Playwright run",
      "desc": "--specs",
      "cmd": "cd ~/Projects/playwright-cli-select && npx playwright-cli-select run --specs",
      "tags": ["testing", "projects"]
    },
    {
      "title": "Say Goodbye",
      "desc": "Print Goodbye",
      "cmd": "echo Goodbye",
      "tags": ["echo"]
    },
    {
      "title": "List Files",
      "desc": "List files in current dir",
      "cmd": "ls -la"
    },
    { "title": "Current Directory", "desc": "Print working dir", "cmd": "pwd" },
    { "title": "Date", "desc": "Show current date/time", "cmd": "date" },
    {
      "title": "Go home",
      "desc": "home directory",
      "cmd": "cd:~",
      "tags": ["nav"]
    },
    { "title": "Whoami", "desc": "Show current user", "cmd": "whoami" },
    {
      "title": "Git Version",
      "desc": "Show installed Git version",
      "cmd": "git --version",
      "tags": ["git"]
    },
    {
      "title": "pwtree",
      "desc": "See an interactive tree view of Playwright suite",
      "cmd": "cd ~/Projects/playwright-cli-select && pwtree",
      "tags": ["playwright", "projects"]
    },
    {
      "title": "List Docker Containers",
      "desc": "Running containers",
      "cmd": "docker ps"
    },
    {
      "title": "List Aliases",
      "desc": "Show shell aliases",
      "cmd": "zsh -i -c alias"
    }
  ]
}
```

---

## Command line arguments

### Help mode

All available commands are included in the help menu:

```bash
mash --help
```

### Global

To display the items in the home `.config/mash/config.json` file from any directory:

```bash
mash --global
```

### Tag

To display only items that have a tag in the respective config file:

```bash
mash --tag
```

To display only items that have a tag in the global config file:

```bash
mash --tag --global
```

To display only items with specific tags, quoted and separated by semicolon:

```bash
mash --tag "dev;nav"
```

To display a tree view of items in the respective config file:

### Tree

> [!NOTE]  
> You can set the color of the tag, title, and description shown in the tree view within the config file.

```bash
mash --tree
```

To display a tree view of only tagged items in the respective config file:

```bash
mash --tree --tag
```

To display a tree view of items with specific items in the respective config file:

```bash
mash --tree --tag="dev"
```

### Skip intro

To skip the intro screen:

```bash
mash --skip-intro
```
