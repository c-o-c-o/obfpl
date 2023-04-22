package exts

import (
	"path/filepath"
)

type ExtsPool struct {
	pool   map[string]map[string]string
	getExt func(file string) (string, int, error)
}

type StartingFiles struct {
	Name  string
	Files map[string]string
}

func NewExtsPool(getExt func(string) (string, int, error)) *ExtsPool {
	return &ExtsPool{
		pool:   map[string]map[string]string{},
		getExt: getExt,
	}
}

func (ep *ExtsPool) Add(file string) (*StartingFiles, error) {
	e, l, err := ep.getExt(file)
	if err != nil {
		return nil, err
	}

	n := ep.GetName(file)
	if _, ok := ep.pool[n]; ok {
		ep.pool[n][e] = file
	} else {
		ep.pool[n] = map[string]string{e: file}
	}

	if len(ep.pool[n]) == l {
		fs := ep.pool[n]
		delete(ep.pool, n)

		return &StartingFiles{
			Name:  n,
			Files: fs,
		}, nil
	}

	return nil, nil
}

func (ep *ExtsPool) GetName(file string) string {
	e := filepath.Ext(file)
	n := filepath.Base(file)
	return n[:len(n)-len(e)]
}
