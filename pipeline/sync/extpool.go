package sync

import (
	"obfpl/libcode/pathlib"
	"obfpl/pipeline/apply/find"
	"path/filepath"
)

type ExtPool struct {
	ext  map[string]string
	list map[string]map[string]string
}

func NewExtPool(ext map[string]string) *ExtPool {
	list := make(map[string]map[string]string)

	return &ExtPool{
		ext:  ext,
		list: list,
	}
}

func (ep *ExtPool) Add(filename string) (map[string]string, error) {
	ext, err := find.MatchExt(ep.ext, filepath.Ext(filename)[1:])
	if err != nil {
		return nil, err
	}

	name := pathlib.WithoutExt(filename)
	_, exist := ep.list[name]
	if !exist {
		ep.list[name] = make(map[string]string)
	}

	ep.list[name][ext] = filename

	if len(ep.list[name]) == len(ep.ext) {
		group := ep.list[name]
		delete(ep.list, name)

		return group, nil
	}

	return nil, nil
}
