package analyze

import (
	"obfpl/data"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

func ReadProfile(path string) (*data.Profile, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	r := data.Profile{}
	err = yaml.Unmarshal(b, &r)
	if err != nil {
		return nil, err
	}

	return &r, nil
}

func ExpandVar(pf data.Profile, exedp string) *data.Profile {
	for k, v := range pf.Var {
		pf.Var[k] = strings.ReplaceAll(v, ToVarName("edr"), exedp)
	}

	for i := range pf.Proc {
		pf.Proc[i].Cmd = VerReflection(pf.Proc[i].Cmd, pf.Var)
	}

	for i := range pf.Notify {
		pf.Notify[i] = VerReflection(pf.Notify[i], pf.Var)
	}

	return &pf
}

func ToVarName(name string) string {
	return "{@" + name + "}"
}

func VerReflection(tgt string, vars map[string]string) string {
	for k, v := range vars {
		tgt = strings.ReplaceAll(tgt, ToVarName(k), v)
	}

	return tgt
}
