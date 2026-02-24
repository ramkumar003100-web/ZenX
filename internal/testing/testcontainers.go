package testing

import (
	"context"
	"os/exec"
)

func StartContainer(ctx context.Context, image, name string, args ...string) error {
	all := append([]string{"run", "-d", "--rm", "--name", name}, args...)
	all = append(all, image)
	cmd := exec.CommandContext(ctx, "docker", all...)
	return cmd.Run()
}
