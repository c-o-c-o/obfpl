package apply

import (
	"obfpl/libcode/pathlib"
	"obfpl/libcode/storage"
	"os"
	"path/filepath"
	"syscall"
)

func Output(ctx *Context) error {
	outlist := make(map[string]string, len(ctx.group))

	err := os.MkdirAll(ctx.outpath, 0777)
	if err != nil {
		return err
	}

	sf := getSuffix(ctx)
	src := ctx.swap.list[ctx.swap.idx%len(ctx.swap.list)]

	for k, v := range ctx.group {
		outpath := filepath.Join(ctx.outpath, pathlib.AddSuffix(v, sf))
		err := outFile(filepath.Join(src, v), outpath, ctx.ctimes[k])
		if err != nil {
			return err
		}

		outlist[k] = outpath

		if ctx.Loging != nil {
			ctx.Loging(v)
		}
	}

	outedUpdate(ctx, outlist)
	return cleanup(ctx)
}

func getSuffix(ctx *Context) string {
	fl := make([]string, 0, len(ctx.group))
	for _, v := range ctx.group {
		fl = append(fl, v)
	}

	sf := storage.GetFileNameSuffix(ctx.outpath, fl)
	return sf
}

func outFile(srcpath string, outpath string, time syscall.Filetime) error {
	err := os.Rename(
		srcpath,
		outpath)
	if err != nil {
		return err
	}

	err = storage.UpdateCreateTime(outpath, time)
	if err != nil {
		return err
	}

	return nil
}
