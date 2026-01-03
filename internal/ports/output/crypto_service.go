package output

// CryptoService defines the output port for encryption/decryption operations.
// This interface abstracts cryptographic operations for securing sensitive data.
type CryptoService interface {
	// Encrypt encrypts the given plaintext using AES-256-GCM.
	// Returns the encrypted ciphertext (including nonce) or an error.
	Encrypt(plaintext []byte) ([]byte, error)

	// Decrypt decrypts the given ciphertext using AES-256-GCM.
	// Returns the decrypted plaintext or an error if decryption fails.
	Decrypt(ciphertext []byte) ([]byte, error)

	// EncryptString is a convenience method for encrypting strings.
	// It converts the string to bytes, encrypts it, and returns base64-encoded ciphertext.
	EncryptString(plaintext string) (string, error)

	// DecryptString is a convenience method for decrypting strings.
	// It decodes the base64 ciphertext, decrypts it, and returns the plaintext string.
	DecryptString(ciphertext string) (string, error)

	// EncryptJSON encrypts a JSON-serializable object.
	// It marshals the object to JSON, encrypts it, and returns the ciphertext.
	EncryptJSON(data interface{}) ([]byte, error)

	// DecryptJSON decrypts ciphertext and unmarshals it into the target object.
	// The target must be a pointer to the destination struct.
	DecryptJSON(ciphertext []byte, target interface{}) error
}
