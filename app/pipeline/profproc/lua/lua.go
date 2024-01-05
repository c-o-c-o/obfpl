package lua

import (
	"obfpl/-packages/dict"
	"obfpl/-packages/sync"
	"obfpl/app/pipeline/profproc/-packages/context"
	"obfpl/app/pipeline/profproc/lua/lua"
	"os"
	"path/filepath"
	"regexp"
)

type Profile struct {
	isAsync  bool
	tempPath string
	basis    string
	vari     map[string]string
	exts     map[string]string
}

type LuaProfProc struct {
	script  string
	outPath string
	profile *Profile
}

func NewLuaProfProc(profilePath string, outPath string) (*LuaProfProc, error) {
	p := &LuaProfProc{
		outPath: outPath,
		profile: &Profile{
			vari: map[string]string{},
			exts: map[string]string{},
		},
	}

	//スクリプト読み込み
	bin, err := os.ReadFile(profilePath)
	if err != nil {
		return nil, err
	}

	p.script = string(bin)
	//----------------------------------------------------------------

	err = lua.DoScript(p.script+" Setup()", GetSetupMethods(p.profile))
	if err != nil {
		return nil, err
	}

	return p, nil
}

func (p *LuaProfProc) GetType() string {
	return "Lua"
}

func (p *LuaProfProc) SelectExt(filePath string) (string, error) {
	for ext, ptn := range p.profile.exts {
		ok, err := regexp.Match(ptn, []byte(filePath))

		if err != nil {
			return "", err
		}
		if ok {
			return ext, nil
		}
	}

	return "", nil
}

func (p *LuaProfProc) CreateExtList() []string {
	return dict.Keys(p.profile.exts)
}

func (p *LuaProfProc) Call(name string, extGroup map[string]string, waiter *sync.Waiter) {
	if p.profile.isAsync {
		go p.call(name, extGroup, waiter)
	} else {
		p.call(name, extGroup, waiter)
	}
}

func (p *LuaProfProc) call(name string, extGroup map[string]string, waiter *sync.Waiter) {
	defer waiter.Destroy()

	ctx, err := context.NewContext(
		name,
		dict.Values(extGroup),
		p.profile.tempPath,
		p.profile.exts,
		p.profile.vari)
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

	err = lua.DoScript(p.script+" Process()", GetProcessMethods(p, waiter, ctx))
	if err != nil {
		waiter.Error(err)
		return
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

	err = lua.DoScript(p.script+" Notify()", GetNotifyMethods(p, ctx))
	if err != nil {
		waiter.Error(err)
		return
	}

	waiter.Close()
}
