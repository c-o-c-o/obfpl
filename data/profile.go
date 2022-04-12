package data

type Profile struct {
	Env struct {
		Temp     string `yaml:"temp"`
		ExecRule string `yaml:"exec-rule"`
	} `yaml:"env"`
	Ext    map[string]string `yaml:"ext"`
	Name   string            `yaml:"name"`
	Var    map[string]string `yaml:"var"`
	Proc   []Process         `yaml:"process"`
	Notify []string          `yaml:"notify"`
}

type Process struct {
	Ptn    string            `yaml:"ptn"`
	Trg    string            `yaml:"trg"`
	Enc    string            `yaml:"enc"`
	Ext    map[string]string `yaml:"ext"`
	IsWait bool              `yaml:"wait"`
	Cmd    string            `yaml:"cmd"`
}
