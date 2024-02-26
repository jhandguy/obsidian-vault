package push

import (
	"github.com/jhandguy/obsidian-vault/internal/env"
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

var password string

func init() {
	Cmd.Flags().StringVarP(&password, "password", "p", "", "password to encrypt the obsidian vault")
	Cmd.MarkFlagRequired("password")
}

func push(cmd *cobra.Command, _ []string) error {
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
