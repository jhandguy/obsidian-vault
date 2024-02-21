package clone

import (
	"path/filepath"

	"github.com/jhandguy/obsidian-vault/cmd/flags"
	"github.com/jhandguy/obsidian-vault/internal/gh"
	"github.com/jhandguy/obsidian-vault/internal/vault"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

var Cmd = &cobra.Command{
	Use:           "clone",
	Short:         "Create and clone private GitHub repository",
	RunE:          clone,
	SilenceUsage:  true,
	SilenceErrors: true,
}

func init() {
	bindBoolPFlag("create", false, "should create repository before cloning")
}

func bindBoolPFlag(name string, value bool, usage string) {
	Cmd.PersistentFlags().Bool(name, value, usage)
	if err := viper.BindPFlag(name, Cmd.PersistentFlags().Lookup(name)); err != nil {
		zap.S().Fatalf("failed to bind bool persistent flag '%s': %v", name, err)
	}
}

func clone(*cobra.Command, []string) error {
	shell, err := flags.GetString("shell")
	if err != nil {
		return err
	}

	path, err := flags.GetString("path")
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

	if viper.GetBool("create") {
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
