package utils

import "math/rand"

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyz123456789")

func InsecureRandom(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}
