package extgroup

import "path/filepath"

/*
ファイルパスから拡張子なしのファイル名を取得します
*/
func GetName(filePath string) string {
	e := filepath.Ext(filePath)
	n := filepath.Base(filePath)
	return n[:len(n)-len(e)]
}
