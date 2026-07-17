// Package utils provides RSA encryption and challenge-response utilities for secure authentication
package utils

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/hex"
	"encoding/pem"
	"sync"
	"time"

	"github.com/marmotedu/errors"
)

var (
	privateKey     *rsa.PrivateKey
	publicKey      *rsa.PublicKey
	keyOnce        sync.Once
	challenges     = make(map[string]challenge) // sessionID -> challenge data
	challengesLock sync.RWMutex
)

// challenge stores the challenge string and its creation time for replay attack prevention
type challenge struct {
	Challenge string
	Timestamp time.Time
}

// IsExpired checks if the challenge has expired (5 minutes lifetime)
func (c challenge) IsExpired() bool {
	return time.Since(c.Timestamp) > 5*time.Minute
}

// GenerateRSAKeys generates RSA key pair if not exists
func GenerateRSAKeys() {
	keyOnce.Do(func() {
		var err error
		privateKey, err = rsa.GenerateKey(rand.Reader, 2048)
		if err != nil {
			panic("failed to generate RSA key: " + err.Error())
		}
		publicKey = &privateKey.PublicKey
	})
}

// GenerateChallenge creates a new random challenge string (16 bytes hex-encoded)
func GenerateChallenge() (string, error) {
	bytes := make([]byte, 16)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", errors.Wrap(err, "failed to generate challenge")
	}
	return hex.EncodeToString(bytes), nil
}

// GetPublicKeyWithChallenge returns public key and a new challenge for a session.
// If the session already has a valid (non-expired) challenge, it will be reused.
func GetPublicKeyWithChallenge(sessionID string) (string, string, error) {
	if publicKey == nil {
		GenerateRSAKeys()
	}

	challengesLock.Lock()
	c, ok := challenges[sessionID]
	if !ok || c.IsExpired() {
		// Generate new challenge for new or expired sessions
		newChallengeStr, err := GenerateChallenge()
		if err != nil {
			challengesLock.Unlock()
			return "", "", errors.Wrap(err, "failed to generate challenge")
		}
		c = challenge{
			Challenge: newChallengeStr,
			Timestamp: time.Now(),
		}
		challenges[sessionID] = c
	}
	challengesLock.Unlock()

	pubBytes, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		return "", "", errors.Wrap(err, "failed to marshal public key")
	}
	block := &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: pubBytes,
	}
	return string(pem.EncodeToMemory(block)), c.Challenge, nil
}

// VerifyChallenge checks if the provided challenge matches for a session.
// The challenge is consumed (deleted) after verification, whether successful or not.
func VerifyChallenge(sessionID, challengeStr string) bool {
	challengesLock.RLock()
	c, exists := challenges[sessionID]
	challengesLock.RUnlock()

	if !exists {
		return false
	}

	// Check if challenge is expired
	if c.IsExpired() {
		// Clean up expired challenge
		challengesLock.Lock()
		delete(challenges, sessionID)
		challengesLock.Unlock()
		return false
	}

	// Challenge is valid, clean it up after use (one-time use)
	challengesLock.Lock()
	delete(challenges, sessionID)
	challengesLock.Unlock()

	return c.Challenge == challengeStr
}

// CleanExpiredChallenges cleans up expired challenges from the map.
// Should be called periodically by a background goroutine.
func CleanExpiredChallenges() {
	challengesLock.Lock()
	defer challengesLock.Unlock()

	now := time.Now()
	for sessionID, c := range challenges {
		if now.Sub(c.Timestamp) > 5*time.Minute {
			delete(challenges, sessionID)
		}
	}
}

// GetPublicKeyPEM returns the public key in PEM format (legacy method without challenge)
func GetPublicKeyPEM() (string, error) {
	if publicKey == nil {
		GenerateRSAKeys()
	}
	pubBytes, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		return "", errors.Wrap(err, "failed to marshal public key")
	}
	block := &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: pubBytes,
	}
	return string(pem.EncodeToMemory(block)), nil
}

// RSADecrypt decrypts base64-encoded RSA encrypted data using PKCS1v15
func RSADecrypt(ciphertextBase64 string) (string, error) {
	if privateKey == nil {
		GenerateRSAKeys()
	}

	ciphertext, err := base64.StdEncoding.DecodeString(ciphertextBase64)
	if err != nil {
		return "", errors.New("invalid base64 encoding, check input format")
	}

	plaintext, err := rsa.DecryptPKCS1v15(rand.Reader, privateKey, ciphertext)
	if err != nil {
		return "", errors.New("RSA decryption failed, check encrypted data or regenerate")
	}

	return string(plaintext), nil
}

// RSAEncrypt encrypts data using a PEM-encoded public key with PKCS1v15
func RSAEncrypt(publicKeyPEM, plaintext string) (string, error) {
	block, _ := pem.Decode([]byte(publicKeyPEM))
	if block == nil {
		return "", errors.New("failed to parse PEM block containing the public key")
	}

	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return "", errors.Wrap(err, "failed to parse DER encoded public key")
	}

	rsaPub, ok := pub.(*rsa.PublicKey)
	if !ok {
		return "", errors.New("not an RSA public key")
	}

	ciphertext, err := rsa.EncryptPKCS1v15(rand.Reader, rsaPub, []byte(plaintext))
	if err != nil {
		return "", errors.Wrap(err, "failed to encrypt data")
	}

	return base64.StdEncoding.EncodeToString(ciphertext), nil
}
