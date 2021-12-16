package process

import (
	"errors"
	"obfpl/data"
	"obfpl/libcode/strfuncs"
	"os"
	"os/exec"
	"regexp"
	"syscall"

	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"
)

func Call(name string, proc data.Process, repList map[string]string) (bool, error) {
	match, err := matching(
		name,
		strfuncs.ReplaceMap(proc.Ptn, repList),
		strfuncs.ReplaceMap(proc.Trg, repList),
		proc.Enc)
	if err != nil {
		return false, err
	}
	if !match {
		return false, nil
	}

	return true, CallExec(strfuncs.ReplaceMap(proc.Cmd, repList))
}

func matching(str string, ptn string, trg string, enc string) (bool, error) {
	if trg != "" {
		buf, err := os.ReadFile(trg)
		if err != nil {
			return false, err
		}

		conv, exist := map[string]encoding.Encoding{
			"":          japanese.ShiftJIS,
			"utf-8":     nil,
			"shift-jis": japanese.ShiftJIS,
		}[enc]
		if !exist {
			return false, err
		}

		if conv != nil {
			str, _, err = transform.String(conv.NewDecoder(), string(buf))
			if err != nil {
				return false, err
			}
		} else {
			str = string(buf)
		}
	}

	reg, err := regexp.Compile(ptn)
	if err != nil {
		return false, err
	}

	return reg.MatchString(str), nil
}

func CallExec(cmd string) error {
	args := SplitCommand(cmd)
	if len(args) < 1 {
		return errors.New("the command is specified incorrectly")
	}

	inst := exec.Command(args[0], args[1:]...)
	inst.SysProcAttr = &syscall.SysProcAttr{
		HideWindow: true,
	}
	_, err := inst.Output()
	if err != nil {
		eerr, ok := err.(*exec.ExitError)
		if ok {
			return errors.New("error > " + args[0] + "\n" + string(eerr.Stderr))
		}
		return err
	}
	return nil
}

func SplitCommand(cmd string) []string {
	rcmd := []rune(cmd)
	scmds := make([]string, 0, 50)
	str := make([]rune, 0, len(rcmd))
	instr := false

	for _, v := range rcmd {
		if v == '"' {
			instr = !instr
		}

		if v == ' ' && !instr {
			scmds = append(scmds, string(DeleteDQ(str)))
			str = make([]rune, 0, len(rcmd))
			continue
		}

		str = append(str, v)
	}

	if len(str) != 0 {
		scmds = append(scmds, string(DeleteDQ(str)))
	}

	return scmds
}

func DeleteDQ(s []rune) []rune {
	if len(s) >= 2 && s[0] == '"' && s[len(s)-1] == '"' {
		s = s[1 : len(s)-1]
	}

	return s
}
