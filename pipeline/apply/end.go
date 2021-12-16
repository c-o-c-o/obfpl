package apply

import (
	"obfpl/libcode/pathlib"
	"obfpl/libcode/storage"
	"obfpl/libcode/strfuncs"
	"obfpl/pipeline/process"
	"os"
	"path/filepath"
)

func End(ctx *ApplyContext, outPath string, erch chan error) {
	src := ctx.repList["{@src}"]
	err := os.MkdirAll(outPath, 0777)
	if err != nil {
		erch <- err
		return
	}

	fl := make([]string, 0, len(ctx.group))
	for _, v := range ctx.group {
		fl = append(fl, v)
	}

	sf := storage.GetFileNameSuffix(outPath, fl)

	repl := make(map[string]string, len(ctx.group))
	for k, v := range ctx.group {
		dstp := filepath.Join(outPath, pathlib.AddSuffix(v, sf))
		err := os.Rename(
			filepath.Join(src, v),
			dstp)
		if err != nil {
			erch <- err
			return
		}

		repl["{@"+k+"}"] = dstp

		err = storage.UpdateCreateTime(dstp, ctx.ctimes[k])
		if err != nil {
			erch <- err
		}

		//ログ出力
		if ctx.Loging != nil {
			ctx.Loging(v)
		}
	}

	for _, v := range ctx.profile.Notify {
		err = process.CallExec(strfuncs.ReplaceMap(v, repl))
		if err != nil {
			erch <- err
		}
	}

	err = os.RemoveAll(filepath.Join(ctx.profile.Env.Temp, ctx.name))
	if err != nil {
		erch <- err
		return
	}
}
