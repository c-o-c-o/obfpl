package apply

import (
	"errors"
	"obfpl/data"
	"obfpl/libcode/retry"
	"obfpl/libcode/storage"
	"os"
	"path/filepath"
	"syscall"
	"time"
)

func Init(pf *data.Profile, name string, lastdst string, pgroup map[string]string) (*ApplyContext, error) {
	name = storage.GetExistsName(pf.Env.Temp, name)
	swaps, err := storage.MakeDirList(filepath.Join(pf.Env.Temp, name), []string{"a", "b"})
	if err != nil {
		return nil, err
	}

	ctimes, err := getCreateTimes(pgroup)
	if err != nil {
		return nil, err
	}

	err = moveGroup(pgroup, swaps[0])
	if err != nil {
		return nil, err
	}

	group := toNameGroup(pgroup)

	return &ApplyContext{
		temp:    name,
		swaps:   swaps,
		name:    name,
		group:   group,
		ctimes:  ctimes,
		ext:     pf.Ext,
		lastdst: lastdst,
		profile: pf,
	}, nil
}

func getCreateTimes(pgroup map[string]string) (map[string]syscall.Filetime, error) {
	ctimes := make(map[string]syscall.Filetime, len(pgroup))

	for k, v := range pgroup {
		info, err := os.Stat(v)
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

func toNameGroup(pgroup map[string]string) map[string]string {
	group := make(map[string]string, len(pgroup))
	for k, v := range pgroup {
		group[k] = filepath.Base(v)
	}
	return group
}

func moveGroup(pgroup map[string]string, dstp string) error {
	for _, v := range pgroup {
		err := retry.CountRetry(
			5,
			func(c int) error {
				return os.Rename(v, filepath.Join(dstp, filepath.Base(v)))
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
