package apply

import (
	"errors"
	"obfpl/app/pipeline/apply/lua"
	"obfpl/app/pipeline/sync"
	"os"
	"path/filepath"
	"regexp"
)

type LuaApply struct {
	exts    map[string]string
	vari    map[string]string
	isAsync bool
	out     string
	temp    string
	basis   string
	script  string
}

func NewLuaApply(ppath string, dpath string) (*LuaApply, error) {
	la := &LuaApply{
		exts:    map[string]string{},
		out:     dpath,
		vari:    map[string]string{},
		isAsync: false,
	}

	bin, err := os.ReadFile(ppath)
	if err != nil {
		return nil, err
	}

	la.script = string(bin)

	if err := lua.DoScript(la.script+" Setup()", []*lua.Function{
		lua.NewFunction("SetEnv", func(flow *lua.Flow) error {
			k := flow.GetString(0)
			v := flow.GetString(1)

			switch k {
			case "temp":
				temp, err := filepath.Abs(v)
				if err != nil {
					return err
				}
				la.temp = temp
			case "exec-rule":
				la.isAsync = v == "async"
			default:
				return errors.New("an undefined key \"" + k + "\" was specified")
			}

			return nil
		}),
		lua.NewFunction("SetVar", func(flow *lua.Flow) error {
			la.vari[flow.GetString(0)] = flow.GetString(1)
			return nil
		}),
		lua.NewFunction("SetExt", func(flow *lua.Flow) error {
			la.exts[flow.GetString(0)] = flow.GetString(1)
			return nil
		}),
		lua.NewFunction("SetName", func(flow *lua.Flow) error {
			la.basis = flow.GetString(0)
			return nil
		}),
	}); err != nil {
		return nil, err
	}

	return la, nil
}

func (la *LuaApply) IsAsync() bool {
	return la.isAsync
}

func (la *LuaApply) GetExt(file string) (string, int, error) {
	e, err := la.getExt(la.exts, file)
	return e, len(la.exts), err
}

func (la *LuaApply) getExt(exts map[string]string, file string) (string, error) {
	for k, ptn := range exts {
		med, err := regexp.Match(ptn, []byte(file))
		if err != nil {
			return "", err
		}
		if med {
			return k, nil
		}
	}

	return "", nil
}

func (la *LuaApply) Run(files *StartingFiles, waiter *sync.Waiter) {
	defer waiter.Destroy()

	ctx, err := NewContext(files, la.temp, la.exts, la.vari)
	if err != nil {
		waiter.Error(err)
		return
	}
	defer ctx.Cleanup(func(err error) {
		waiter.Error(err)
	})

	if err := lua.DoScript(la.script+" Process()", []*lua.Function{
		lua.NewFunction("Match", func(flow *lua.Flow) error {
			ptn := flow.GetString(0)
			path := flow.GetString(1)
			enc := flow.GetString(2)

			sg, err := ctx.Temp.GetGroup("src", WithExts(ctx.Exts, la.getExt))
			if err != nil {
				return err
			}
			name := GetName(sg, la.basis)

			isMatch, err := Match(ptn, name, path, enc)
			if err != nil {
				return err
			}

			flow.ReturnBool(isMatch)
			return nil
		}),
		lua.NewFunction("Wait", func(flow *lua.Flow) error {
			waiter.Wait()
			return nil
		}),
		lua.NewFunction("Execute", func(flow *lua.Flow) error {
			cmd := flow.GetString(0)

			src, dst := ctx.Temp.GetPaths()
			sg, err := ctx.Temp.GetGroup("src", WithExts(ctx.Exts, la.getExt))
			if err != nil {
				return err
			}
			ctx.Var.Update(map[string]string{
				"src":  src,
				"dst":  dst,
				"out":  la.out,
				"name": GetName(sg, la.basis),
			}, sg)

			if err := Call(ctx.Var.Apply(cmd)); err != nil {
				return err
			}

			if _, err := ctx.Temp.Exchange(WithExts(ctx.Exts, la.getExt)); err != nil {
				return err
			}

			return nil
		}),
		lua.NewFunction("SetExt", func(flow *lua.Flow) error {
			ctx.Exts[flow.GetString(0)] = flow.GetString(1)
			return nil
		}),
		lua.NewFunction("ClearExt", func(flow *lua.Flow) error {
			ctx.Exts = map[string]string{}
			return nil
		}),
	}); err != nil {
		waiter.Error(err)
		return
	}

	sg, err := ctx.Temp.GetGroup("src", WithExts(ctx.Exts, la.getExt))
	if err != nil {
		waiter.Error(err)
		return
	}
	ctx.Temp.Output(la.out, toArrayFromValue(sg))

	for _, file := range sg {
		waiter.Log(file)
	}

	exts := map[string]string{}

	for k, v := range sg {
		exts[k] = filepath.Join(la.out, v)
	}

	ctx.Var.Update(map[string]string{}, exts)

	if err := lua.DoScript(la.script+" Notify()", []*lua.Function{
		lua.NewFunction("Match", func(flow *lua.Flow) error {
			ptn := flow.GetString(0)
			path := flow.GetString(1)
			enc := flow.GetString(2)

			sg, err := ctx.Temp.GetGroup("src", WithExts(ctx.Exts, la.getExt))
			if err != nil {
				return err
			}
			name := GetName(sg, la.basis)

			isMatch, err := Match(ptn, name, path, enc)
			if err != nil {
				return err
			}

			flow.ReturnBool(isMatch)
			return nil
		}),
		lua.NewFunction("Execute", func(flow *lua.Flow) error {
			cmd := flow.GetString(0)

			if err := Call(ctx.Var.Apply(cmd)); err != nil {
				return err
			}
			return nil
		}),
	}); err != nil {
		waiter.Error(err)
		return
	}

	waiter.Close()
}
