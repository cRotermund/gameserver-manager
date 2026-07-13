package harness

import (
	"fmt"
	"os/exec"
)

func KustomizeAvailable() bool {
	return exec.Command("kustomize", "version").Run() == nil
}

func KustomizeBuild(dir string) (string, error) {
	cmd := exec.Command("kustomize", "build", dir)
	out, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("kustomize build %s: %w", dir, err)
	}
	return string(out), nil
}