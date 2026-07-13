package harness

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/fatih/color"
)

func Run(cmd string, args ...string) error {
	full := strings.Join(append([]string{cmd}, args...), " ")
	color.Cyan("  >>> %s", full)
	c := exec.Command(cmd, args...)
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	return c.Run()
}

func RunWithInput(input string, cmd string, args ...string) error {
	full := strings.Join(append([]string{cmd}, args...), " ")
	color.Cyan("  >>> %s", full)
	c := exec.Command(cmd, args...)
	c.Stdin = strings.NewReader(input)
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	return c.Run()
}

func RunQuiet(cmd string, args ...string) (string, error) {
	c := exec.Command(cmd, args...)
	out, err := c.Output()
	if err != nil {
		return "", fmt.Errorf("%s %s: %w", cmd, strings.Join(args, " "), err)
	}
	return strings.TrimSpace(string(out)), nil
}