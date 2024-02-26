package pull

import (
	"github.com/jhandguy/obsidian-vault/internal/env"
	"github.com/jhandguy/obsidian-vault/internal/vault"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

var Cmd = &cobra.Command{
	Use:           "pull",
	Short:         "Pull and decrypt remote vault from Git",
	RunE:          pull,
	SilenceUsage:  true,
	SilenceErrors: true,
}

var password string

func init() {
	Cmd.Flags().StringVarP(&password, "password", "p", "", "password to decrypt the obsidian vault")
	Cmd.MarkFlagRequired("password")
}

func pull(cmd *cobra.Command, _ []string) error {
	shell := env.GetShell()

	path, err := cmd.InheritedFlags().GetString("path")
	if err != nil {
		return err
	}

	source, err := vault.GetGitRepositoryPath(path)
	if err != nil {
		return err
	}

	target, err := vault.GetObsidianVaultPath(path)
	if err != nil {
		return err
	}

	v := vault.New(source, target)

	if err = v.Pull(shell); err != nil {
		return err
	}

	if err = v.Scan(); err != nil {
		return err
	}

	if err = v.Clean(); err != nil {
		return err
	}

	if err = v.Decrypt(password); err != nil {
		return err
	}

	zap.S().Info("✅ vault sync successful")

	return nil
}
