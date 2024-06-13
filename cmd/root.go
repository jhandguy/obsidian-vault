package cmd

import (
	"fmt"
	"time"

	"github.com/jhandguy/obsidian-vault/cmd/clone"
	"github.com/jhandguy/obsidian-vault/cmd/pull"
	"github.com/jhandguy/obsidian-vault/cmd/push"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var debug bool

var cmd = &cobra.Command{
	Use:   "ov",
	Short: "CLI to backup Obsidian encrypted notes in GitHub",
	Long:  "obsidian-vault is a CLI to backup your Obsidian notes in GitHub using AES-GCM-256 authenticated encryption.",
	CompletionOptions: cobra.CompletionOptions{
		DisableDefaultCmd: true,
	},
}

func init() {
	cobra.OnInitialize(setup)

	cmd.AddCommand(clone.Cmd)
	cmd.AddCommand(pull.Cmd)
	cmd.AddCommand(push.Cmd)

	cmd.PersistentFlags().BoolVarP(&debug, "debug", "d", false, "debug for ov")
	cmd.PersistentFlags().String("path", ".", "path to the obsidian vault")
	cmd.PersistentFlags().String("config", ".obsidian", "name of the config folder")
}

func setup() {
	if err := setupLogger(); err != nil {
		fmt.Printf("failed to setup logger: %v", err)
	}
}

func setupLogger() error {
	config := zap.NewProductionConfig()
	config.Encoding = "console"
	config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	config.EncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout(time.TimeOnly)
	config.DisableStacktrace = true
	config.DisableCaller = true

	if debug {
		config.Level.SetLevel(zap.DebugLevel)
	} else {
		config.Level.SetLevel(zap.InfoLevel)
	}

	logger, err := config.Build()
	if err != nil {
		return err
	}
	zap.ReplaceGlobals(logger)
	return nil
}

func Execute(version string) {
	cmd.Version = version
	if err := cmd.Execute(); err != nil {
		zap.S().Fatalf("❌ %v", err)
	}
}
