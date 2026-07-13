package harness

import (
	"fmt"
	"os"
	"path/filepath"
)

var (
	ServiceDirs = map[string]string{
		"api": "src/services/control-plane-api",
		"web": "src/services/control-plane-web",
		"bot": "src/discord/control-bot",
	}

	Images = map[string]string{
		"api": "control-plane-api",
		"web": "control-plane-web",
		"bot": "control-bot",
	}

	ServicePorts = map[string]int{
		"api": 8080,
		"web": 3000,
	}
)

func ProjectRoot() (string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("get working directory: %w", err)
	}

	dir := cwd
	for {
		if _, err := os.Stat(filepath.Join(dir, ".git")); err == nil {
			return dir, nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return "", fmt.Errorf("could not find project root (no .git directory found)")
		}
		dir = parent
	}
}