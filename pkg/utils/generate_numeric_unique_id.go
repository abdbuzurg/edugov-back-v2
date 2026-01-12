package utils

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"strings"
)

func GenerateNumericUniqueID() (string, error) {
	// Initialize a strings.Builder for efficient string concatenation.
	// We know the final string length will be 9 (8 digits + 1 hyphen).
	var builder strings.Builder
	builder.Grow(9)

	// Loop 8 times to generate each of the 8 required digits.
	for i := 0; i < 8; i++ {
		// After the first 4 digits, append a hyphen to match the XXXX-XXXX format.
		if i == 4 {
			builder.WriteByte('-')
		}

		// Generate a random digit (0-9).
		// rand.Int(rand.Reader, big.NewInt(10)) generates a cryptographically
		// secure random integer in the range [0, 10), which means 0 through 9.
		// This method ensures a uniform distribution and avoids any statistical bias.
		digitBig, err := rand.Int(rand.Reader, big.NewInt(10))
		if err != nil {
			// If there's an error during random number generation, return an error.
			return "", fmt.Errorf("failed to generate random digit: %w", err)
		}

		// Convert the *big.Int result to an int64 and then to a string,
		// appending it to our string builder.
		builder.WriteString(fmt.Sprintf("%d", digitBig.Int64()))
	}

	// Return the final generated string and no error (nil).
	return builder.String(), nil
}
