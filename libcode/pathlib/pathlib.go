package pathlib

import "path/filepath"

func AddSuffix(path string, sf string) string {
	return WithoutExt(path) + sf + filepath.Ext(path)
}

func WithoutExt(path string) string {
	return filepath.Base(path[:len(path)-len(filepath.Ext(path))])
}
