package crypto

import (
	"encoding/base64"
	"encoding/json"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewAESCrypto(t *testing.T) {
	tests := []struct {
		name    string
		key     string
		wantErr bool
		errMsg  string
	}{
		{
			name:    "valid 32-byte key",
			key:     "12345678901234567890123456789012",
			wantErr: false,
		},
		{
			name:    "key too short (16 bytes)",
			key:     "1234567890123456",
			wantErr: true,
			errMsg:  "encryption key must be exactly 32 bytes",
		},
		{
			name:    "key too long (64 bytes)",
			key:     "1234567890123456789012345678901212345678901234567890123456789012",
			wantErr: true,
			errMsg:  "encryption key must be exactly 32 bytes",
		},
		{
			name:    "empty key",
			key:     "",
			wantErr: true,
			errMsg:  "encryption key must be exactly 32 bytes",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			crypto, err := NewAESCrypto(tt.key)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
				assert.Nil(t, crypto)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, crypto)
				assert.Equal(t, []byte(tt.key), crypto.key)
			}
		})
	}
}

func TestAESCrypto_Encrypt_Decrypt(t *testing.T) {
	key := "12345678901234567890123456789012"
	crypto, err := NewAESCrypto(key)
	require.NoError(t, err)

	tests := []struct {
		name      string
		plaintext []byte
	}{
		{
			name:      "simple text",
			plaintext: []byte("Hello, World!"),
		},
		{
			name:      "empty string",
			plaintext: []byte(""),
		},
		{
			name:      "unicode text",
			plaintext: []byte("你好世界 🌍"),
		},
		{
			name:      "json data",
			plaintext: []byte(`{"user":"admin","password":"secret123"}`),
		},
		{
			name:      "long text",
			plaintext: []byte(strings.Repeat("A", 1000)),
		},
		{
			name:      "binary data",
			plaintext: []byte{0x00, 0x01, 0x02, 0xFF, 0xFE, 0xFD},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Encrypt
			ciphertext, err := crypto.Encrypt(tt.plaintext)
			require.NoError(t, err)
			assert.NotNil(t, ciphertext)

			// Ciphertext should be different from plaintext
			if len(tt.plaintext) > 0 {
				assert.NotEqual(t, tt.plaintext, ciphertext)
			}

			// Decrypt
			decrypted, err := crypto.Decrypt(ciphertext)
			require.NoError(t, err)
			// Handle empty byte slice comparison (nil vs empty slice)
			if len(tt.plaintext) == 0 {
				assert.Empty(t, decrypted)
			} else {
				assert.Equal(t, tt.plaintext, decrypted)
			}
		})
	}
}

func TestAESCrypto_Encrypt_Randomness(t *testing.T) {
	key := "12345678901234567890123456789012"
	crypto, err := NewAESCrypto(key)
	require.NoError(t, err)

	plaintext := []byte("test message")

	// Encrypt the same plaintext multiple times
	ciphertext1, err := crypto.Encrypt(plaintext)
	require.NoError(t, err)

	ciphertext2, err := crypto.Encrypt(plaintext)
	require.NoError(t, err)

	// Each encryption should produce different ciphertext (due to random nonce)
	assert.NotEqual(t, ciphertext1, ciphertext2,
		"encrypting the same plaintext twice should produce different ciphertexts")

	// Both should decrypt to the same plaintext
	decrypted1, err := crypto.Decrypt(ciphertext1)
	require.NoError(t, err)
	assert.Equal(t, plaintext, decrypted1)

	decrypted2, err := crypto.Decrypt(ciphertext2)
	require.NoError(t, err)
	assert.Equal(t, plaintext, decrypted2)
}

