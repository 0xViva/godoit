# 📝 GoDoIt - ToDo for your terminal

A simple ToDo TUI written in Go using Bubbletea.

## Features

- ✅ **Add** new tasks  
- ✏️ **Edit** existing tasks inline  
- ❌ **Delete** tasks  
- ✔️ **Mark** tasks as done/undone  
- 🔤 **Navigate** tasks with `j/k` or arrow keys  
- 💾 **Persistent storage**: Saves tasks to:`$XDG_DATA_HOME/godoit/todos.json` or `/home/<user>/.local/share/godoit/todos.json`

## Installation

> **Prerequisite:** Go must be installed on your system. You can download it from [golang.org](https://golang.org/dl/).

### 1. Install using Go directly (recommended)

If you have Go installed, you can install directly with:

```zsh
go install github.com/0xViva/godoit@latest
```
### 2. Install using Go from source:

```zsh
git clone https://github.com/0xViva/godoit.git
cd godoit
go install
```

### 3. Install by downloading from releases:

You can download precompiled binaries from [Releases page](https://github.com/0xViva/godoit/releases)

## Run

>⚠️ Important: Make sure your Go bin directory is in your shell $PATH.
>For most systems, add this to your .zshrc or .bashrc:

```zsh
export PATH="$PATH:$(go env GOPATH)/bin"
```

Once installed and your `$PATH` includes your `GOPATH/bin` or `GOBIN`, you can run GoDoIt simply by typing:

```zsh
godoit
```

## Wanna contribute to development?

> prerequisites:

- https://taskfile.dev/
- https://github.com/air-verse/air

### run locally with hotreloading:
`air`

### A quick cleanup if needed:
If you want to remove GoDoIt from your system, including the binary and the saved `todos.json` file, you can run:
```zsh
go-task remove
```

#### Different cmds for development, building and release can be found in TaskFile.yml. .air.toml is for hotreload behavior.
