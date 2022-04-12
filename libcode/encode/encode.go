package encode

import (
	"errors"

	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"
)

func Convert(str string, enc string) (string, error) {
	conv, exist := map[string]encoding.Encoding{
		"":          japanese.ShiftJIS,
		"utf-8":     nil,
		"shift-jis": japanese.ShiftJIS,
	}[enc]
	if !exist {
		return "", errors.New("the specified encoding was not found")
	}

	if conv == nil {
		return str, nil
	}

	str, _, err := transform.String(conv.NewDecoder(), str)
	if err != nil {
		return "", err
	}

	return str, nil
}
