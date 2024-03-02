package git

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

var git = New("echo", "/tmp")

func TestAdd(t *testing.T) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	err := git.Add(&stdout, &stderr)
	assert.Nil(t, err, "expected no error when adding git changes")
	assert.Equal(t, "-c git --git-dir /tmp/.git --work-tree /tmp add .\n", stdout.String())
	assert.Empty(t, stderr.String(), "expected no error when adding git changes")
}

func TestCommit(t *testing.T) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	err := git.Commit(&stdout, &stderr, "commit message")
	assert.Nil(t, err, "expected no error when committing git changes")
	assert.Equal(t, `-c git --git-dir /tmp/.git --work-tree /tmp commit -m "commit message"`+"\n", stdout.String())
	assert.Empty(t, stderr.String(), "expected no error when committing git changes")
}

func TestPush(t *testing.T) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	err := git.Push(&stdout, &stderr)
	assert.Nil(t, err, "expected no error when pushing git changes")
	assert.Equal(t, `-c git --git-dir /tmp/.git --work-tree /tmp push origin main`+"\n", stdout.String())
	assert.Empty(t, stderr.String(), "expected no error when pushing git changes")
}

func TestPull(t *testing.T) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	err := git.Pull(&stdout, &stderr)
	assert.Nil(t, err, "expected no error when pulling git changes")
	assert.Equal(t, `-c git --git-dir /tmp/.git --work-tree /tmp pull origin main`+"\n", stdout.String())
	assert.Empty(t, stderr.String(), "expected no error when pulling git changes")
}
