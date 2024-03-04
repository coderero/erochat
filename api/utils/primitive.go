package utils

import (
	"os"
)

func GetFile(filename string) ([]byte, error) {
	var privateKey []byte

	file, err := os.ReadFile(filename)
	if err != nil {
		return privateKey, err
	}

	privateKey = file
	return privateKey, nil
}
