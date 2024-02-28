package git

import (
	"fmt"
	"os/exec"
	"path/filepath"

	"go.uber.org/zap"
)

const GitFolder = ".git"

func Add(shell, path string) error {
	folder := filepath.Join(path, GitFolder)
	command := fmt.Sprintf("git --git-dir %s --work-tree %s add .", folder, path)
	cmd := exec.Command(shell, "-c", command)

	bytes, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to add git changes: %s", bytes)
	}

	return nil
}

func Commit(shell, path, msg string) error {
	folder := filepath.Join(path, GitFolder)
	command := fmt.Sprintf("git --git-dir %s --work-tree %s commit -m \"%s\"", folder, path, msg)
	cmd := exec.Command(shell, "-c", command)

	bytes, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to commit git changes: %s", bytes)
	}

	return nil
}

func Push(shell, path string) error {
	zap.S().Info("🚀 pushing vault to GitHub")

	folder := filepath.Join(path, GitFolder)
	command := fmt.Sprintf("git --git-dir %s --work-tree %s push origin main", folder, path)
	cmd := exec.Command(shell, "-c", command)

	bytes, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to push git changes: %s", bytes)
	}

	return nil
}

func Pull(shell, path string) error {
	zap.S().Info("📡 pulling vault from GitHub")

	folder := filepath.Join(path, GitFolder)
	command := fmt.Sprintf("git --git-dir %s --work-tree %s pull origin main", folder, path)
	cmd := exec.Command(shell, "-c", command)

	bytes, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to pull git changes: %s", bytes)
	}

	return nil
}
