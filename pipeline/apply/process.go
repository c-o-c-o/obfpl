package apply

import (
	"obfpl/analyze"
	"obfpl/data"
	"obfpl/libcode/storage"
	"obfpl/pipeline/apply/find"
	"obfpl/pipeline/apply/process"
	"os"
	"path/filepath"
)

func Call(ctx *Context, proc data.Process) error {
	err := process.Call(analyze.VerReflection(proc.Cmd, ctx.vars))
	if err != nil {
		return err
	}

	src := ctx.vars["src"]
	dst := ctx.vars["dst"]
	group, err := moveGroupFiles(ctx, src, dst)
	if err != nil {
		return err
	}

	err = storage.RemoveDirInner(src)
	if err != nil {
		return err
	}

	calledUpdate(ctx, group)
	return nil
}

func cleanup(ctx *Context) error {
	return os.RemoveAll(filepath.Join(filepath.Join(ctx.profile.Env.Temp, ctx.name)))
}

func Notify(ctx *Context, cmd string) error {
	err := process.Call(analyze.VerReflection(cmd, ctx.vars))
	if err != nil {
		return err
	}

	return nil
}

func moveGroupFiles(ctx *Context, src string, dst string) (map[string]string, error) {
	sgroup, err := find.File(src, ctx.ext)
	if err != nil {
		return nil, err
	}

	dgroup, err := find.File(dst, ctx.ext)
	if err != nil {
		return nil, err
	}

	for k := range ctx.group {
		_, exist := dgroup[k]
		if exist {
			continue
		}

		srcfn := sgroup[k]
		err := os.Rename(filepath.Join(src, srcfn), filepath.Join(dst, srcfn))
		if err != nil {
			return nil, err
		}
		dgroup[k] = srcfn
	}
	return dgroup, nil

}
