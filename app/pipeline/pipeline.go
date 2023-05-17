package pipeline

import (
	"errors"
	"obfpl/app/pipeline/apply"
	"obfpl/app/pipeline/apply/sync"
	"obfpl/app/pipeline/exts"
	"path/filepath"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type Pipeline struct {
	Type  string
	apply apply.Apply
}

/*
ppath プロファイルパス

dpath 出力先パス
*/
func NewPipeline(ppath string, dpath string) (*Pipeline, error) {
	p := &Pipeline{}
	p.Type = cases.Title(language.English).String(filepath.Ext(ppath)[1:])

	var err error = nil
	switch p.Type {
	case "Lua":
		p.apply, err = apply.NewLuaApply(ppath, dpath)
	case "Yaml", "Yml":
		p.apply, err = apply.NewYmlApply(ppath, dpath)
	default:
		err = errors.New("please specify the correct profile")
	}

	if err != nil {
		return nil, err
	}

	return p, nil
}

/*
ファイルの拡張子一致チェック

指定したファイルが揃うとProcess開始
*/
func (p Pipeline) CreateChecking(msgch chan<- string) func(string) {
	waiter := sync.NewClosedWaiter(msgch)
	extsPool := exts.NewExtsPool(p.apply.GetExt)

	return func(file string) {
		sf, err := extsPool.Add(file)
		if err != nil {
			msgch <- err.Error()
			return
		}
		if sf == nil {
			return
		}

		waiter = waiter.Next()
		asf := &apply.StartingFiles{
			Name:  sf.Name,
			Files: sf.Files,
		}

		if p.apply.IsAsync() {
			go p.apply.Run(asf, waiter)
		} else {
			p.apply.Run(asf, waiter)
		}
	}
}
