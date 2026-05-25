package runner

import (
	_ "embed"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

type Config struct {
	RunBuildImage bool `json:"run_build_image"`
	RunLint       bool `json:"run_lints"`
	RunTests      bool `json:"run_tests"`
	RunVulnCheck  bool `json:"run_vulnerability_check"`
}

type runner struct {
	config *Config
}

func NewRunner(c *Config) *runner {
	return &runner{
		config: c,
	}
}

func (r *runner) Execute() error {
	// Run static checks first
	if r.config.RunLint || r.config.RunVulnCheck {
		if err := r.executeStaticChecks(); err != nil {
			return fmt.Errorf("static checks failed: %w", err)
		}
	}

	c := r.createCommands()
	if c == "" {
		return nil
	}

	wd, _ := os.Getwd()

	cmd := exec.Command("docker", "run",
		"--rm", "-i",
		"-v", wd+":/workspace",
		"-w", "/workspace",
		"golang:1.26",
		"sh",
	)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	stdin, err := cmd.StdinPipe()
	if err != nil {
		return fmt.Errorf("creating stdin pipe: %w", err)
	}

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("starting container: %w", err)
	}

	if _, err := stdin.Write([]byte(c)); err != nil {
		return fmt.Errorf("writing to stdin: %w", err)
	}

	stdin.Close()

	return cmd.Wait()
}

func (r *runner) executeStaticChecks() error {
	if r.config.RunLint {
		cmd := exec.Command("golangci-lint", "run", "./...")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("lint failed: %w", err)
		}

		fmt.Print("Lint success")
	}

	if r.config.RunVulnCheck {
		cmd := exec.Command("govulncheck", "./...")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("vulnerability check failed: %w", err)
		}
	}

	return nil
}

func (r *runner) createCommands() string {
	commands := []string{"go build ./..."}

	if r.config.RunTests {
		commands = append(commands, "go test ./...")
	}

	return strings.Join(commands, " && ")
}
