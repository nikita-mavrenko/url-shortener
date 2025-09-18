package utils

import (
	"math/rand"
	"time"
)

func GenerateAlphabet(length int) []rune {
	rnd := rand.New(rand.NewSource(time.Now().UnixMilli()))

	alphabet := make([]rune, length)

	chars := []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789")

	for i := range alphabet {
		alphabet[i] = chars[rnd.Intn(len(chars))]
	}

	return alphabet
}
