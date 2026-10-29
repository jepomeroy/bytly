package utils

import (
	"math/rand"
	"time"
)

var runes = []rune("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func RandomUrl(size int) string {
	str := make([]rune, size)

	// time.Now().UnixNano(), which yields a constantly-changing number.
	rand.Seed(time.Now().UnixNano())

	for i := range str {
		str[i] = runes[rand.Intn(len(runes))]
	}

	return string(str)
}
