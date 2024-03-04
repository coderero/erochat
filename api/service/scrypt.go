package service

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"fmt"
	"strings"

	"golang.org/x/crypto/scrypt"
)

// ScryptService represents a scrypt service.
type ScryptService struct {
	// N is the CPU/memory cost parameter.
	N int

	// R is the block size parameter.
	R int

	// P is the parallelization parameter.
	P int

	// KeyLen is the length of the derived key.
	KeyLen int

	// SaltLen is the length of the salt.
	SaltLen int
}

// NewScryptService creates a new ScryptService.
func NewScryptService(n, r, p, keyLen, saltLen int) *ScryptService {
	return &ScryptService{
		N:       n,
		R:       r,
		P:       p,
		KeyLen:  keyLen,
		SaltLen: saltLen,
	}
}

// Hash hashes a password.
func (s *ScryptService) Hash(password string) (string, error) {
	var (
		salt []byte
		err  error
		dk   []byte

		encodedSalt string
		encodedDK   string

		hash string
	)

	// Generate a random salt.
	if salt, err = s.generatesSalt(s.SaltLen); err != nil {
		return "", err
	}

	// Derive the key.
	dk, err = scrypt.Key([]byte(password), salt, s.N, s.R, s.P, s.KeyLen)
	if err != nil {
		return "", err
	}

	// Encode the salt and derived key.
	encodedSalt = base64.RawStdEncoding.EncodeToString(salt)
	encodedDK = base64.RawStdEncoding.EncodeToString(dk)

	// Concatenate the encoded salt and derived key.
	hash = fmt.Sprintf("%s:%s", encodedSalt, encodedDK)

	return hash, nil
}

// Compare compares a password with a hash.
func (s *ScryptService) Compare(password, hash string) bool {
	var (
		salt  []byte
		dk    []byte
		newDk []byte
		err   error
	)

	// Split the hash into the encoded salt and derived key.
	parts := strings.Split(hash, ":")
	if len(parts) != 2 {
		return false
	}

	// Decode the encoded salt and derived key.
	salt, err = base64.RawStdEncoding.DecodeString(parts[0])
	if err != nil {
		return false
	}
	dk, err = base64.RawStdEncoding.DecodeString(parts[1])
	if err != nil {
		return false
	}

	// Derive the key.
	newDk, err = scrypt.Key([]byte(password), salt, s.N, s.R, s.P, s.KeyLen)
	if err != nil {
		return false
	}

	// Compare the derived key with the decoded derived key.
	return subtle.ConstantTimeCompare(newDk, dk) == 1
}

// GenerateSalt generates a random salt.
func (s *ScryptService) generatesSalt(l int) ([]byte, error) {
	length := l * 3 / 4
	salt := make([]byte, length)
	_, err := rand.Read(salt)
	if err != nil {
		return nil, err
	}

	return salt, nil
}
