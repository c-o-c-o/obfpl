package apply

import (
	"errors"
	"obfpl/data"
	"obfpl/libcode/pathlib"
	"obfpl/libcode/retry"
	"obfpl/libcode/storage"
	"os"
	"path/filepath"
	"syscall"
	"time"
)

func CreateContext(pf *data.Profile, obspath string, outpath string, group map[string]string) (*Context, error) {
	fname, err := getGroupValue(group, pf.Name)
	if err != nil {
		return nil, err
	}
	name := pathlib.WithoutExt(fname)
	sname := storage.GetExistsName(pf.Env.Temp, name)
	swaps, err := storage.MakeDirList(filepath.Join(pf.Env.Temp, sname), []string{"a", "b"})
	if err != nil {
		return nil, err
	}

	ctimes, err := getCreateTimes(obspath, group)
	if err != nil {
		return nil, err
	}

	err = moveFiles(obspath, swaps[0], group)
	if err != nil {
		return nil, err
	}

	ctx := &Context{
		profile: pf,
		swap: Swap{
			list: swaps,
			idx:  0,
		},
		ext:     pf.Ext,
		group:   group,
		ctimes:  ctimes,
		name:    name,
		outpath: outpath,
	}

	return ctx, nil
}

func getCreateTimes(path string, files map[string]string) (map[string]syscall.Filetime, error) {
	ctimes := make(map[string]syscall.Filetime, len(files))

	for k, fn := range files {
		info, err := os.Stat(filepath.Join(path, fn))
		if err != nil {
			return nil, err
		}

		ct, ok := info.Sys().(*syscall.Win32FileAttributeData)
		if !ok {
			return nil, errors.New("")
		}
		ctimes[k] = ct.CreationTime
	}

	return ctimes, nil
}

func moveFiles(srcpath string, dstpath string, group map[string]string) error {
	for _, v := range group {
		err := retry.CountRetry(
			5,
			func(c int) error {
				return os.Rename(
					filepath.Join(srcpath, v),
					filepath.Join(dstpath, v))
			},
			func(c int) {
				time.Sleep(time.Millisecond * 500)
			})
		if err != nil {
			return err
		}
	}
	return nil
}
