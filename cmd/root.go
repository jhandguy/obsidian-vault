package cmd

import (
	"log"
	"strings"
	"time"

	"github.com/jhandguy/obsidian-vault/cmd/clone"
	"github.com/jhandguy/obsidian-vault/cmd/pull"
	"github.com/jhandguy/obsidian-vault/cmd/push"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var cmd = &cobra.Command{
	Use:   "ov",
	Short: "CLI to backup Obsidian encrypted notes in GitHub",
	Long:  "obsidian-vault is a CLI to backup your Obsidian notes in GitHub with AES-GCM-256 authenticated encryption.",
	CompletionOptions: cobra.CompletionOptions{
		DisableDefaultCmd: true,
	},
}

func init() {
	cobra.OnInitialize(setup)

	viper.AutomaticEnv()
	viper.SetEnvPrefix("ov")
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	viper.MustBindEnv("shell", "SHELL")

	cmd.AddCommand(clone.Cmd)
	cmd.AddCommand(pull.Cmd)
	cmd.AddCommand(push.Cmd)

	bindStringFlag("path", ".", "path to the obsidian vault")
	bindStringPFlag("password", "p", "", "password to encrypt/decrypt the obsidian vault")
}

func setup() {
	if err := setupLogger(); err != nil {
		log.Fatalf("failed to setup logger: %v", err)
	}
}

func setupLogger() error {
	config := zap.NewProductionConfig()
	config.Encoding = "console"
	config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	config.EncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout(time.TimeOnly)
	config.DisableStacktrace = true
	config.DisableCaller = true
	logger, err := config.Build()
	if err != nil {
		return err
	}
	zap.ReplaceGlobals(logger)
	return nil
}

func bindStringFlag(name, value, usage string) {
	cmd.PersistentFlags().String(name, value, usage)
	if err := viper.BindPFlag(name, cmd.PersistentFlags().Lookup(name)); err != nil {
		zap.S().Fatalf("failed to bind string persistent flag '%s': %v", name, err)
	}
}

func bindStringPFlag(name, shorthand, value, usage string) {
	cmd.PersistentFlags().StringP(name, shorthand, value, usage)
	if err := viper.BindPFlag(name, cmd.PersistentFlags().Lookup(name)); err != nil {
		zap.S().Fatalf("failed to bind string persistent flag '%s': %v", name, err)
	}
}

func Execute() {
	if err := cmd.Execute(); err != nil {
		zap.S().Fatal(err)
	}
}
