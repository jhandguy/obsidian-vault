package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"io"

	"golang.org/x/crypto/scrypt"
)

type Crypto struct{}

func New() *Crypto {
	return &Crypto{}
}

func (c *Crypto) Decrypt(data []byte, password string, fileName string) ([]byte, error) {
	key, err := c.deriveKey(password, fileName)
	if err != nil {
		return nil, err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonceSize := gcm.NonceSize()
	nonce, ciphertext := data[:nonceSize], data[nonceSize:]

	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, err
	}

	return plaintext, nil
}

func (c *Crypto) Encrypt(plaintext []byte, password, fileName string) ([]byte, error) {
	key, err := c.deriveKey(password, fileName)
	if err != nil {
		return nil, err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	return gcm.Seal(nonce, nonce, plaintext, nil), nil
}

func (c *Crypto) deriveKey(password, fileName string) ([]byte, error) {
	return scrypt.Key([]byte(password), []byte(fileName), 32768, 8, 1, 32)
}
