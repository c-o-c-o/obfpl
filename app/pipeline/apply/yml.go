package apply

import (
	"obfpl/app/pipeline/apply/sync"
	"obfpl/app/pipeline/apply/yml"
	"path/filepath"
	"strings"
)

type YmlApply struct {
	exts    map[string]string
	vari    map[string]string
	isAsync bool
	out     string
	temp    string
	basis   string
	notify  []string
	procs   []yml.Process
}

func NewYmlApply(ppath string, dpath string) (*YmlApply, error) {
	p, err := yml.LoadProfile(ppath)
	if err != nil {
		return nil, err
	}

	temp, err := filepath.Abs(p.Env.Temp)
	if err != nil {
		return nil, err
	}

	return &YmlApply{
		exts:    p.Ext,
		out:     dpath,
		vari:    p.Var,
		isAsync: p.Env.ExecRule == "async",
		temp:    temp,
		basis:   p.Name,
		notify:  p.Notify,
		procs:   p.Proc,
	}, nil
}

func (ya *YmlApply) IsAsync() bool {
	return ya.isAsync
}

func (ya *YmlApply) GetExt(file string) (string, int, error) {
	e, err := ya.getExt(ya.exts, file)
	return e, len(ya.exts), err
}

func (ya *YmlApply) getExt(exts map[string]string, file string) (string, error) {
	for key, sense := range exts {
		e := filepath.Ext(file)
		if len(e) <= 1 {
			continue
		}
		if strings.Contains(sense, e[1:]) {
			return key, nil
		}
	}
	return "", nil
}

func (ya *YmlApply) Run(files *StartingFiles, waiter *sync.Waiter) {
	defer waiter.Destroy()

	ctx, err := NewContext(files, ya.temp, ya.exts, ya.vari)
	if err != nil {
		waiter.Error(err)
		return
	}
	defer ctx.Cleanup(func(err error) {
		waiter.Error(err)
	})

	//プロセス
	sg, err := ctx.Temp.GetGroup("src", WithExts(ctx.Exts, ya.getExt))
	if err != nil {
		waiter.Error(err)
		return
	}

	for _, p := range ya.procs {
		name := GetName(sg, ya.basis)

		isCall, err := Match(p.Ptn, name, p.Trg, p.Enc)
		if err != nil {
			waiter.Error(err)
			return
		}
		if !isCall {
			continue
		}

		src, dst := ctx.Temp.GetPaths()
		ctx.Var.Update(map[string]string{
			"src":  src,
			"dst":  dst,
			"out":  ya.out,
			"name": name,
		}, sg)

		if p.IsWait {
			waiter.Wait()
		}

		if err := Call(ctx.Var.Apply(p.Cmd)); err != nil {
			waiter.Error(err)
			return
		}

		if p.Ext != nil {
			ctx.Exts = p.Ext
		}

		sg, err = ctx.Temp.Exchange(WithExts(ctx.Exts, ya.getExt))
		if err != nil {
			waiter.Error(err)
			return
		}
	}

	if err := ctx.Temp.Output(ya.out, toArrayFromValue(sg)); err != nil {
		waiter.Error(err)
		return
	}

	for _, file := range sg {
		waiter.Log(file)
	}

	exts := map[string]string{}

	for k, v := range sg {
		exts[k] = filepath.Join(ya.out, v)
	}

	ctx.Var.Update(map[string]string{}, exts)

	//通知
	for _, n := range ya.notify {
		if err := Call(ctx.Var.Apply(n)); err != nil {
			waiter.Error(err)
			return
		}
	}

	waiter.Close()
}

func toArrayFromValue[Value any](dict map[string]Value) []Value {
	array := []Value{}
	for _, v := range dict {
		array = append(array, v)
	}
	return array
}
