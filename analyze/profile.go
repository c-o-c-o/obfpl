package analyze

import (
	"obfpl/data"
	"obfpl/libcode/strfuncs"
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
		pf.Var[k] = strings.ReplaceAll(v, "{@edr}", exedp)
	}

	vari := make(map[string]string, len(pf.Var))
	for k, v := range pf.Var {
		vari["{@"+k+"}"] = v
	}

	for i := range pf.Proc {
		pf.Proc[i].Cmd = strfuncs.ReplaceMap(pf.Proc[i].Cmd, vari)
	}

	for i := range pf.Notify {
		pf.Notify[i] = strfuncs.ReplaceMap(pf.Notify[i], vari)
	}

	return &pf
}
