package pipeline

import (
	"obfpl/data"
	"obfpl/libcode/pathlib"
	"obfpl/pipeline/apply"
	"obfpl/pipeline/observe"
	"path/filepath"
)

type Pipeline struct {
	OutputPath string
	Loging     func(string)
	detectList map[string]map[string]string
	profile    data.Profile
}

func FromPipeline(op string, pf *data.Profile) (*Pipeline, error) {
	r := Pipeline{
		OutputPath: op,
		detectList: map[string]map[string]string{},
		profile:    *pf,
	}

	return &r, nil
}

func (pl *Pipeline) ObsFolder(obsp string, erch chan error) {
	apply := func(path string, erch chan error) {
		if pl.profile.Env.ExecRule == "async" {
			go pl.apply(path, erch)
		} else {
			pl.apply(path, erch)
		}
	}

	observe.Folder(obsp, apply, erch)
}

func (pl *Pipeline) apply(path string, erch chan error) {
	name, pgroup, isapply := pl.checkDetectList(path)
	if !isapply {
		return
	}

	ctx, err := apply.Init(&pl.profile, name, pl.OutputPath, pgroup)
	if err != nil {
		erch <- err
		return
	}
	ctx.Loging = pl.Loging

	isbreak := apply.Run(ctx, erch)
	if isbreak {
		return
	}

	apply.End(ctx, erch)
}

func (pl *Pipeline) checkDetectList(path string) (string, map[string]string, bool) {
	etype, err := apply.MatchExt(pl.profile.Ext, filepath.Ext(path)[1:])
	if err != nil {
		return "", nil, false
	}
	name := pathlib.WithoutExt(path)
	_, exist := pl.detectList[name]
	if !exist {
		pl.detectList[name] = make(map[string]string)
	}
	pl.detectList[name][etype] = path

	if len(pl.detectList[name]) == len(pl.profile.Ext) {
		pgroup := pl.detectList[name]
		delete(pl.detectList, name)

		return name, pgroup, true
	}

	return name, nil, false
}
