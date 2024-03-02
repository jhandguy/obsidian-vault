package vault

import (
	"os"
	"path/filepath"
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

	// TODO: check git files were created and encrypted
	// TODO: remove local files

	err = v.Pull(password)
	assert.Nil(t, err, "expected no error when pushing vault")

	// TODO: check local files were recreated and decrypted
}
