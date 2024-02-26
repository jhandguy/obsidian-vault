package vault

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/jhandguy/obsidian-vault/internal/crypto"
	"github.com/jhandguy/obsidian-vault/internal/git"
	"go.uber.org/zap"
)

const obsidianFolder = ".obsidian"

type Vault struct {
	directories []string
	files       []string
	sourcePath  string
	targetPath  string
}

func New(sourcePath, targetPath string) *Vault {
	zap.S().Debugf("found source vault: %s", sourcePath)
	zap.S().Debugf("found target vault: %s", targetPath)

	return &Vault{
		sourcePath: sourcePath,
		targetPath: targetPath,
	}
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

	// TODO: parallelize with goroutines
	for _, file := range v.files {
		sourceFile := filepath.Join(v.sourcePath, file)
		targetFile := filepath.Join(v.targetPath, file)

		data, err := os.ReadFile(sourceFile)
		if err != nil {
			return fmt.Errorf("failed to read file %s: %w", sourceFile, err)
		}

		encrypted, err := crypto.Encrypt(data, password, file)
		if err != nil {
			return fmt.Errorf("failed to encrypt file %s: %w", sourceFile, err)
		}

		if err := os.WriteFile(targetFile, encrypted, 0644); err != nil {
			return fmt.Errorf("failed to write file %s: %w", targetFile, err)
		}

		zap.S().Debugf("encrypted file: %s (%dB)", targetFile, len(encrypted))
	}

	return nil
}

func (v *Vault) Decrypt(password string) error {
	zap.S().Infof("🔑 decrypting vault: %s", v.targetPath)

	// TODO: parallelize with goroutines
	for _, file := range v.files {
		sourceFile := filepath.Join(v.sourcePath, file)
		targetFile := filepath.Join(v.targetPath, file)

		data, err := os.ReadFile(sourceFile)
		if err != nil {
			return fmt.Errorf("failed to read file %s: %w", sourceFile, err)
		}

		decrypted, err := crypto.Decrypt(data, password, file)
		if err != nil {
			return fmt.Errorf("failed to decrypt file %s: %w", sourceFile, err)
		}

		if err := os.WriteFile(targetFile, decrypted, 0644); err != nil {
			return fmt.Errorf("failed to write file %s: %w", targetFile, err)
		}

		zap.S().Debugf("decrypted file: %s (%dB)", targetFile, len(decrypted))
	}

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

func GetObsidianVaultPath(path string) (string, error) {
	abs, err := filepath.Abs(path)
	if err != nil {
		return "", fmt.Errorf("failed to get obsidian vault path: %w", err)
	}

	return abs, nil
}

func GetGitRepositoryPath(path string) (string, error) {
	abs, err := GetObsidianVaultPath(path)
	if err != nil {
		return "", err
	}

	return filepath.Join(abs, fmt.Sprintf(".git-%s", filepath.Base(abs))), nil
}
