package crypto

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEncryptionAndDecryptionPreserveOriginalData(t *testing.T) {
	c := New()
	plaintext := []byte("Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua.")
	password := "consectetur-adipiscing-elit"
	fileName := "LoremIpsum.md"

	data, err := c.Encrypt(plaintext, password, fileName)

	assert.NoError(t, err)
	assert.NotEmpty(t, data)
	assert.NotEqual(t, plaintext, data)

	decrypted, err := c.Decrypt(data, password, fileName)

	assert.NoError(t, err)
	assert.NotEmpty(t, decrypted)
	assert.Equal(t, plaintext, decrypted)

	encrypted, err := c.Encrypt(plaintext, password, fileName)

	assert.NoError(t, err)
	assert.NotEmpty(t, encrypted)
	assert.NotEqual(t, plaintext, encrypted)
	assert.NotEqual(t, data, encrypted)
}
