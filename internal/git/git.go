package git

import (
	"fmt"
	"io"
	"path/filepath"

	"github.com/jhandguy/obsidian-vault/internal/cmd"
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
	err := cmd.Run(g.shell, command, stdout, stderr)
	if err != nil {
		return fmt.Errorf("failed to add git changes: %v", err)
	}

	return nil
}

func (g *Git) Commit(stdout, stderr io.Writer, msg string) error {
	folder := filepath.Join(g.path, HiddenFolder)
	command := fmt.Sprintf("git --git-dir %s --work-tree %s commit -m \"%s\"", folder, g.path, msg)
	err := cmd.Run(g.shell, command, stdout, stderr)
	if err != nil {
		return fmt.Errorf("failed to commit git changes: %v", err)
	}

	return nil
}

func (g *Git) Push(stdout, stderr io.Writer) error {
	folder := filepath.Join(g.path, HiddenFolder)
	command := fmt.Sprintf("git --git-dir %s --work-tree %s push origin main", folder, g.path)
	err := cmd.Run(g.shell, command, stdout, stderr)
	if err != nil {
		return fmt.Errorf("failed to push git changes: %v", err)
	}

	return nil
}

func (g *Git) Pull(stdout, stderr io.Writer) error {
	folder := filepath.Join(g.path, HiddenFolder)
	command := fmt.Sprintf("git --git-dir %s --work-tree %s pull origin main", folder, g.path)
	err := cmd.Run(g.shell, command, stdout, stderr)
	if err != nil {
		return fmt.Errorf("failed to pull git changes: %v", err)
	}

	return nil
}
