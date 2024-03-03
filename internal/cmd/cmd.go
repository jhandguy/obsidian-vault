package cmd

import (
	"io"
	"os/exec"

	"go.uber.org/zap"
)

func Run(shell, command string, stdout, stderr io.Writer) error {
	zap.S().Debug(command)
	cmd := exec.Command(shell, "-c", command)
	cmd.Stdout = stdout
	cmd.Stderr = stderr
	return cmd.Run()
}
