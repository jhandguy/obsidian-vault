package clean

import (
	"github.com/jhandguy/obsidian-vault/internal/vault"
	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use:           "clean",
	Short:         "Clean and remove local vaults",
	RunE:          clean,
	SilenceUsage:  true,
	SilenceErrors: true,
}

var remove bool

func init() {
	Cmd.Flags().BoolVar(&remove, "remove", false, "should remove git vault after cleaning")
}

func clean(cmd *cobra.Command, _ []string) error {
	path, err := cmd.InheritedFlags().GetString("path")
	if err != nil {
		return err
	}

	config, err := cmd.InheritedFlags().GetString("config")
	if err != nil {
		return err
	}

	v, err := vault.New(path, config)
	if err != nil {
		return err
	}

	return v.Clean(remove)
}
