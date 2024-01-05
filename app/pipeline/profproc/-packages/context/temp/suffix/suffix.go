package suffix

import (
	"errors"
	"math/rand"
	"os"
	"path/filepath"
	"time"
)

var seed *rand.Rand = nil

func init() {
	seed = rand.New(rand.NewSource(time.Now().UnixMilli()))
}

func With(path string, suffix string) string {
	ext := filepath.Ext(path)
	return path[:len(path)-len(ext)] + suffix + ext
}

func TryWithFiles(base string, files []string, suffix string) error {
	for _, file := range files {
		path := filepath.Join(base, With(file, suffix))
		info, _ := os.Stat(path)
		if info != nil {
			return errors.New("file already exists. : " + info.Name())
		}
	}
	return nil
}

func Get(len int) string {
	return randomString(len)
}

func randomString(n int) string {
	var letter = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

	b := make([]rune, n)
	for i := range b {
		b[i] = letter[seed.Intn(len(letter))]
	}
	return string(b)
}
