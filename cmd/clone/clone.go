package clone

import (
	"github.com/jhandguy/obsidian-vault/internal/vault"
	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use:           "clone",
	Short:         "Create and clone private GitHub repository",
	RunE:          clone,
	SilenceUsage:  true,
	SilenceErrors: true,
}

var create bool

func init() {
	Cmd.Flags().BoolVar(&create, "create", false, "should create repository before cloning")
}

func clone(cmd *cobra.Command, _ []string) error {
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

	return v.Clone(create)
}
