package plugins

import (
	"context"
	"os/exec"
)

// WASM runtime integration point (e.g. wazero/wasmtime wrapper).
type WASMRunner struct{}

func (WASMRunner) Run(_ context.Context, modulePath string, args ...string) error {
	cmd := exec.Command(modulePath, args...)
	return cmd.Run()
}
