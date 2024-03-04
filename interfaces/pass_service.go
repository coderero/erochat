package interfaces

type PassService interface {
	// Hash hashes a password.
	Hash(password string) (string, error)

	// Compare compares a password with a hash.
	Compare(password, hash string) bool
}
