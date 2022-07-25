package util

import (
	"math/rand"
	"time"
)

const (
	base32Encoder = `ABCDEFGHIJKLMNOPQRSTUVWXYZ234567`
)

// randString 随机字符串
func randString(encoder string, length int) string {
	rand.Seed(time.Now().UnixNano())
	letters := []rune(encoder)
	b := make([]rune, length)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

// RandBase32String 指定长度随机base32字符串
func RandBase32String(length int) string {
	return randString(base32Encoder, length)
}
