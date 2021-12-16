package strfuncs

import (
	"math/rand"
	"strings"
)

func ReplaceMap(str string, repl map[string]string) string {
	for k, v := range repl {
		str = strings.ReplaceAll(str, k, v)
	}

	return str
}

func RandomString(n int) string {
	var letter = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

	b := make([]rune, n)
	for i := range b {
		b[i] = letter[rand.Intn(len(letter))]
	}
	return string(b)
}
