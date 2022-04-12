package process

import (
	"errors"
	"os/exec"
	"syscall"
)

func Call(cmd string) error {
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

	{
		instr := false

		for _, v := range rcmd {
			if v == '"' || v == '\'' {
				instr = !instr
			}

			if v == ' ' && !instr {
				scmds = append(scmds, string(DeleteQuotation(str)))
				str = make([]rune, 0, len(rcmd))
				continue
			}

			str = append(str, v)
		}
	}

	if len(str) != 0 {
		scmds = append(scmds, string(DeleteQuotation(str)))
	}

	return scmds
}

func DeleteQuotation(s []rune) []rune {
	if len(s) < 2 {
		return s
	}

	issq := (s[0] == '\'' && s[len(s)-1] == '\'')
	isdq := (s[0] == '"' && s[len(s)-1] == '"')

	if issq || isdq {
		s = s[1 : len(s)-1]
	}

	return s
}
