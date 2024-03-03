package vault

import (
	"fmt"
	"io"
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
	localPath   string
	gitPath     string
	gh          *gh.GitHub
	git         *git.Git
	crypto      *crypto.Crypto
	stdout      io.Writer
	stderr      io.Writer
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
	repoName := filepath.Base(localPath)
	shell := getShell()

	zap.S().Debugf("local vault path: %s", localPath)
	zap.S().Debugf("git vault path: %s", gitPath)

	var stdout, stderr io.Writer
	if zap.S().Level() == zap.DebugLevel {
		stdout = os.Stdout
		stderr = os.Stderr
	}

	return &Vault{
		localPath: localPath,
		gitPath:   gitPath,
		gh:        gh.New(shell, gitPath, repoName),
		git:       git.New(shell, gitPath),
		crypto:    crypto.New(),
		stdout:    stdout,
		stderr:    stderr,
	}, nil
}

func (v *Vault) Clone(create bool) error {
	if create {
		zap.S().Info("🐙 creating github repository")
		if err := v.gh.CreateRepository(v.stdout, v.stderr); err != nil {
			return err
		}
	}

	zap.S().Info("📦 cloning github repository")
	if err := v.gh.CloneRepository(v.stdout, v.stderr); err != nil {
		return err
	}

	zap.S().Info("✅ github clone successful")
	return nil
}

func (v *Vault) Pull(password string) error {
	zap.S().Info("📡 pulling vault from GitHub")
	if err := v.git.Pull(v.stdout, v.stderr); err != nil {
		return err
	}

	if err := v.scan(vaultTypeGit); err != nil {
		return err
	}

	if err := v.clean(vaultTypeLocal); err != nil {
		return err
	}

	zap.S().Infof("🔑 decrypting vault: %s", v.gitPath)
	if err := v.decrypt(password); err != nil {
		return err
	}

	zap.S().Info("✅ vault sync successful")
	return nil
}

func (v *Vault) Push(password string) error {
	if err := v.scan(vaultTypeLocal); err != nil {
		return err
	}

	if err := v.clean(vaultTypeGit); err != nil {
		return err
	}

	zap.S().Infof("🔒 encrypting vault: %s", v.localPath)
	if err := v.encrypt(password); err != nil {
		return err
	}

	zap.S().Info("🚀 pushing vault to GitHub")
	if err := v.git.Add(v.stdout, v.stderr); err != nil {
		return err
	}

	msg := fmt.Sprintf("[%s] obsidian-vault backup", time.Now().Format(time.DateTime))
	if err := v.git.Commit(v.stdout, v.stderr, msg); err != nil {
		return err
	}

	if err := v.git.Push(v.stdout, v.stderr); err != nil {
		return err
	}

	zap.S().Info("✅ vault backup successful")
	return nil
}

func (v *Vault) scan(t vaultType) error {
	path, err := v.getVaultPath(t)
	if err != nil {
		return err
	}

	if _, err := os.Stat(filepath.Join(path, ".obsidian")); os.IsNotExist(err) {
		return fmt.Errorf("not an obsidian vault: %s", path)
	}

	v.directories = []string{}
	v.files = []string{}
	fn := func(p string, d fs.DirEntry, _ error) error {
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

	fn := func(p string, d fs.DirEntry, _ error) error {
		if p == path {
			return nil
		}

		if d.Name() == git.HiddenFolder || p == v.gitPath {
			return filepath.SkipDir
		}

		if d.IsDir() {
			if err := os.RemoveAll(p); err != nil {
				return fmt.Errorf("failed to remove directory %s: %w", p, err)
			}

			zap.S().Debugf("removed directory: %s", p)
			return filepath.SkipDir
		}

		if err := os.Remove(p); err != nil {
			return fmt.Errorf("failed to remove file %s: %w", p, err)
		}

		zap.S().Debugf("removed file: %s", p)
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

		zap.S().Debugf("created directory: %s", dirPath)
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

	encrypted, err := v.crypto.Encrypt(data, password, fileName)
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

	decrypted, err := v.crypto.Decrypt(data, password, fileName)
	if err != nil {
		return fmt.Errorf("failed to decrypt file %s: %w", gitFile, err)
	}

	if err := os.WriteFile(localFile, decrypted, 0644); err != nil {
		return fmt.Errorf("failed to write file %s: %w", localFile, err)
	}

	zap.S().Debugf("decrypted file: %s (%dB)", localFile, len(decrypted))
	return nil
}

func getShell() string {
	shell, ok := os.LookupEnv("SHELL")
	if !ok {
		return "/bin/sh"
	}

	return shell
}
