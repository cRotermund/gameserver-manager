package harness

import (
	"fmt"
	"os/exec"
)

func ContainerCLI() (string, error) {
	for _, candidate := range []string{"podman", "docker"} {
		if _, err := exec.LookPath(candidate); err == nil {
			return candidate, nil
		}
	}
	return "", fmt.Errorf("neither podman nor docker is installed or available on PATH")
}