package clone

import (
	"path/filepath"

	"github.com/jhandguy/obsidian-vault/internal/env"
	"github.com/jhandguy/obsidian-vault/internal/gh"
	"github.com/jhandguy/obsidian-vault/internal/vault"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

var Cmd = &cobra.Command{
	Use:           "clone",
	Short:         "Create and clone private GitHub repository",
	RunE:          clone,
	SilenceUsage:  true,
	SilenceErrors: true,
}

var shouldCreateRepository bool

func init() {
	Cmd.Flags().BoolVar(&shouldCreateRepository, "create", false, "should create repository before cloning")
}

func clone(cmd *cobra.Command, _ []string) error {
	shell := env.GetShell()

	path, err := cmd.InheritedFlags().GetString("path")
	if err != nil {
		return err
	}

	source, err := vault.GetObsidianVaultPath(path)
	if err != nil {
		return err
	}

	target, err := vault.GetGitRepositoryPath(path)
	if err != nil {
		return err
	}

	name := filepath.Base(source)

	if shouldCreateRepository {
		if err := gh.CreateRepository(shell, name); err != nil {
			return err
		}
	}

	if err := gh.CloneRepository(shell, name, target); err != nil {
		return err
	}

	zap.S().Info("✅ github clone successful")

	return nil
}
