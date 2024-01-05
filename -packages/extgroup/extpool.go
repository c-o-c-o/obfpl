package extgroup

type ExtPool struct {
	// pool[name][ext] = "filePath"
	pool map[string]map[string]string
}

type PopFiles struct {
	Name     string
	ExtGroup map[string]string
}

func NewExtPool() *ExtPool {
	return &ExtPool{
		pool: map[string]map[string]string{},
	}
}

/*
ext group が全て揃っていた場合はそれを返します
*/
func (ep *ExtPool) Push(filePath string, ext string, extList []string) *PopFiles {

	name := GetName(filePath)
	if _, ok := ep.pool[name]; ok {
		ep.pool[name][ext] = filePath
	} else {
		ep.pool[name] = map[string]string{ext: filePath}
	}

	//本来は ep.pool[name][x] == extList[x] でチェックしないといけない
	//しかし ext が不正な値で無ければこれで十分なのと面倒なので数をチェックだけ
	if len(ep.pool[name]) < len(extList) {
		return nil
	}

	extGroup := ep.pool[name]
	delete(ep.pool, name)

	return &PopFiles{
		Name:     name,
		ExtGroup: extGroup,
	}
}
