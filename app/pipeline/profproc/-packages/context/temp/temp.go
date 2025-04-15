package temp

import (
	"io/fs"
	"obfpl/-packages/array"
	"obfpl/app/pipeline/profproc/-packages/context/temp/suffix"
	"os"
	"path/filepath"
	"strings"
	"syscall"
	"time"
)

type Temporary struct {
	baseDir  string
	swapList []string
	swapIdx  int
	cTime    time.Time
}

/*
２つ以上のディレクトリで「入力元」「出力先」を示す２つの一時ディレクトリを提供する

* base テンポラリファイルを生成するディレクトリ

* name 任意のフォルダ名、被るとサフィックスを付けて生成する

* fileNames 管理するファイル名のリスト

* swaps 行き来するディレクトリの名前
*/
func NewTemporary(base string, name string, fileNames []string, swaps []string) (*Temporary, error) {
	if len(swaps) < 2 {
		panic("code error : 2 or more 'swaps' are always required")
	}

	baseDir, err := getValidDirPath(base, strings.TrimSpace(name))
	if err != nil {
		return nil, err
	}

	for _, swap := range swaps {
		if err := os.MkdirAll(filepath.Join(baseDir, swap), 0777); err != nil {
			return nil, err
		}
	}

	for _, fileName := range fileNames {
		//ファイルを一時フォルダへ移動を最大5回リトライ
		err := retry(
			5,
			func(c int) error {
				return os.Rename(fileName, filepath.Join(baseDir, swaps[0], filepath.Base(fileName)))
			},
			func(c int) {
				time.Sleep(time.Millisecond * 500)
			})
		if err != nil {
			return nil, err
		}
	}

	return &Temporary{
		baseDir:  baseDir,
		swapList: swaps,
		swapIdx:  0,
		cTime:    time.Now(),
	}, nil
}

/*
現在の一時フォルダーパスを取得する

	return src, dst
*/
func (t *Temporary) GetPaths() (string, string) {
	src := t.swapList[t.swapIdx%len(t.swapList)]
	dst := t.swapList[(t.swapIdx+1)%len(t.swapList)]

	return filepath.Join(t.baseDir, src), filepath.Join(t.baseDir, dst)
}

/*
一時フォルダー内のファイルを取得する

target 'src' | 'dst'

	return []FileName, error
*/
func (t *Temporary) LoadFileNames(target string) ([]string, error) {
	i, ok := map[string]int{
		"src": 0,
		"dst": 1,
	}[target]
	if !ok {
		panic("'target' must be either 'src' or 'dst'")
	}

	tgt := t.swapList[(t.swapIdx+i)%len(t.swapList)]

	files, err := os.ReadDir(filepath.Join(t.baseDir, tgt))
	if err != nil {
		return nil, err
	}

	return array.Map(files, func(file fs.DirEntry) string {
		return file.Name()
	}), nil
}

/*
一時フォルダーを切り替える、その際に一定のルールに従ってファイルを移動させる。移動しなかったファイルは削除する
*/
func (t *Temporary) Exchange(getMoveList func(srcFileNames []string, dstFileNames []string) ([]string, error)) error {

	srcFileNames, err := t.LoadFileNames("src")
	if err != nil {
		return err
	}

	dstFileNames, err := t.LoadFileNames("dst")
	if err != nil {
		return err
	}

	moveList, err := getMoveList(srcFileNames, dstFileNames)
	if err != nil {
		return err
	}

	src, dst := t.GetPaths()
	for _, moveFileName := range moveList {
		err := os.Rename(filepath.Join(src, moveFileName), filepath.Join(dst, moveFileName))
		if err != nil {
			return err
		}
	}

	//srcディレクトリ内 全削除
	files, err := os.ReadDir(src)
	if err != nil {
		return err
	}

	for _, file := range files {
		err := os.RemoveAll(filepath.Join(src, file.Name()))
		if err != nil {
			return err
		}
	}

	t.swapIdx += 1
	return nil
}

func (t *Temporary) Output(outdir string, files []string) error {
	srcdir, _ := t.GetPaths()
	sf, err := getValidFilesSuffix(outdir, files)
	if err != nil {
		return err
	}

	for _, file := range files {
		src := filepath.Join(srcdir, file)
		out := filepath.Join(outdir, suffix.With(file, sf))
		if err != nil {
			return err
		}
		if err := os.Rename(src, out); err != nil {
			return err
		}
	}

	//生成時刻を処理開始時に設定
	ctime := syscall.NsecToFiletime(t.cTime.UnixNano())
	for _, file := range files {
		fh, err := syscall.Open(filepath.Join(outdir, file), os.O_RDWR, 0755)
		if err != nil {
			return err
		}
		defer syscall.Close(fh)

		if err := syscall.SetFileTime(
			fh,
			&ctime,
			nil,
			nil,
		); err != nil {
			return err
		}
	}

	return nil
}

func (t *Temporary) Cleanup() error {
	return os.RemoveAll(t.baseDir)
}

func getValidFilesSuffix(base string, files []string) (string, error) {
	var result = ""
	err := retry(
		10,
		func(c int) error {
			var sf = ""
			if c != 0 {
				sf = suffix.Get(8)
			}

			if err := suffix.TryWithFiles(base, files, sf); err != nil {
				return err
			}
			result = sf
			return nil
		},
		func(c int) {
			time.Sleep(time.Millisecond * 100)
		})
	if err != nil {
		return "", err
	}
	return result, nil
}

func getValidDirPath(base string, name string) (string, error) {
	var result = ""
	err := retry(
		10,
		func(c int) error {
			cand := filepath.Join(base, suffix.With(name, suffix.Get(16)))
			if err := os.MkdirAll(cand, 0777); err != nil {
				return err
			}
			result = cand
			return nil
		},
		func(c int) {
			time.Sleep(time.Millisecond * 100)
		})
	if err != nil {
		return "", err
	}
	return result, nil
}

func retry(c int, f func(c int) error, d func(c int)) error {
	for i := 0; i < c; i++ {
		err := f(i)

		if err == nil {
			break
		}

		if i+1 >= c {
			return err
		}
		d(i)
	}
	return nil
}
