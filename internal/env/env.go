package env

import "os"

func GetShell() string {
	shell, ok := os.LookupEnv("SHELL")
	if !ok {
		return "/bin/sh"
	}

	return shell
}
