package apply

import (
	"errors"
	"obfpl/libcode/maplib"
	"obfpl/libcode/storage"
	"obfpl/pipeline/process"
	"os"
	"path/filepath"
	"strings"
)

func Run(ctx *ApplyContext, erch chan error) bool {
	callc := 0
	for _, v := range ctx.profile.Proc {
		ctx.repList = makeReplaceList(ctx, callc)
		src := ctx.repList["{@src}"]
		dst := ctx.repList["{@dst}"]

		iscall, err := process.Call(maplib.Choice(ctx.group), v, ctx.repList)
		if err != nil {
			erch <- errors.New(err.Error())
		}
		if !iscall {
			continue
		}
		callc += 1

		ctx.group, err = complDstGroup(ctx, src, dst)
		if err != nil {
			erch <- err
			return true
		}

		storage.RemoveDirInner(src, func(e error) { erch <- e })
	}

	ctx.repList = makeReplaceList(ctx, callc)
	return false
}

func makeReplaceList(ctx *ApplyContext, callc int) map[string]string {
	repl := make(map[string]string, len(ctx.group)+3)

	repl["{@name}"] = ctx.name
	repl["{@src}"] = ctx.swaps[(callc)%len(ctx.swaps)]
	repl["{@dst}"] = ctx.swaps[(callc+1)%len(ctx.swaps)]

	for k, v := range ctx.group {
		repl["{@"+k+"}"] = v
	}

	return repl
}

func complDstGroup(ctx *ApplyContext, src string, dst string) (map[string]string, error) {
	sdetg, err := getGroup(src, ctx.profile.Ext)
	if err != nil {
		return nil, err
	}

	dgroup, err := getGroup(dst, ctx.profile.Ext)
	if err != nil {
		return nil, err
	}

	for k := range ctx.group {
		_, exist := dgroup[k]
		if exist {
			continue
		}

		srcfn := sdetg[k]
		err := os.Rename(filepath.Join(src, srcfn), filepath.Join(dst, srcfn))
		if err != nil {
			return nil, err
		}
		dgroup[k] = srcfn
	}
	return dgroup, nil
}

func getGroup(path string, ext map[string]string) (map[string]string, error) {
	dir, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}

	group := make(map[string]string)
	for _, v := range dir {
		k, err := MatchExt(ext, filepath.Ext(v.Name())[1:])
		if err != nil {
			continue
		}
		group[k] = v.Name()
	}
	return group, nil
}

func MatchExt(gfill map[string]string, ext string) (string, error) {
	if strings.Contains(ext, ",") {
		return "", errors.New("the extension did not match")
	}

	for k, v := range gfill {
		if strings.Contains(v, ext) {
			return k, nil
		}
	}

	return "", errors.New("the extension did not match")
}
