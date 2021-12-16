package storage

import (
	"io/ioutil"
	"obfpl/libcode/pathlib"
	"obfpl/libcode/strfuncs"
	"os"
	"path/filepath"
	"syscall"
)

func UpdateCreateTime(path string, ctime syscall.Filetime) error {
	fh, err := syscall.Open(path, os.O_RDWR, 0755)
	if err != nil {
		return err
	}
	defer syscall.Close(fh)

	err = syscall.SetFileTime(fh, &ctime, nil, nil)
	if err != nil {
		return err
	}

	return nil
}

func GetExistsName(path string, name string) string {
	suffix := ""
	woname := pathlib.WithoutExt(name)
	ext := filepath.Ext(name)

	for {
		_, err := os.Stat(filepath.Join(path, woname+suffix+ext))
		if os.IsNotExist(err) {
			break
		}
		suffix = strfuncs.RandomString(4)
	}

	return woname + suffix + ext
}

func GetFileNameSuffix(path string, names []string) string {
	duplicate := func(sf string) bool {
		rslt := make([]string, len(names))
		for i, v := range names {
			cname := pathlib.AddSuffix(v, sf)
			_, err := os.Stat(filepath.Join(path, cname))
			if !os.IsNotExist(err) {
				return true
			}
			rslt[i] = cname
		}
		return false
	}

	suffix := ""
	for {
		if !duplicate(suffix) {
			return suffix
		}
		suffix = strfuncs.RandomString(4)
	}
}

func MakeDirList(path string, namel []string) ([]string, error) {
	r := make([]string, len(namel))

	for i, v := range namel {
		r[i] = filepath.Join(path, v)
	}

	for _, v := range r {
		err := os.MkdirAll(v, 0777)
		if err != nil {
			return nil, err
		}
	}

	return r, nil
}

func RemoveDirInner(src string, throw func(error)) {
	fl, err := ioutil.ReadDir(src)
	if err != nil {
		throw(err)
	}

	for _, v := range fl {
		err = os.RemoveAll(filepath.Join(src, v.Name()))
		if err != nil {
			throw(err)
		}
	}
}
