package services

import (
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"

	"golang.org/x/crypto/argon2"
)

type PasswordService struct {
	// There should be argon2 configuration in here, but for simplicity let's not do that for today.
	RandomService *RandomService
}

type Argon2Configuration struct {
	HashRaw    []byte
	Salt       []byte
	TimeCost   uint32
	MemoryCost uint32
	Threads    uint8
	KeyLength  uint32
}

// HashPassword uses Argon2ID to hash a password, and returns a saved hash that combines all information
// for later verification.
func (s *PasswordService) HashPassword(password string) (string, error) {
	config := &Argon2Configuration{
		TimeCost:   2,
		MemoryCost: 64 * 1024,
		Threads:    4,
		KeyLength:  32,
	}

	salt, err := s.RandomService.GenerateSecretKey(16)
	if err != nil {
		return "", fmt.Errorf("password hashing failed: %w", err)
	}
	config.Salt = salt

	// Execute Argon2id hashing algorithm
	config.HashRaw = argon2.IDKey(
		[]byte(password),
		config.Salt,
		config.TimeCost,
		config.MemoryCost,
		config.Threads,
		config.KeyLength,
	)

	// Generate standardized hash format
	encodedHash := fmt.Sprintf(
		"$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s",
		argon2.Version,
		config.MemoryCost,
		config.TimeCost,
		config.Threads,
		base64.RawStdEncoding.EncodeToString(config.Salt),
		base64.RawStdEncoding.EncodeToString(config.HashRaw),
	)

	return encodedHash, nil
}

func (s *PasswordService) parseArgon2Config(encodedHash string) (*Argon2Configuration, error) {
	components := strings.Split(encodedHash, "$")
	if len(components) != 6 {
		return nil, errors.New("invalid hash format structure")
	}

	// Validate algorithm identifier
	if !strings.HasPrefix(components[1], "argon2id") {
		return nil, errors.New("unsupported algorithm variant")
	}

	// Extract version information
	var version int
	_, err := fmt.Sscanf(components[2], "v=%d", &version)
	if err != nil {
		return nil, err
	}

	// Parse configuration parameters
	config := &Argon2Configuration{}
	_, err = fmt.Sscanf(components[3], "m=%d,t=%d,p=%d",
		&config.MemoryCost, &config.TimeCost, &config.Threads)
	if err != nil {
		return nil, err
	}

	// Decode salt component
	salt, err := base64.RawStdEncoding.DecodeString(components[4])
	if err != nil {
		return nil, fmt.Errorf("salt decoding failed: %w", err)
	}
	config.Salt = salt

	// Decode hash component
	hash, err := base64.RawStdEncoding.DecodeString(components[5])
	if err != nil {
		return nil, fmt.Errorf("hash decoding failed: %w", err)
	}
	config.HashRaw = hash
	config.KeyLength = uint32(len(hash))

	return config, nil
}

// VerifyPassword verifies a password against a hash.
func (s *PasswordService) VerifyPassword(storedHash string, providedPassword string) (bool, error) {
	// Parse stored hash parameters
	config, err := s.parseArgon2Config(storedHash)
	if err != nil {
		return false, fmt.Errorf("hash parsing failed: %w", err)
	}

	// Generate hash using identical parameters
	computedHash := argon2.IDKey(
		[]byte(providedPassword),
		config.Salt,
		config.TimeCost,
		config.MemoryCost,
		config.Threads,
		config.KeyLength,
	)

	// Perform constant-time comparison to prevent timing attacks
	match := subtle.ConstantTimeCompare(config.HashRaw, computedHash) == 1
	return match, nil
}
