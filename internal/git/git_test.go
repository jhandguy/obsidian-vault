package git

import (
	"bytes"
	"fmt"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

var tmpPath = filepath.Clean("/tmp")
var git = New("echo", tmpPath)
var gitCommand = fmt.Sprintf("git --git-dir %s --work-tree %s", filepath.Join(tmpPath, HiddenFolder), tmpPath)

func TestAdd(t *testing.T) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	err := git.Add(&stdout, &stderr)
	assert.NoError(t, err)
	assert.Equal(t, fmt.Sprintf("-c %s add .\n", gitCommand), stdout.String())
	assert.Empty(t, stderr.String())
}

func TestCommit(t *testing.T) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	msg := "commit message"
	err := git.Commit(&stdout, &stderr, msg)
	assert.NoError(t, err)
	assert.Equal(t, fmt.Sprintf("-c %s commit -m \"%s\"\n", gitCommand, msg), stdout.String())
	assert.Empty(t, stderr.String())
}

func TestPush(t *testing.T) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	err := git.Push(&stdout, &stderr)
	assert.NoError(t, err)
	assert.Equal(t, fmt.Sprintf("-c %s push origin main\n", gitCommand), stdout.String())
	assert.Empty(t, stderr.String())
}

func TestPull(t *testing.T) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	err := git.Pull(&stdout, &stderr)
	assert.NoError(t, err)
	assert.Equal(t, fmt.Sprintf("-c %s pull origin main\n", gitCommand), stdout.String())
	assert.Empty(t, stderr.String())
}
