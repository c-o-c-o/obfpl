package variable

import (
	"os"
	"path/filepath"
	"strings"
)

type Variable struct {
	values map[string]string
	group  map[string]string
	system map[string]string
}

func NewVariable(vari map[string]string) (*Variable, error) {
	v := &Variable{
		values: map[string]string{},
		group:  map[string]string{},
		system: map[string]string{},
	}

	path, err := os.Executable()
	if err != nil {
		return nil, err
	}

	edir := filepath.Dir(path)
	for key, val := range vari {
		v.values[key] = strings.ReplaceAll(val, "{@edr}", edir)
	}

	return v, nil
}

func (v *Variable) Apply(str string) string {
	for _, vars := range []map[string]string{v.values, v.group, v.system} {
		for k, v := range vars {
			str = strings.ReplaceAll(str, "{@"+k+"}", v)
		}
	}

	return str
}

func (v *Variable) Update(system map[string]string, group map[string]string) {
	if group != nil {
		v.group = group
	}

	if v.system != nil {
		v.system = system
	}
}
