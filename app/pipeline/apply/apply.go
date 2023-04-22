package apply

import (
	"errors"
	"obfpl/app/pipeline/apply/temp"
	"obfpl/app/pipeline/sync"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"syscall"

	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"
)

type Apply interface {
	IsAsync() bool

	GetExt(file string) (string, int, error)

	Run(files *StartingFiles, waiter *sync.Waiter)
}

type StartingFiles struct {
	Name  string
	Files map[string]string
}

/* ************************ */

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
	v.group = group
	v.system = system
}

/* ************************ */

type Context struct {
	Name string
	Temp *temp.Temporary
	Exts map[string]string
	Var  *Variable
}

func NewContext(files *StartingFiles, tmp string, exts map[string]string, vari map[string]string) (*Context, error) {
	t, err := temp.NewTemporary(
		tmp,
		files.Name,
		toArrayFromValue(files.Files),
		[]string{"a", "b"})
	if err != nil {
		return nil, err
	}

	v, err := NewVariable(vari)
	if err != nil {
		return nil, err
	}

	return &Context{
		Name: files.Name,
		Temp: t,
		Exts: exts,
		Var:  v,
	}, nil
}

func (c *Context) Cleanup(e func(error)) {
	if err := c.Temp.Cleanup(); err != nil {
		e(err)
	}
}

/* ************************ */

func WithExts(exts map[string]string, getExt func(exts map[string]string, file string) (string, error)) func(string) (string, error) {
	return func(file string) (string, error) {
		e, err := getExt(exts, file)
		return e, err
	}
}

func GetName(sgroup map[string]string, basis string) string {
	if basis == "" {
		//適当な奴を返す
		for _, n := range sgroup {
			return n
		}
		return ""
	}

	n, ok := sgroup[basis]
	if !ok {
		return ""
	}
	return n[:len(n)-len(filepath.Ext(n))]
}

func Match(ptn string, str string, tpath string, enc string) (bool, error) {
	if tpath == "" {
		return regexp.MatchString(ptn, str)
	}

	bin, err := os.ReadFile(tpath)
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

func Call(cmd string) error {
	p := []rune{}
	var delim *rune = nil
	cmds := []string{}

	for _, s := range cmd {
		if s == '\'' || s == '"' {
			if delim == nil {
				delim = &s
				continue
			}
			if *delim == s {
				delim = nil
			}
		}

		if delim == nil && (s == ' ' || (s == '\'' || s == '"')) {
			if len(p) != 0 {
				cmds = append(cmds, string(p))
			}
			p = []rune{}
			continue
		}

		p = append(p, s)
	}
	if len(p) != 0 {
		cmds = append(cmds, string(p))
	}

	inst := exec.Command(cmds[0], cmds[1:]...)
	inst.SysProcAttr = &syscall.SysProcAttr{
		HideWindow: true,
	}

	if _, err := inst.Output(); err != nil {
		eerr, ok := err.(*exec.ExitError)
		if ok {
			return errors.New("error > " + cmds[0] + "\n" + string(eerr.Stderr))
		}
		return err
	}

	return nil
}
