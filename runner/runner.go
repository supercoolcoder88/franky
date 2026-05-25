package runner

import (
	_ "embed"
	"os"
	"os/exec"
)

type runner struct {
	projectPath string
}

func NewRunner(p string) *runner {
	return &runner{
		projectPath: p,
	}
}

func (r *runner) Execute() error {
	cmd := exec.Command("docker", "run",
		"--rm", "-i",
		"-v", r.projectPath+":/workspace",
		"-w", "/workspace",
		"golang:1.26",
		"sh",
	)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}
