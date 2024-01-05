package util

import (
	"obfpl/-packages/array"
	"os"
	"path/filepath"
	"regexp"

	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"
)

func GetBasisValue(extGroup map[string]string, basis string) string {
	if basis == "" {
		//適当な奴を返す
		for _, v := range extGroup {
			return v
		}
		return ""
	}

	n, ok := extGroup[basis]
	if !ok {
		return ""
	}
	return n[:len(n)-len(filepath.Ext(n))]
}

func Match(ptn string, str string, textPath string, enc string) (bool, error) {
	if textPath == "" {
		return regexp.MatchString(ptn, str)
	}

	bin, err := os.ReadFile(textPath)
	if err != nil {
		return false, err
	}

	str, err = func() (string, error) {
		switch enc {
		case "shift-jis":
			s, _, err := transform.String(japanese.ShiftJIS.NewDecoder(), string(bin))
			return s, err
		case "utf-8":
			return string(bin), nil
		default:
			return string(bin), nil
		}
	}()
	if err != nil {
		return false, err
	}

	return regexp.MatchString(ptn, str)
}

func CreateGetMoveList(selectExt func(string) (string, error)) func([]string, []string) ([]string, error) {
	return func(srcFileNames, dstFileNames []string) ([]string, error) {
		extList, err := array.MapE(dstFileNames, func(val string) (string, error) {
			return selectExt(val)
		})
		if err != nil {
			return nil, err
		}

		return array.FilterE(srcFileNames, func(val string) (bool, error) {
			ext, err := selectExt(val)
			if err != nil {
				return false, err
			}
			return !array.In(extList, ext), nil
		})
	}
}
