package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func RunNew(args []string) {
	if len(args) < 1 {
		fmt.Println("Usage: zenx new <project>")
		return
	}

	project := strings.TrimSpace(args[0])
	if project == "" {
		fmt.Println("project name cannot be empty")
		return
	}

	dirs := []string{
		project,
		filepath.Join(project, "cmd", "server"),
		filepath.Join(project, "internal", "handlers"),
		filepath.Join(project, "internal", "services"),
		filepath.Join(project, "internal", "repositories"),
		filepath.Join(project, "configs"),
		filepath.Join(project, "migrations"),
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			fmt.Printf("failed to create %s: %v\n", dir, err)
			return
		}
	}

	mainContent := `package main

import "fmt"

func main() {
	fmt.Println("Welcome to ZenX")
}
`

	goMod := fmt.Sprintf("module %s\n\ngo 1.22\n", project)

	writeFile(filepath.Join(project, "cmd", "server", "main.go"), mainContent)
	writeFile(filepath.Join(project, "go.mod"), goMod)
	fmt.Printf("ZenX project created: %s\n", project)
}

func writeFile(path, content string) {
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		fmt.Printf("failed to write %s: %v\n", path, err)
	}
}
