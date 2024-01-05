package profproc

import (
	"errors"
	"obfpl/-packages/sync"
	"obfpl/app/pipeline/profproc/lua"
	"obfpl/app/pipeline/profproc/yml"
	"path/filepath"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type ProfProc interface {
	GetType() string
	SelectExt(filePath string) (string, error)
	CreateExtList() []string
	Call(name string, extGroup map[string]string, waiter *sync.Waiter)
}

func NewProfProc(profilePath string, outPath string) (ProfProc, error) {
	profileType := cases.Title(language.English).String(filepath.Ext(profilePath)[1:])

	switch profileType {
	case "Lua":
		return lua.NewLuaProfProc(profilePath, outPath)
	case "Yaml", "Yml":
		return yml.NewYmlProfProc(profilePath, outPath)
	default:
		return nil, errors.New("please specify the correct profile")
	}
}
