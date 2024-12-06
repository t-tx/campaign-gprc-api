package pkg

import (
	"crypto/rand"
	"math/big"
)

func GenerateUniqueID(length int) (string, error) {
	// Define the alphabet (uppercase and lowercase letters)
	const alphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

	// Create a slice to store the random characters
	result := make([]byte, length)

	// Randomly pick characters from the alphabet
	for i := 0; i < length; i++ {
		// Generate a random number to pick an index from the alphabet string
		randomIndex, err := rand.Int(rand.Reader, big.NewInt(int64(len(alphabet))))
		if err != nil {
			return "", err
		}
		// Assign a random character from the alphabet to the result slice
		result[i] = alphabet[randomIndex.Int64()]
	}

	// Convert the slice to a string and return it
	return string(result), nil
}
