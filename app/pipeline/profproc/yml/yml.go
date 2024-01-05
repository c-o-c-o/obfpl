package yml

import (
	"obfpl/-packages/dict"
	"obfpl/-packages/exec"
	"obfpl/-packages/sync"
	"obfpl/app/pipeline/profproc/-packages/context"
	"obfpl/app/pipeline/profproc/-packages/util"
	"path/filepath"
	"strings"
)

type YmlProfProc struct {
	outPath string
	profile Profile
}

func NewYmlProfProc(profilePath string, outPath string) (*YmlProfProc, error) {
	p, err := LoadProfile(profilePath)
	if err != nil {
		return nil, err
	}

	temp, err := filepath.Abs(p.Env.Temp)
	if err != nil {
		return nil, err
	}
	p.Env.Temp = temp

	return &YmlProfProc{
		outPath: outPath,
		profile: *p,
	}, nil
}

func (p *YmlProfProc) GetType() string {
	return "Yaml"
}

func (p *YmlProfProc) SelectExt(filePath string) (string, error) {
	for ext, sense := range p.profile.Ext {
		e := filepath.Ext(filePath)

		if len(e) <= 1 {
			continue
		}

		if strings.Contains(sense, e[1:]) {
			return ext, nil
		}
	}

	return "", nil
}

func (p *YmlProfProc) CreateExtList() []string {
	return dict.Keys(p.profile.Ext)
}

func (p *YmlProfProc) Call(name string, extGroup map[string]string, waiter *sync.Waiter) {
	if p.profile.Env.ExecRule == "async" {
		go p.call(name, extGroup, waiter)
	} else {
		p.call(name, extGroup, waiter)
	}
}

func (p *YmlProfProc) call(name string, extGroup map[string]string, waiter *sync.Waiter) {
	defer waiter.Destroy()

	ctx, err := context.NewContext(
		name,
		dict.Values(extGroup),
		p.profile.Env.Temp,
		p.profile.Ext,
		p.profile.Var)
	if err != nil {
		waiter.Error(err)
		return
	}
	defer (func() {
		err := ctx.Temp.Cleanup()
		if err != nil {
			waiter.Error(err)
		}
	})()

	for _, proc := range p.profile.Proc {
		basisVal := util.GetBasisValue(p.profile.Ext, p.profile.Name)
		matched, err := util.Match(
			proc.Ptn,
			basisVal,
			proc.Trg,
			proc.Enc,
		)
		if err != nil {
			waiter.Error(err)
			return
		}
		if !matched {
			continue
		}

		fileNames, err := ctx.Temp.LoadFileNames("src")
		if err != nil {
			waiter.Error(err)
			return
		}

		group, err := dict.FromArrayE(fileNames, func(val string) (string, string, error) {
			k, err := p.SelectExt(val)
			return k, val, err
		})
		if err != nil {
			waiter.Error(err)
			return
		}

		src, dst := ctx.Temp.GetPaths()
		ctx.Vari.Update(map[string]string{
			"src":  src,
			"dst":  dst,
			"out":  p.outPath,
			"name": basisVal,
		}, group)

		if proc.IsWait {
			waiter.Wait()
		}

		err = exec.Call(ctx.Vari.Apply(proc.Cmd))
		if err != nil {
			waiter.Error(err)
			return
		}

		if proc.Ext != nil {
			ctx.Exts = proc.Ext
		}

		err = ctx.Temp.Exchange(util.CreateGetMoveList(p.SelectExt))
		if err != nil {
			waiter.Error(err)
			return
		}
	}

	fileNames, err := ctx.Temp.LoadFileNames("src")
	if err != nil {
		waiter.Error(err)
		return
	}

	err = ctx.Temp.Output(p.outPath, fileNames)
	if err != nil {
		waiter.Error(err)
		return
	}

	for _, fileName := range fileNames {
		waiter.Log(fileName)
	}

	group, err := dict.FromArrayE(fileNames, func(val string) (string, string, error) {
		k, err := p.SelectExt(val)
		return k, filepath.Join(p.outPath, val), err
	})
	if err != nil {
		waiter.Error(err)
		return
	}

	ctx.Vari.Update(map[string]string{}, group)

	for _, cmd := range p.profile.Notify {
		err := exec.Call(ctx.Vari.Apply(cmd))
		if err != nil {
			waiter.Error(err)
			return
		}
	}

	waiter.Close()
}
