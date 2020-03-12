package main

import (
	"math/rand"
	"time"
)

const base64 = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789-_"

// GenerateID generates an n-character long ID based on Base64.
func GenerateID(n int) string {
	rand.Seed(time.Now().UnixNano())

	alphabet := []rune(base64)
	result := ""

	for i := 0; i < n; i++ {
		randomChar := alphabet[rand.Intn(64)]
		result += string(randomChar)
	}

	return result
}
