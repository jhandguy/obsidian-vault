package vault

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExample(t *testing.T) {
	pwd, err := os.Getwd()
	assert.Nil(t, err, "expected no error when getting current working directory")
	path := filepath.Join(pwd, "../../example")
	password := "consectetur-adipiscing-elit"

	err = os.Setenv("SHELL", "echo")
	assert.Nil(t, err, "expected no error when setting SHELL environment variable")

	v, err := New(path)
	assert.Nil(t, err, "expected no error when creating vault")

	gitPath, err := v.getVaultPath(vaultTypeGit)
	assert.Nil(t, err, "expected no error when getting git vault path")

	err = os.MkdirAll(gitPath, os.ModePerm)
	assert.Nil(t, err, "expected no error when creating git vault folder")
	defer os.RemoveAll(gitPath)

	err = v.Push(password)
	assert.Nil(t, err, "expected no error when pushing vault")

	tmpPath := filepath.Join(pwd, "tmp")
	err = os.MkdirAll(tmpPath, os.ModePerm)
	assert.Nil(t, err, "expected no error when creating tmp folder")
	defer os.RemoveAll(tmpPath)

	for _, file := range append(v.directories, v.files...) {
		if strings.Contains(filepath.ToSlash(file), "/") {
			continue
		}
		err := os.Rename(filepath.Join(v.localPath, file), filepath.Join(tmpPath, file))
		assert.Nil(t, err, "expected no error when renaming file")
	}

	err = v.Pull(password)
	assert.Nil(t, err, "expected no error when pushing vault")

	for _, file := range v.files {
		plaintext, err := os.ReadFile(filepath.Join(tmpPath, file))
		assert.Nil(t, err, "expected no error when reading plaintext file")

		decrypted, err := os.ReadFile(filepath.Join(v.localPath, file))
		assert.Nil(t, err, "expected no error when reading decrypted file")

		assert.Equal(t, string(plaintext), string(decrypted), "expected plaintext and decrypted files to be equal")
	}
}
