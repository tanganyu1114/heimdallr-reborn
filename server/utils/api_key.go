package utils

import (
	"crypto/rand"
	"encoding/hex"

	"golang.org/x/crypto/bcrypt"
)

// GenerateAPIKey generates a random API key with "ak_" prefix
func GenerateAPIKey() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return "ak_" + hex.EncodeToString(bytes), nil
}

// GenerateAPISecret generates a random API secret with "sk_" prefix
func GenerateAPISecret() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return "sk_" + hex.EncodeToString(bytes), nil
}

// BcryptHash encrypts a string using bcrypt
func BcryptHash(str string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(str), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

// BcryptCheck verifies a string against a bcrypt hash
func BcryptCheck(plain, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(plain))
	return err == nil
}
