package gh

import (
	"fmt"
	"io"

	"github.com/jhandguy/obsidian-vault/internal/cmd"
)

type GitHub struct {
	shell string
	path  string
	name  string
}

func New(shell, path, name string) *GitHub {
	return &GitHub{shell: shell, path: path, name: name}
}

func (g *GitHub) CreateRepository(stdout, stderr io.Writer) error {
	description := fmt.Sprintf("Encrypted backup of %s, created with obsidian-vault.", g.name)
	command := fmt.Sprintf("gh repo create %s --description \"%s\" --private --disable-issues --disable-wiki", g.name, description)
	err := cmd.Run(g.shell, command, stdout, stderr)
	if err != nil {
		return fmt.Errorf("failed to create github repository: %v", err)
	}

	return nil
}

func (g *GitHub) CloneRepository(stdout, stderr io.Writer) error {
	command := fmt.Sprintf("gh repo clone %s %s", g.name, g.path)
	err := cmd.Run(g.shell, command, stdout, stderr)
	if err != nil {
		return fmt.Errorf("failed to clone github repository: %v", err)
	}

	return nil
}
