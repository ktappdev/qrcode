package helpers

import (
	"math/rand"
	"time"
)

const charset = "abcdefghijklmnopqrstuvwxyz0123456789"

func GenerateRandomString(length int) string {
	var result string
	seededRand := rand.New(rand.NewSource(time.Now().UnixNano()))

	for i := 0; i < length; i++ {
		index := seededRand.Intn(len(charset))
		result += string(charset[index])
	}
	return result
}
