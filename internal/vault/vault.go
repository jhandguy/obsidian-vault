package vault

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/jhandguy/obsidian-vault/internal/crypto"
	"github.com/jhandguy/obsidian-vault/internal/env"
	"github.com/jhandguy/obsidian-vault/internal/gh"
	"github.com/jhandguy/obsidian-vault/internal/git"
	"go.uber.org/zap"
)

type Vault struct {
	directories []string
	files       []string
	localPath   string
	gitPath     string
}

type vaultType string

const (
	vaultTypeLocal vaultType = "local"
	vaultTypeGit   vaultType = "git"
)

func New(path string) (*Vault, error) {
	localPath, err := filepath.Abs(path)
	if err != nil {
		return nil, fmt.Errorf("failed to get local vault path: %w", err)
	}

	gitPath := filepath.Join(localPath, fmt.Sprintf(".git-%s", filepath.Base(localPath)))

	zap.S().Debugf("local vault path: %s", localPath)
	zap.S().Debugf("git vault path: %s", gitPath)

	return &Vault{
		localPath: localPath,
		gitPath:   gitPath,
	}, nil
}

func (v *Vault) Clone(create bool) error {
	shell := env.GetShell()
	name := filepath.Base(v.localPath)

	if create {
		if err := gh.CreateRepository(shell, name); err != nil {
			return err
		}
	}

	if err := gh.CloneRepository(shell, name, v.gitPath); err != nil {
		return err
	}

	return nil
}

func (v *Vault) Pull(password string) error {
	shell := env.GetShell()
	if err := git.Pull(shell, v.gitPath); err != nil {
		return err
	}

	if err := v.scan(vaultTypeGit); err != nil {
		return err
	}

	if err := v.clean(vaultTypeLocal); err != nil {
		return err
	}

	if err := v.decrypt(password); err != nil {
		return err
	}

	return nil
}

func (v *Vault) Push(password string) error {
	shell := env.GetShell()
	if err := v.scan(vaultTypeLocal); err != nil {
		return err
	}

	if err := v.clean(vaultTypeGit); err != nil {
		return err
	}

	if err := v.encrypt(password); err != nil {
		return err
	}

	if err := git.Add(shell, v.gitPath); err != nil {
		return err
	}

	msg := fmt.Sprintf("[%s] obsidian-vault backup", time.Now().Format(time.DateTime))
	if err := git.Commit(shell, v.gitPath, msg); err != nil {
		return err
	}

	return git.Push(shell, v.gitPath)
}

func (v *Vault) scan(t vaultType) error {
	path, err := v.getVaultPath(t)
	if err != nil {
		return err
	}

	if _, err := os.Stat(filepath.Join(path, ".obsidian")); os.IsNotExist(err) {
		return fmt.Errorf("not an obsidian vault: %s", path)
	}

	fn := func(p string, d fs.DirEntry, e error) error {
		if p == path {
			return nil
		}

		if strings.HasPrefix(d.Name(), ".git") {
			return filepath.SkipDir
		}

		relativePath, err := filepath.Rel(path, p)
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
	if err := filepath.WalkDir(path, fn); err != nil {
		return fmt.Errorf("failed to scan vault: %w", err)
	}

	zap.S().Debugf("scanned %d directories: %v", len(v.directories), v.directories)
	zap.S().Debugf("scanned %d files: %v", len(v.files), v.files)

	return nil
}

func (v *Vault) clean(t vaultType) error {
	path, err := v.getVaultPath(t)
	if err != nil {
		return err
	}

	fn := func(p string, d fs.DirEntry, e error) error {
		if p == path {
			return nil
		}

		if strings.HasPrefix(d.Name(), ".") {
			return filepath.SkipDir
		}

		if d.IsDir() {
			if err := os.RemoveAll(p); err != nil {
				return fmt.Errorf("failed to remove directory %s: %w", p, err)
			}
			return filepath.SkipDir
		}

		if err := os.Remove(p); err != nil {
			return fmt.Errorf("failed to remove file %s: %w", p, err)
		}

		return nil
	}

	if err := filepath.WalkDir(path, fn); err != nil {
		return fmt.Errorf("failed to clean vault: %w", err)
	}

	for _, dir := range v.directories {
		dirPath := filepath.Join(path, dir)

		if err := os.MkdirAll(dirPath, os.ModePerm); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dirPath, err)
		}
	}

	return nil
}

func (v *Vault) getVaultPath(t vaultType) (string, error) {
	switch t {
	case vaultTypeLocal:
		return v.localPath, nil
	case vaultTypeGit:
		return v.gitPath, nil
	default:
		return "", fmt.Errorf("unknown vault type: %s", t)
	}
}

func (v *Vault) encrypt(password string) error {
	zap.S().Infof("🔒 encrypting vault: %s", v.localPath)

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
	localFile := filepath.Join(v.localPath, fileName)
	gitFile := filepath.Join(v.gitPath, fileName)

	data, err := os.ReadFile(localFile)
	if err != nil {
		return fmt.Errorf("failed to read file %s: %w", localFile, err)
	}

	encrypted, err := crypto.Encrypt(data, password, fileName)
	if err != nil {
		return fmt.Errorf("failed to encrypt file %s: %w", localFile, err)
	}

	if err := os.WriteFile(gitFile, encrypted, 0644); err != nil {
		return fmt.Errorf("failed to write file %s: %w", gitFile, err)
	}

	zap.S().Debugf("encrypted file: %s (%dB)", gitFile, len(encrypted))

	return nil
}

func (v *Vault) decrypt(password string) error {
	zap.S().Infof("🔑 decrypting vault: %s", v.gitPath)

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
	gitFile := filepath.Join(v.gitPath, fileName)
	localFile := filepath.Join(v.localPath, fileName)

	data, err := os.ReadFile(gitFile)
	if err != nil {
		return fmt.Errorf("failed to read file %s: %w", gitFile, err)
	}

	decrypted, err := crypto.Decrypt(data, password, fileName)
	if err != nil {
		return fmt.Errorf("failed to decrypt file %s: %w", gitFile, err)
	}

	if err := os.WriteFile(localFile, decrypted, 0644); err != nil {
		return fmt.Errorf("failed to write file %s: %w", localFile, err)
	}

	zap.S().Debugf("decrypted file: %s (%dB)", localFile, len(decrypted))

	return nil
}