func TestAESCrypto_Decrypt_Errors(t *testing.T) {
	key := "12345678901234567890123456789012"
	crypto, err := NewAESCrypto(key)
	require.NoError(t, err)

	tests := []struct {
		name       string
		ciphertext []byte
		wantErr    bool
		errMsg     string
	}{
		{
			name:       "empty ciphertext",
			ciphertext: []byte{},
			wantErr:    true,
			errMsg:     "ciphertext too short",
		},
		{
			name:       "too short ciphertext",
			ciphertext: []byte{0x01, 0x02},
			wantErr:    true,
			errMsg:     "ciphertext too short",
		},
		{
			name:       "corrupted ciphertext",
			ciphertext: []byte("this is not a valid ciphertext at all"),
			wantErr:    true,
			errMsg:     "failed to decrypt",
		},
		{
			name:       "tampered ciphertext",
			ciphertext: func() []byte {
				ct, _ := crypto.Encrypt([]byte("original message"))
				// Tamper with the ciphertext
				ct[len(ct)-1] ^= 0xFF
				return ct
			}(),
			wantErr: true,
			errMsg:  "failed to decrypt",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			plaintext, err := crypto.Decrypt(tt.ciphertext)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
				assert.Nil(t, plaintext)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestAESCrypto_EncryptString_DecryptString(t *testing.T) {
	key := "12345678901234567890123456789012"
	crypto, err := NewAESCrypto(key)
	require.NoError(t, err)

	tests := []struct {
		name      string
		plaintext string
	}{
		{
			name:      "simple string",
			plaintext: "Hello, World!",
		},
		{
			name:      "empty string",
			plaintext: "",
		},
		{
			name:      "unicode string",
			plaintext: "こんにちは世界",
		},
		{
			name:      "special characters",
			plaintext: "!@#$%^&*()_+-=[]{}|;:',.<>?/~`",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Encrypt
			ciphertext, err := crypto.EncryptString(tt.plaintext)
			require.NoError(t, err)
			assert.NotEmpty(t, ciphertext)

			// Verify base64 encoding
			_, err = base64.StdEncoding.DecodeString(ciphertext)
			assert.NoError(t, err, "ciphertext should be valid base64")

			// Decrypt
			decrypted, err := crypto.DecryptString(ciphertext)
			require.NoError(t, err)
			assert.Equal(t, tt.plaintext, decrypted)
		})
	}
}

func TestAESCrypto_DecryptString_InvalidBase64(t *testing.T) {
	key := "12345678901234567890123456789012"
	crypto, err := NewAESCrypto(key)
	require.NoError(t, err)

	// Invalid base64 string
	_, err = crypto.DecryptString("not-valid-base64!!!")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to decode base64")
}

func TestAESCrypto_EncryptJSON_DecryptJSON(t *testing.T) {
	key := "12345678901234567890123456789012"
	crypto, err := NewAESCrypto(key)
	require.NoError(t, err)

	type TestStruct struct {
		Name  string `json:"name"`
		Age   int    `json:"age"`
		Email string `json:"email"`
	}

	tests := []struct {
		name string
		data interface{}
	}{
		{
			name: "simple struct",
			data: TestStruct{
				Name:  "John Doe",
				Age:   30,
				Email: "john@example.com",
			},
		},
		{
			name: "map",
			data: map[string]interface{}{
				"username": "admin",
				"password": "secret123",
				"roles":    []string{"admin", "user"},
			},
		},
		{
			name: "array",
			data: []string{"apple", "banana", "cherry"},
		},
		{
			name: "nested struct",
			data: map[string]interface{}{
				"user": map[string]interface{}{
					"name": "Alice",
					"credentials": map[string]string{
						"username": "alice",
						"password": "pass123",
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Encrypt
			ciphertext, err := crypto.EncryptJSON(tt.data)
			require.NoError(t, err)
			assert.NotNil(t, ciphertext)

			// Decrypt
			switch tt.data.(type) {
			case TestStruct:
				var result TestStruct
				err = crypto.DecryptJSON(ciphertext, &result)
				require.NoError(t, err)
				assert.Equal(t, tt.data, result)
			case map[string]interface{}:
				var result map[string]interface{}
				err = crypto.DecryptJSON(ciphertext, &result)
				require.NoError(t, err)
				// Compare JSON representation for deep equality
				// JSON unmarshaling loses type info for nested structures
				expectedJSON, _ := json.Marshal(tt.data)
				resultJSON, _ := json.Marshal(result)
				assert.JSONEq(t, string(expectedJSON), string(resultJSON))
			case []string:
				var result []string
				err = crypto.DecryptJSON(ciphertext, &result)
				require.NoError(t, err)
				assert.Equal(t, tt.data, result)
			}
		})
	}
}

func TestAESCrypto_EncryptJSON_InvalidData(t *testing.T) {
	key := "12345678901234567890123456789012"
	crypto, err := NewAESCrypto(key)
	require.NoError(t, err)

	// Invalid data that cannot be marshaled to JSON (function)
	invalidData := func() {}

	_, err = crypto.EncryptJSON(invalidData)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to marshal JSON")
}

func TestAESCrypto_DecryptJSON_InvalidJSON(t *testing.T) {
	key := "12345678901234567890123456789012"
	crypto, err := NewAESCrypto(key)
	require.NoError(t, err)

	// Encrypt invalid JSON data
	ciphertext, err := crypto.Encrypt([]byte("not valid json"))
	require.NoError(t, err)

	var result map[string]interface{}
	err = crypto.DecryptJSON(ciphertext, &result)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to unmarshal JSON")
}

func TestAESCrypto_DifferentKeys(t *testing.T) {
	key1 := "12345678901234567890123456789012"
	key2 := "abcdefghijklmnopqrstuvwxyz123456"

	crypto1, err := NewAESCrypto(key1)
	require.NoError(t, err)

	crypto2, err := NewAESCrypto(key2)
	require.NoError(t, err)

	plaintext := []byte("secret message")

	// Encrypt with crypto1
	ciphertext, err := crypto1.Encrypt(plaintext)
	require.NoError(t, err)

	// Try to decrypt with crypto2 (different key) - should fail
	_, err = crypto2.Decrypt(ciphertext)
	assert.Error(t, err, "decrypting with wrong key should fail")
	assert.Contains(t, err.Error(), "failed to decrypt")
}

func BenchmarkAESCrypto_Encrypt(b *testing.B) {
	key := "12345678901234567890123456789012"
	crypto, _ := NewAESCrypto(key)
	plaintext := []byte("benchmark test data")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = crypto.Encrypt(plaintext)
	}
}

func BenchmarkAESCrypto_Decrypt(b *testing.B) {
	key := "12345678901234567890123456789012"
	crypto, _ := NewAESCrypto(key)
	plaintext := []byte("benchmark test data")
	ciphertext, _ := crypto.Encrypt(plaintext)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = crypto.Decrypt(ciphertext)
	}
}

func BenchmarkAESCrypto_EncryptString(b *testing.B) {
	key := "12345678901234567890123456789012"
	crypto, _ := NewAESCrypto(key)
	plaintext := "benchmark test string"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = crypto.EncryptString(plaintext)
	}
}

func BenchmarkAESCrypto_DecryptString(b *testing.B) {
	key := "12345678901234567890123456789012"
	crypto, _ := NewAESCrypto(key)
	plaintext := "benchmark test string"
	ciphertext, _ := crypto.EncryptString(plaintext)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = crypto.DecryptString(ciphertext)
	}
}
