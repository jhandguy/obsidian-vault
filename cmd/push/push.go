package push

import (
	"github.com/jhandguy/obsidian-vault/internal/vault"
	"github.com/spf13/cobra"
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
	path, err := cmd.InheritedFlags().GetString("path")
	if err != nil {
		return err
	}

	v, err := vault.New(path)
	if err != nil {
		return err
	}

	return v.Push(password)
}
