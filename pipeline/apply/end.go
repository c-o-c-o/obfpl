package apply

import (
	"obfpl/libcode/pathlib"
	"obfpl/libcode/storage"
	"obfpl/pipeline/apply/process"
	"os"
	"path/filepath"
)

func End(ctx *ApplyContext, outPath string, erch chan error) {
	src := ctx.vars["src"]
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

	vars := make(map[string]string, len(ctx.group))
	for k, v := range ctx.group {
		dstp := filepath.Join(outPath, pathlib.AddSuffix(v, sf))
		err := os.Rename(
			filepath.Join(src, v),
			dstp)
		if err != nil {
			erch <- err
			return
		}

		vars[k] = dstp

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
		err = process.CallExec(v, vars)
		if err != nil {
			erch <- err
		}
	}

	err = os.RemoveAll(filepath.Join(ctx.profile.Env.Temp, ctx.temp))
	if err != nil {
		erch <- err
		return
	}
}
