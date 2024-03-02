# obsidian-vault

obsidian-vault is a CLI to backup your Obsidian notes in GitHub using AES-GCM-256 authenticated encryption.

## Requirements

The following CLIs are required to use obsidian-vault:

- [GitHub](https://cli.github.com/)
- [Git](https://git-scm.com/)

## Usage

```shell
➜ ov
obsidian-vault is a CLI to backup your Obsidian notes in GitHub using AES-GCM-256 authenticated encryption.

Usage:
  ov [command]

Available Commands:
  clone       Create and clone private GitHub repository
  help        Help about any command
  pull        Pull and decrypt remote vault from Git
  push        Encrypt and push local vault to Git

Flags:
  -h, --help          help for ov
      --path string   path to the obsidian vault (default ".")

Use "ov [command] --help" for more information about a command.
```
