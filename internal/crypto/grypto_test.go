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

	assert.Nil(t, err, "expected no error when encrypting data")
	assert.NotEmpty(t, data, "expected encrypted data to be non-empty")
	assert.NotEqualf(t, plaintext, data, "expected encrypted data to be different from plaintext")

	decrypted, err := c.Decrypt(data, password, fileName)

	assert.Nil(t, err, "expected no error when decrypting data")
	assert.NotEmpty(t, decrypted, "expected decrypted data to be non-empty")
	assert.Equalf(t, plaintext, decrypted, "expected decrypted data to be the same as plaintext")

	encrypted, err := c.Encrypt(plaintext, password, fileName)

	assert.Nil(t, err, "expected no error when encrypting data")
	assert.NotEmpty(t, encrypted, "expected encrypted data to be non-empty")
	assert.NotEqualf(t, plaintext, encrypted, "expected encrypted data to be different from plaintext")
	assert.NotEqualf(t, data, encrypted, "expected encrypted data to be different from previous encrypted data")
}
