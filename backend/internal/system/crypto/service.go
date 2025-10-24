/*
 * Copyright (c) 2025, WSO2 LLC. (http://www.wso2.com).
 *
 * WSO2 LLC. licenses this file to you under the Apache License,
 * Version 2.0 (the "License"); you may not use this file except
 * in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

// Package crypto provides cryptographic functionality with algorithm agility.
package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"sync"

	"github.com/asgardeo/thunder/internal/system/config"
	"github.com/asgardeo/thunder/internal/system/hash"
	"github.com/asgardeo/thunder/internal/system/log"
)

const (
	// aesgcmAlgorithm represents AES-GCM algorithm
	aesgcmAlgorithm = "AES-GCM"
	// defaultKeySize defines the default key size for AES-GCM
	defaultKeySize = 32
)

// CryptoService provides cryptographic operations.
type CryptoService struct {
	Key []byte
	Kid string
}

var (
	// instance is the singleton instance of CryptoService
	instance *CryptoService
	// once ensures the singleton is initialized only once
	once sync.Once
)

// GetCryptoService creates and returns a singleton instance of the CryptoService.
func GetCryptoService() *CryptoService {
	logger := log.GetLogger().With(log.String(log.LoggerKeyComponentName, "CryptoService"))
	once.Do(func() {
		var err error
		instance, err = initCryptoService()
		if err != nil {
			logger.Error("Failed to initialize CryptoService: %v", log.Error(err))
		}
	})
	return instance
}

// initCryptoService initializes the CryptoService from configuration sources.
func initCryptoService() (*CryptoService, error) {
	logger := log.GetLogger().With(log.String(log.LoggerKeyComponentName, "CryptoService"))
	// Try to get key from the application configuration
	config := config.GetThunderRuntime().Config.Crypto.Key // Use the correct config getter

	// Check if crypto configuration exists
	if config != "" {
		key, err := hex.DecodeString(config)
		if err == nil {
			logger.Debug("Using crypto key from configuration")
			return NewCryptoService(key)
		}
		logger.Warn("Invalid crypto key in configuration, generating a new key")
	}

	// Generate new key as fallback for development
	logger.Warn("No valid crypto key found in configuration, generating a new one")

	key, err := generateRandomKey(defaultKeySize)
	if err != nil {
		return nil, err
	}

	// Print the generated key for development purposes
	encodedKey := hex.EncodeToString(key)
	logger.Debug("Generated new crypto key (hex): %s", log.String("logKey", encodedKey))

	return NewCryptoService(key)
}

// NewCryptoService creates a new instance of CryptoService with the provided key.
func NewCryptoService(key []byte) (*CryptoService, error) {
	// Check key size for algorithm

	return &CryptoService{
		Key: key,
		Kid: getKeyID(key), // Generate a unique key ID
	}, nil
}

// Encrypt encrypts the given plaintext and returns a JSON string
// containing the encrypted data.
func (cs *CryptoService) Encrypt(plaintext []byte) (string, error) {
	// Create AES cipher
	block, err := aes.NewCipher(cs.Key)
	if err != nil {
		return "", err
	}

	// Create GCM mode
	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	// Create a nonce
	nonce := make([]byte, aesgcm.NonceSize())
	if _, err = rand.Read(nonce); err != nil {
		return "", err
	}

	// Encrypt and authenticate plaintext, prepend nonce
	ciphertext := aesgcm.Seal(nonce, nonce, plaintext, nil)

	// Create metadata structure
	encData := EncryptedData{
		Algorithm:  aesgcmAlgorithm,
		Ciphertext: base64.StdEncoding.EncodeToString(ciphertext),
		KeyID:      cs.Kid, // Unique identifier for the key
	}

	// Serialize to JSON
	jsonData, err := json.Marshal(encData)
	if err != nil {
		return "", err
	}

	return string(jsonData), nil
}

// Decrypt decrypts the given JSON string produced by Encrypt
// and returns the original plaintext.
func (cs *CryptoService) Decrypt(encodedData string) ([]byte, error) {
	// Deserialize JSON
	var encData EncryptedData
	if err := json.Unmarshal([]byte(encodedData), &encData); err != nil {
		return nil, fmt.Errorf("invalid data format: %w", err)
	}

	// Verify algorithm
	if encData.Algorithm != aesgcmAlgorithm {
		return nil, fmt.Errorf("unsupported algorithm: %s", encData.Algorithm)
	}

	// Decode the payload
	ciphertext, err := base64.StdEncoding.DecodeString(encData.Ciphertext)
	if err != nil {
		return nil, fmt.Errorf("invalid payload encoding: %w", err)
	}

	// Create AES cipher
	block, err := aes.NewCipher(cs.Key)
	if err != nil {
		return nil, err
	}

	// Create GCM mode
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	// Verify ciphertext length
	nonceSize := aesGCM.NonceSize()
	if len(ciphertext) < nonceSize {
		return nil, errors.New("ciphertext too short")
	}

	// Extract nonce and decrypt
	nonce, encryptedData := ciphertext[:nonceSize], ciphertext[nonceSize:]
	plaintext, err := aesGCM.Open(nil, nonce, encryptedData, nil)
	if err != nil {
		return nil, err
	}

	return plaintext, nil
}

// EncryptString encrypts the given plaintext string and returns a
// JSON string containing the encrypted data.
func (cs *CryptoService) EncryptString(plaintext string) (string, error) {
	return cs.Encrypt([]byte(plaintext))
}

// DecryptString decrypts the given JSON string produced by Encrypt
// and returns the original plaintext string.
func (cs *CryptoService) DecryptString(ciphertext string) (string, error) {
	plaintext, err := cs.Decrypt(ciphertext)
	if err != nil {
		return "", err
	}
	return string(plaintext), nil
}

// generateRandomKey generates a random key of the specified size.
func generateRandomKey(keySize int) ([]byte, error) {
	key := make([]byte, keySize)
	_, err := rand.Read(key)
	if err != nil {
		return nil, err
	}
	return key, nil
}

// getKeyID generates a unique identifier for the key.
func getKeyID(key []byte) string {
	return hash.GenerateThumbprint(key)
}
