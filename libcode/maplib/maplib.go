package maplib

import "errors"

func Choice(m map[string]string) (string, error) {
	for _, v := range m {
		return v, nil
	}

	return "", errors.New("the map is empty")
}
