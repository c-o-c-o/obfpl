package apply

import (
	"obfpl/analyze"
	"obfpl/data"
	"obfpl/libcode/encode"
	"os"
	"regexp"
)

func Match(ctx *Context, proc data.Process) (bool, error) {
	str, err := getMatchString(
		ctx,
		analyze.VerReflection(proc.Trg, ctx.vars),
		proc.Enc)
	if err != nil {
		return false, err
	}

	reg, err := regexp.Compile(analyze.VerReflection(proc.Ptn, ctx.vars))
	if err != nil {
		return false, err
	}

	return reg.MatchString(str), nil
}

func getMatchString(ctx *Context, trg string, enc string) (string, error) {
	if trg == "" {
		return getGroupValue(ctx.group, ctx.profile.Name)
	}

	return readText(trg, enc)
}

func readText(path string, enc string) (string, error) {
	tbuf, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}

	return encode.Convert(string(tbuf), enc)
}
