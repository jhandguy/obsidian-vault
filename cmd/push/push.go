package push

import (
	"github.com/jhandguy/obsidian-vault/cmd/flags"
	"github.com/jhandguy/obsidian-vault/internal/vault"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

var Cmd = &cobra.Command{
	Use:           "push",
	Short:         "Encrypt and push local vault to Git",
	RunE:          push,
	SilenceUsage:  true,
	SilenceErrors: true,
}

func push(*cobra.Command, []string) error {
	shell, err := flags.GetString("shell")
	if err != nil {
		return err
	}

	password, err := flags.GetString("password")
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

	v := vault.New(source, target)

	if err = v.Scan(); err != nil {
		return err
	}

	if err = v.Clean(); err != nil {
		return err
	}

	if err = v.Encrypt(password); err != nil {
		return err
	}

	if err = v.Push(shell); err != nil {
		return err
	}

	zap.S().Info("✅ vault backup successful")

	return nil
}
