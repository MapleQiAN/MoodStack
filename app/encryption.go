package app

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io"

	"golang.org/x/crypto/pbkdf2"
)

const (
	keySize    = 32 // AES-256
	saltSize   = 16
	ivSize     = 16
	iterations = 100000
)

// EncryptionKey represents an encryption key with salt
type EncryptionKey struct {
	Key  []byte
	Salt []byte
}

// GenerateKey generates a new encryption key from password
func GenerateKey(password string) (*EncryptionKey, error) {
	salt := make([]byte, saltSize)
	if _, err := rand.Read(salt); err != nil {
		return nil, fmt.Errorf("failed to generate salt: %v", err)
	}

	key := pbkdf2.Key([]byte(password), salt, iterations, keySize, sha256.New)

	return &EncryptionKey{
		Key:  key,
		Salt: salt,
	}, nil
}

// DeriveKey derives encryption key from password and salt
func DeriveKey(password string, salt []byte) []byte {
	return pbkdf2.Key([]byte(password), salt, iterations, keySize, sha256.New)
}

// EncryptData encrypts data using AES-256-GCM
func EncryptData(data []byte, key []byte) ([]byte, []byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create cipher: %v", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create GCM: %v", err)
	}

	// Generate random IV
	iv := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, nil, fmt.Errorf("failed to generate IV: %v", err)
	}

	// Encrypt data
	ciphertext := gcm.Seal(nil, iv, data, nil)

	return ciphertext, iv, nil
}

// DecryptData decrypts data using AES-256-GCM
func DecryptData(ciphertext []byte, key []byte, iv []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("failed to create cipher: %v", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("failed to create GCM: %v", err)
	}

	// Decrypt data
	plaintext, err := gcm.Open(nil, iv, ciphertext, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt data: %v", err)
	}

	return plaintext, nil
}

// EncryptString encrypts a string and returns base64 encoded result
func EncryptString(text string, key []byte) (string, string, error) {
	ciphertext, iv, err := EncryptData([]byte(text), key)
	if err != nil {
		return "", "", err
	}

	return base64.StdEncoding.EncodeToString(ciphertext),
		base64.StdEncoding.EncodeToString(iv), nil
}

// DecryptString decrypts a base64 encoded string
func DecryptString(ciphertext string, ivStr string, key []byte) (string, error) {
	ciphertextBytes, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return "", fmt.Errorf("failed to decode ciphertext: %v", err)
	}

	iv, err := base64.StdEncoding.DecodeString(ivStr)
	if err != nil {
		return "", fmt.Errorf("failed to decode IV: %v", err)
	}

	plaintext, err := DecryptData(ciphertextBytes, key, iv)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}

// HashPassword creates a SHA-256 hash of the password with salt
func HashPassword(password string, salt []byte) string {
	hash := sha256.Sum256(append([]byte(password), salt...))
	return base64.StdEncoding.EncodeToString(hash[:])
}

// VerifyPassword verifies a password against its hash
func VerifyPassword(password string, hash string, salt []byte) bool {
	expectedHash := HashPassword(password, salt)
	return expectedHash == hash
}
