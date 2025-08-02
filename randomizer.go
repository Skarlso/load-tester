package main

import (
	"crypto/rand"
	"fmt"
	"math/big"
)

const charset = "abcdefghijklmnopqrstuvwxyz0123456789"

// GenerateRandomString creates a random string of specified length
func GenerateRandomString(length int) string {
	result := make([]byte, length)
	for i := range result {
		num, _ := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		result[i] = charset[num.Int64()]
	}
	return string(result)
}

// GenerateRandomName creates a randomized name with prefix
func GenerateRandomName(prefix string, suffixLength int) string {
	suffix := GenerateRandomString(suffixLength)
	return fmt.Sprintf("%s-%s", prefix, suffix)
}
