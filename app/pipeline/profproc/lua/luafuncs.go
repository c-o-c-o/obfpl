package lua

import (
	"errors"
	"obfpl/-packages/dict"
	"obfpl/-packages/exec"
	"obfpl/-packages/sync"
	"obfpl/app/pipeline/profproc/-packages/context"
	"obfpl/app/pipeline/profproc/-packages/util"
	"obfpl/app/pipeline/profproc/lua/lua"
	"path/filepath"
)

func GetSetupMethods(profile *Profile) []*lua.Function {
	return []*lua.Function{
		lua.NewFunction("SetEnv", func(flow *lua.Flow) error {
			k := flow.GetString(0)
			v := flow.GetString(1)

			switch k {
			case "temp":
				temp, err := filepath.Abs(v)
				if err != nil {
					return err
				}
				profile.tempPath = temp
			case "exec-rule":
				profile.isAsync = v == "async"
			default:
				return errors.New("an undefined key \"" + k + "\" was specified")
			}

			return nil
		}),
		lua.NewFunction("SetVar", func(flow *lua.Flow) error {
			profile.vari[flow.GetString(0)] = flow.GetString(1)
			return nil
		}),
		lua.NewFunction("SetExt", func(flow *lua.Flow) error {
			profile.exts[flow.GetString(0)] = flow.GetString(1)
			return nil
		}),
		lua.NewFunction("SetName", func(flow *lua.Flow) error {
			profile.basis = flow.GetString(0)
			return nil
		}),
	}
}

func GetProcessMethods(profProc *LuaProfProc, waiter *sync.Waiter, ctx *context.Context) []*lua.Function {
	return []*lua.Function{
		lua.NewFunction("Match", func(flow *lua.Flow) error {
			ptn := flow.GetString(0)
			path := flow.GetString(1)
			enc := flow.GetString(2)

			val := ctx.Vari.GetBasisValue(profProc.profile.basis)
			matched, err := util.Match(ptn, val, path, enc)
			if err != nil {
				return err
			}

			flow.ReturnBool(matched)
			return nil
		}),
		lua.NewFunction("Wait", func(flow *lua.Flow) error {
			waiter.Wait()
			return nil
		}),
		lua.NewFunction("Execute", func(flow *lua.Flow) error {
			cmd := flow.GetString(0)

			fileNames, err := ctx.Temp.LoadFileNames("src")
			if err != nil {
				return err
			}

			group, err := dict.FromArrayE(fileNames, func(val string) (string, string, error) {
				k, err := profProc.SelectExt(val)
				return k, val, err
			})
			if err != nil {
				return err
			}

			src, dst := ctx.Temp.GetPaths()
			ctx.Vari.Update(map[string]string{
				"src":  src,
				"dst":  dst,
				"out":  profProc.outPath,
				"name": util.GetBasisValue(group, profProc.profile.basis),
			}, group)

			err = exec.Call(ctx.Vari.Apply(cmd))
			if err != nil {
				return err
			}

			err = ctx.Temp.Exchange(util.CreateGetMoveList(profProc.SelectExt))
			if err != nil {
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
	}
}

func GetNotifyMethods(profProc *LuaProfProc, ctx *context.Context) []*lua.Function {
	return []*lua.Function{
		lua.NewFunction("Match", func(flow *lua.Flow) error {
			ptn := flow.GetString(0)
			path := flow.GetString(1)
			enc := flow.GetString(2)

			val := ctx.Vari.GetBasisValue(profProc.profile.basis)
			matched, err := util.Match(ptn, val, path, enc)
			if err != nil {
				return err
			}

			flow.ReturnBool(matched)
			return nil
		}),
		lua.NewFunction("Execute", func(flow *lua.Flow) error {
			cmd := flow.GetString(0)

			err := exec.Call(ctx.Vari.Apply(cmd))
			if err != nil {
				return err
			}

			return nil
		}),
	}
}
