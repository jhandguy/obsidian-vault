package gh

import (
	"bytes"
	"fmt"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

var tmpPath = filepath.Clean("/tmp")
var repoName = "test"
var gh = New("echo", tmpPath, repoName)

func TestCreateRepository(t *testing.T) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	err := gh.CreateRepository(&stdout, &stderr)
	assert.NoError(t, err)
	assert.Equal(t, fmt.Sprintf("-c gh repo create %s --description \"Encrypted backup of test, created with obsidian-vault.\" --private --disable-issues --disable-wiki\n", repoName), stdout.String())
	assert.Empty(t, stderr.String())
}

func TestCloneRepository(t *testing.T) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	err := gh.CloneRepository(&stdout, &stderr)
	assert.NoError(t, err)
	assert.Equal(t, fmt.Sprintf("-c gh repo clone %s %s\n", repoName, tmpPath), stdout.String())
	assert.Empty(t, stderr.String())
}
