package exec

import (
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"
)

func GetExecDir() (string, error) {
	p, err := os.Executable()
	if err != nil {
		return "", err
	}
	return filepath.Dir(p), nil
}

func Call(cmd string) error {
	// クォーテーション囲いでの文字列化処理を書いてるので泥臭くなっている
	// 直すなら頑張って

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
