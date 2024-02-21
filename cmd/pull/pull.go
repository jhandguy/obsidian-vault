package pull

import (
	"github.com/jhandguy/obsidian-vault/cmd/flags"
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

func pull(*cobra.Command, []string) error {
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
