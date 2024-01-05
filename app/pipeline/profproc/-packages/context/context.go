package context

import (
	"obfpl/app/pipeline/profproc/-packages/context/temp"
	"obfpl/app/pipeline/profproc/-packages/context/variable"
)

type Context struct {
	Name string
	Temp *temp.Temporary
	Exts map[string]string
	Vari *variable.Variable
}

func NewContext(name string, fileNames []string, tempPath string, exts map[string]string, vari map[string]string) (*Context, error) {
	t, err := temp.NewTemporary(
		tempPath,
		name,
		fileNames,
		[]string{"a", "b"})
	if err != nil {
		return nil, err
	}

	v, err := variable.NewVariable(vari)
	if err != nil {
		return nil, err
	}

	return &Context{
		Name: name,
		Temp: t,
		Exts: exts,
		Vari: v,
	}, nil
}
