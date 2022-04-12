package apply

import (
	"obfpl/data"
	"obfpl/libcode/pathlib"
)

func UpdateContext(ctx *Context, proc *data.Process) error {
	/* ---- */
	if proc != nil && proc.Ext != nil {
		ctx.ext = proc.Ext
	}

	/* ---- */
	vars, err := makeVariables(ctx)
	if err != nil {
		return err
	}
	ctx.vars = vars

	/* ---- */
	return nil
}

func outedUpdate(ctx *Context, outlist map[string]string) {
	ctx.vars = outlist
}

func calledUpdate(ctx *Context, newgroup map[string]string) {
	ctx.swap.idx = (ctx.swap.idx + 1) % len(ctx.swap.list)
	ctx.group = newgroup
}

func makeVariables(ctx *Context) (map[string]string, error) {
	vars := make(map[string]string, len(ctx.group)+4)

	name, err := getGroupValue(ctx.group, ctx.profile.Name)
	if err != nil {
		return nil, err
	}

	vars["name"] = pathlib.WithoutExt(name)
	vars["src"] = ctx.swap.list[(ctx.swap.idx)%len(ctx.swap.list)]
	vars["dst"] = ctx.swap.list[(ctx.swap.idx+1)%len(ctx.swap.list)]
	vars["out"] = ctx.outpath

	for k, v := range ctx.group {
		vars[k] = v
	}

	return vars, nil
}
