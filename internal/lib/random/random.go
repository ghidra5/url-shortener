package random

import (
	"math/rand"
	"time"
)

// TODO: change to slice of runes
const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

var rnd = rand.New(rand.NewSource(time.Now().UnixNano()))

func NewRandomString(length int) string {

	res := make([]byte, length)

	for i := 0; i < length; i++ {
		res[i] = charset[rnd.Intn(len(charset))]
	}

	return string(res)
}
