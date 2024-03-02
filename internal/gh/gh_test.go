package gh

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

var gh = New("echo", "/tmp", "test")

func TestCreateRepository(t *testing.T) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	err := gh.CreateRepository(&stdout, &stderr)
	assert.Nil(t, err, "expected no error when creating repository")
	assert.Equal(t, "-c gh repo create test --description \"Encrypted backup of test, created with obsidian-vault.\" --private --disable-issues --disable-wiki\n", stdout.String())
	assert.Empty(t, stderr.String(), "expected no error when creating repository")
}

func TestCloneRepository(t *testing.T) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	err := gh.CloneRepository(&stdout, &stderr)
	assert.Nil(t, err, "expected no error when cloning repository")
	assert.Equal(t, "-c gh repo clone test /tmp\n", stdout.String())
	assert.Empty(t, stderr.String(), "expected no error when cloning repository")
}
