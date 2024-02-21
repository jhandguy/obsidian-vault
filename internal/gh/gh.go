package gh

import (
	"fmt"
	"os/exec"

	"go.uber.org/zap"
)

func CreateRepository(shell, name string) error {
	zap.S().Info("🐙 creating github repository")

	description := fmt.Sprintf("Encrypted backup of %s, created with obsidian-vault.", name)
	command := fmt.Sprintf("gh repo create %s --description \"%s\" --private --disable-issues --disable-wiki", name, description)
	cmd := exec.Command(shell, "-c", command)

	bytes, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to create github repository: %s", bytes)
	}

	return nil
}

func CloneRepository(shell, name, path string) error {
	zap.S().Info("📦 cloning github repository")

	command := fmt.Sprintf("gh repo clone %s %s", name, path)
	cmd := exec.Command(shell, "-c", command)

	bytes, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to clone github repository: %s", bytes)
	}

	return nil
}
