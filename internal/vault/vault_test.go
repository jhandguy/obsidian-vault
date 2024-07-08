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
	assert.NoError(t, err)

	path := filepath.Join(pwd, "../../example")
	config := ".obsidian"
	password := "consectetur-adipiscing-elit"

	err = os.Setenv("SHELL", "echo")
	assert.NoError(t, err)

	v, err := New(path, config)
	assert.NoError(t, err)

	gitPath, err := v.getVaultPath(vaultTypeGit)
	assert.NoError(t, err)

	err = os.MkdirAll(gitPath, os.ModePerm)
	assert.NoError(t, err)
	defer os.RemoveAll(gitPath)

	err = v.Push(password)
	assert.NoError(t, err)

	tmpPath := filepath.Join(pwd, "tmp")
	err = os.MkdirAll(tmpPath, os.ModePerm)
	assert.NoError(t, err)
	defer os.RemoveAll(tmpPath)

	for _, file := range append(v.directories, v.files...) {
		if strings.Contains(filepath.ToSlash(file), "/") {
			continue
		}
		err := os.Rename(filepath.Join(v.localPath, file), filepath.Join(tmpPath, file))
		assert.NoError(t, err)
	}

	err = v.Pull(password)
	assert.NoError(t, err)

	for _, file := range v.files {
		plaintext, err := os.ReadFile(filepath.Join(tmpPath, file))
		assert.NoError(t, err)

		decrypted, err := os.ReadFile(filepath.Join(v.localPath, file))
		assert.NoError(t, err)
		assert.Equal(t, string(plaintext), string(decrypted))
	}

	v.localPath = tmpPath
	err = v.Clean(true)
	assert.NoError(t, err)

	for _, file := range v.files {
		_, err := os.Stat(filepath.Join(tmpPath, file))
		assert.True(t, os.IsNotExist(err))

		_, err = os.Stat(filepath.Join(v.gitPath, file))
		assert.True(t, os.IsNotExist(err))
	}
}
