package utils

import (
	"fmt"
	"math/rand"
	"time"
)

const base = 1594972193
const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

// RandStringBytes returns a random alphanumeric bytes sequences (base62) of the requested size
func RandStringBytes(n uint) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

// RandomName create a random name with time based segment and random string separated by a dash
func RandomName(n uint) string {
	t := time.Now().Unix() - base
	return fmt.Sprintf("%d-%s", t, RandStringBytes(n))
}
