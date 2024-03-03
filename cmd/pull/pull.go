package pull

import (
	"github.com/jhandguy/obsidian-vault/internal/vault"
	"github.com/spf13/cobra"
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
	path, err := cmd.InheritedFlags().GetString("path")
	if err != nil {
		return err
	}

	v, err := vault.New(path)
	if err != nil {
		return err
	}

	return v.Pull(password)
}
