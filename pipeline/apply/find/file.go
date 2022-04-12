package find

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
)

func File(dirpath string, ext map[string]string) (map[string]string, error) {
	dir, err := os.ReadDir(dirpath)
	if err != nil {
		return nil, err
	}

	files := make(map[string]string)
	for _, entry := range dir {
		k, err := MatchExt(ext, filepath.Ext(entry.Name())[1:])
		if err != nil {
			continue
		}
		files[k] = entry.Name()
	}

	return files, nil
}

func MatchExt(extfilter map[string]string, ext string) (string, error) {
	if strings.Contains(ext, ",") {
		return "", errors.New("the extension did not match")
	}

	for k, v := range extfilter {
		if strings.Contains(v, ext) {
			return k, nil
		}
	}

	return "", errors.New("the extension did not match")
}
