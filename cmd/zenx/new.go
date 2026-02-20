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

import (
	"context"
	"log"
	"net/http"
	"time"

	"` + project + `/internal/handlers"
)

func main() {
	srv := handlers.NewServer()
	httpServer := &http.Server{
		Addr:         ":8080",
		Handler:      srv,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	go func() {
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()
	<-context.Background().Done()
}
`

	handlerContent := `package handlers

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func NewServer() http.Handler {
	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())
	mux.HandleFunc("/health", func(w http.ResponseWriter, _ *http.Request) { w.WriteHeader(http.StatusOK); _, _ = w.Write([]byte("ok")) })
	return mux
}
`

	goMod := fmt.Sprintf("module %s\n\ngo 1.22\n", project)
	writeFile(filepath.Join(project, "cmd", "server", "main.go"), mainContent)
	writeFile(filepath.Join(project, "internal", "handlers", "server.go"), handlerContent)
	writeFile(filepath.Join(project, "go.mod"), goMod)
	fmt.Printf("ZenX project created: %s\n", project)
}

func writeFile(path, content string) {
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		fmt.Printf("failed to write %s: %v\n", path, err)
	}
}
