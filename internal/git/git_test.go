package git

import (
	"bytes"
	"fmt"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

var git = New("echo", filepath.Clean("/tmp"))

func TestAdd(t *testing.T) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	err := git.Add(&stdout, &stderr)
	assert.Nil(t, err, "expected no error when adding git changes")
	assert.Equal(t, fmt.Sprintf("-c git --git-dir %s --work-tree %s add .\n", filepath.Clean("/tmp/.git"), filepath.Clean("/tmp")), stdout.String())
	assert.Empty(t, stderr.String(), "expected no error when adding git changes")
}

func TestCommit(t *testing.T) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	err := git.Commit(&stdout, &stderr, "commit message")
	assert.Nil(t, err, "expected no error when committing git changes")
	assert.Equal(t, fmt.Sprintf("-c git --git-dir %s --work-tree %s commit -m \"commit message\"\n", filepath.Clean("/tmp/.git"), filepath.Clean("/tmp")), stdout.String())
	assert.Empty(t, stderr.String(), "expected no error when committing git changes")
}

func TestPush(t *testing.T) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	err := git.Push(&stdout, &stderr)
	assert.Nil(t, err, "expected no error when pushing git changes")
	assert.Equal(t, fmt.Sprintf("-c git --git-dir %s --work-tree %s push origin main\n", filepath.Clean("/tmp/.git"), filepath.Clean("/tmp")), stdout.String())
	assert.Empty(t, stderr.String(), "expected no error when pushing git changes")
}

func TestPull(t *testing.T) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	err := git.Pull(&stdout, &stderr)
	assert.Nil(t, err, "expected no error when pulling git changes")
	assert.Equal(t, fmt.Sprintf("-c git --git-dir %s --work-tree %s pull origin main\n", filepath.Clean("/tmp/.git"), filepath.Clean("/tmp")), stdout.String())
	assert.Empty(t, stderr.String(), "expected no error when pulling git changes")
}
