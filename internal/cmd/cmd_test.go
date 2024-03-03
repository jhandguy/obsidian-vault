package cmd

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRunSuccessfulCommand(t *testing.T) {
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}

	err := Run("echo", "foo", stdout, stderr)

	assert.NoError(t, err)
	assert.Equal(t, "-c foo\n", stdout.String())
	assert.Empty(t, stderr.String())
}

func TestRunNonexistentCommand(t *testing.T) {
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}

	err := Run("nonexistent", "", stdout, stderr)

	assert.Error(t, err)
	assert.Empty(t, stdout.String())
	assert.Empty(t, stderr.String())
}
