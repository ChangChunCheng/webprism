package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
)

// AESCrypto implements the CryptoService interface using AES-256-GCM.
type AESCrypto struct {
	key []byte
}

// NewAESCrypto creates a new AESCrypto instance.
// The key must be exactly 32 bytes for AES-256.
func NewAESCrypto(key string) (*AESCrypto, error) {
	if len(key) != 32 {
		return nil, fmt.Errorf("encryption key must be exactly 32 bytes, got %d", len(key))
	}

	return &AESCrypto{
		key: []byte(key),
	}, nil
}

// Encrypt encrypts the given plaintext using AES-256-GCM.
// Returns the encrypted ciphertext (nonce + ciphertext) or an error.
func (c *AESCrypto) Encrypt(plaintext []byte) ([]byte, error) {
	// Create AES cipher block
	block, err := aes.NewCipher(c.key)
	if err != nil {
		return nil, fmt.Errorf("failed to create cipher block: %w", err)
	}

	// Create GCM mode
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("failed to create GCM: %w", err)
	}

	// Generate random nonce
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, fmt.Errorf("failed to generate nonce: %w", err)
	}

	// Encrypt and append nonce to the beginning
	ciphertext := gcm.Seal(nonce, nonce, plaintext, nil)

	return ciphertext, nil
}

// Decrypt decrypts the given ciphertext using AES-256-GCM.
// Returns the decrypted plaintext or an error if decryption fails.
func (c *AESCrypto) Decrypt(ciphertext []byte) ([]byte, error) {
	// Create AES cipher block
	block, err := aes.NewCipher(c.key)
	if err != nil {
		return nil, fmt.Errorf("failed to create cipher block: %w", err)
	}

	// Create GCM mode
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("failed to create GCM: %w", err)
	}

	// Check minimum ciphertext length
	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return nil, fmt.Errorf("ciphertext too short")
	}

	// Extract nonce and ciphertext
	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]

	// Decrypt
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt: %w", err)
	}

	return plaintext, nil
}

// EncryptString is a convenience method for encrypting strings.
// It converts the string to bytes, encrypts it, and returns base64-encoded ciphertext.
func (c *AESCrypto) EncryptString(plaintext string) (string, error) {
	ciphertext, err := c.Encrypt([]byte(plaintext))
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// DecryptString is a convenience method for decrypting strings.
// It decodes the base64 ciphertext, decrypts it, and returns the plaintext string.
func (c *AESCrypto) DecryptString(ciphertext string) (string, error) {
	decoded, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return "", fmt.Errorf("failed to decode base64: %w", err)
	}

	plaintext, err := c.Decrypt(decoded)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}

// EncryptJSON encrypts a JSON-serializable object.
// It marshals the object to JSON, encrypts it, and returns the ciphertext.
func (c *AESCrypto) EncryptJSON(data interface{}) ([]byte, error) {
	// Marshal to JSON
	plaintext, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal JSON: %w", err)
	}

	// Encrypt
	return c.Encrypt(plaintext)
}

// DecryptJSON decrypts ciphertext and unmarshals it into the target object.
// The target must be a pointer to the destination struct.
func (c *AESCrypto) DecryptJSON(ciphertext []byte, target interface{}) error {
	// Decrypt
	plaintext, err := c.Decrypt(ciphertext)
	if err != nil {
		return err
	}

	// Unmarshal JSON
	if err := json.Unmarshal(plaintext, target); err != nil {
		return fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	return nil
}
