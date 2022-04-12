package apply

import (
	"errors"
	"obfpl/libcode/maplib"
)

func getGroupValue(group map[string]string, name string) (string, error) {
	if name == "" {
		return maplib.Choice(group)
	}

	str, exist := group[name]
	if !exist {
		return "", errors.New("a name that does not exist in ext was specified")
	}

	return str, nil
}
