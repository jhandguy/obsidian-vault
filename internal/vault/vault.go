package vault

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/jhandguy/obsidian-vault/internal/crypto"
	"github.com/jhandguy/obsidian-vault/internal/gh"
	"github.com/jhandguy/obsidian-vault/internal/git"
	"go.uber.org/zap"
)

type Vault struct {
	directories []string
	files       []string
	sourcePath  string
	targetPath  string
}

type Type string

const (
	obsidianFolder string = ".obsidian"
	LocalVaultType Type   = "local"
	GitVaultType   Type   = "git"
)

func New(path string, vaultType Type) (*Vault, error) {
	v := &Vault{}

	localPath, err := filepath.Abs(path)
	if err != nil {
		return nil, fmt.Errorf("failed to get local vault path: %w", err)
	}

	gitPath := filepath.Join(localPath, fmt.Sprintf(".git-%s", filepath.Base(localPath)))

	switch vaultType {
	case LocalVaultType:
		v.sourcePath = localPath
		v.targetPath = gitPath
	case GitVaultType:
		v.sourcePath = gitPath
		v.targetPath = localPath
	default:
		return nil, fmt.Errorf("unknown vault type: %s", vaultType)
	}

	zap.S().Debugf("vault type: %s", vaultType)
	zap.S().Debugf("source vault path: %s", v.sourcePath)
	zap.S().Debugf("target vault path: %s", v.targetPath)

	return v, nil
}

func (v *Vault) Scan() error {
	if !v.isObsidianVault() {
		return fmt.Errorf("not an obsidian vault: %s", v.sourcePath)
	}

	fn := func(path string, d fs.DirEntry, e error) error {
		if v.sourcePath == path {
			return nil
		}

		if strings.HasPrefix(d.Name(), ".git") {
			return filepath.SkipDir
		}

		relativePath, err := filepath.Rel(v.sourcePath, path)
		if err != nil {
			return err
		}

		if d.IsDir() {
			v.directories = append(v.directories, relativePath)
		} else {
			v.files = append(v.files, relativePath)
		}

		return nil
	}
	if err := filepath.WalkDir(v.sourcePath, fn); err != nil {
		return fmt.Errorf("failed to scan vault: %w", err)
	}

	zap.S().Debugf("scanned %d directories: %v", len(v.directories), v.directories)
	zap.S().Debugf("scanned %d files: %v", len(v.files), v.files)

	return nil
}

func (v *Vault) isObsidianVault() bool {
	_, err := os.Stat(filepath.Join(v.sourcePath, obsidianFolder))
	return !os.IsNotExist(err)
}

func (v *Vault) Clean() error {
	fn := func(path string, d fs.DirEntry, e error) error {
		if v.targetPath == path {
			return nil
		}

		if strings.HasPrefix(d.Name(), ".") {
			return filepath.SkipDir
		}

		if d.IsDir() {
			if err := os.RemoveAll(path); err != nil {
				return fmt.Errorf("failed to remove directory %s: %w", path, err)
			}
			return filepath.SkipDir
		}

		if err := os.Remove(path); err != nil {
			return fmt.Errorf("failed to remove file %s: %w", path, err)
		}

		return nil
	}

	if err := filepath.WalkDir(v.targetPath, fn); err != nil {
		return fmt.Errorf("failed to clean vault: %w", err)
	}

	for _, dir := range v.directories {
		dirPath := filepath.Join(v.targetPath, dir)

		if err := os.MkdirAll(dirPath, os.ModePerm); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dirPath, err)
		}
	}

	return nil
}

func (v *Vault) Encrypt(password string) error {
	zap.S().Infof("🔒 encrypting vault: %s", v.targetPath)

	channel := make(chan error, len(v.files))

	for _, fileName := range v.files {
		go func(fileName string) {
			channel <- v.encryptFile(fileName, password)
		}(fileName)
	}

	for range v.files {
		if err := <-channel; err != nil {
			return err
		}
	}

	return nil
}

func (v *Vault) encryptFile(fileName, password string) error {
	sourceFile := filepath.Join(v.sourcePath, fileName)
	targetFile := filepath.Join(v.targetPath, fileName)

	data, err := os.ReadFile(sourceFile)
	if err != nil {
		return fmt.Errorf("failed to read file %s: %w", sourceFile, err)
	}

	encrypted, err := crypto.Encrypt(data, password, fileName)
	if err != nil {
		return fmt.Errorf("failed to encrypt file %s: %w", sourceFile, err)
	}

	if err := os.WriteFile(targetFile, encrypted, 0644); err != nil {
		return fmt.Errorf("failed to write file %s: %w", targetFile, err)
	}

	zap.S().Debugf("encrypted file: %s (%dB)", targetFile, len(encrypted))

	return nil
}

func (v *Vault) Decrypt(password string) error {
	zap.S().Infof("🔑 decrypting vault: %s", v.targetPath)

	channel := make(chan error, len(v.files))

	for _, fileName := range v.files {
		go func(fileName string) {
			channel <- v.decryptFile(fileName, password)
		}(fileName)
	}

	for range v.files {
		if err := <-channel; err != nil {
			return err
		}
	}

	return nil
}

func (v *Vault) decryptFile(fileName, password string) error {
	sourceFile := filepath.Join(v.sourcePath, fileName)
	targetFile := filepath.Join(v.targetPath, fileName)

	data, err := os.ReadFile(sourceFile)
	if err != nil {
		return fmt.Errorf("failed to read file %s: %w", sourceFile, err)
	}

	decrypted, err := crypto.Decrypt(data, password, fileName)
	if err != nil {
		return fmt.Errorf("failed to decrypt file %s: %w", sourceFile, err)
	}

	if err := os.WriteFile(targetFile, decrypted, 0644); err != nil {
		return fmt.Errorf("failed to write file %s: %w", targetFile, err)
	}

	zap.S().Debugf("decrypted file: %s (%dB)", targetFile, len(decrypted))

	return nil
}

func (v *Vault) Push(shell string) error {
	zap.S().Info("🚀 pushing vault to GitHub")

	if err := git.Add(shell, v.targetPath); err != nil {
		return err
	}

	msg := fmt.Sprintf("[%s] obsidian-vault backup", time.Now().Format(time.DateTime))
	if err := git.Commit(shell, v.targetPath, msg); err != nil {
		return err
	}

	return git.Push(shell, v.targetPath)
}

func (v *Vault) Pull(shell string) error {
	zap.S().Info("📡 pulling vault from GitHub")

	return git.Pull(shell, v.sourcePath)
}

func (v *Vault) Clone(shell string, create bool) error {
	name := filepath.Base(v.sourcePath)

	if create {
		if err := gh.CreateRepository(shell, name); err != nil {
			return err
		}
	}

	if err := gh.CloneRepository(shell, name, v.targetPath); err != nil {
		return err
	}

	return nil
}
