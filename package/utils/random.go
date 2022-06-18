package utils

import "math/rand"

const alphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"


func RandomString(length int) string{
	b := make([]byte, length)
	for i := range b{
		b[i] = alphabet[rand.Intn(len(alphabet))]
	}
	return string(b)
}