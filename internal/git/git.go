package git

import (
	"fmt"
	"io"
	"os/exec"
	"path/filepath"

	"go.uber.org/zap"
)

type Git struct {
	shell string
	path  string
}

func New(shell, path string) *Git {
	return &Git{shell: shell, path: path}
}

const HiddenFolder = ".git"

func (g *Git) Add(stdout, stderr io.Writer) error {
	folder := filepath.Join(g.path, HiddenFolder)
	command := fmt.Sprintf("git --git-dir %s --work-tree %s add .", folder, g.path)
	cmd := exec.Command(g.shell, "-c", command)
	cmd.Stdout = stdout
	cmd.Stderr = stderr

	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to add git changes: %v", err)
	}

	return nil
}

func (g *Git) Commit(stdout, stderr io.Writer, msg string) error {
	folder := filepath.Join(g.path, HiddenFolder)
	command := fmt.Sprintf("git --git-dir %s --work-tree %s commit -m \"%s\"", folder, g.path, msg)
	cmd := exec.Command(g.shell, "-c", command)
	cmd.Stdout = stdout
	cmd.Stderr = stderr

	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to commit git changes: %v", err)
	}

	return nil
}

func (g *Git) Push(stdout, stderr io.Writer) error {
	zap.S().Info("🚀 pushing vault to GitHub")

	folder := filepath.Join(g.path, HiddenFolder)
	command := fmt.Sprintf("git --git-dir %s --work-tree %s push origin main", folder, g.path)
	cmd := exec.Command(g.shell, "-c", command)
	cmd.Stdout = stdout
	cmd.Stderr = stderr

	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to push git changes: %v", err)
	}

	return nil
}

func (g *Git) Pull(stdout, stderr io.Writer) error {
	zap.S().Info("📡 pulling vault from GitHub")

	folder := filepath.Join(g.path, HiddenFolder)
	command := fmt.Sprintf("git --git-dir %s --work-tree %s pull origin main", folder, g.path)
	cmd := exec.Command(g.shell, "-c", command)
	cmd.Stdout = stdout
	cmd.Stderr = stderr

	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to pull git changes: %v", err)
	}

	return nil
}
