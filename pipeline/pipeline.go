package pipeline

import (
	"obfpl/data"
	"obfpl/pipeline/apply"
	"obfpl/pipeline/observe"
	"obfpl/pipeline/sync"
	"path/filepath"
)

type Pipeline struct {
	OutPath string
	Loging  func(string)
	profile data.Profile
}

func FromProfile(op string, pf *data.Profile) (*Pipeline, error) {
	r := Pipeline{
		OutPath: op,
		profile: *pf,
	}

	return &r, nil
}

func (pl *Pipeline) ObsFolder(obsp string, erch chan error) {
	waiter := sync.NewClosedWaiter(erch)
	pool := NewExtPool(pl.profile.Ext)

	observe.Folder(
		obsp,
		func(path string, erch chan error) {
			//ファイルチェック
			group, err := pool.Add(filepath.Base(path))
			if err != nil {
				erch <- err
				return
			}
			if group == nil {
				return
			}
			/* -------- */

			ctx, err := apply.CreateContext(&pl.profile, obsp, pl.OutPath, group)
			if err != nil {
				erch <- err
				return
			}
			waiter = waiter.Next()

			switch pl.profile.Env.ExecRule {
			case "async":
				go pl.apply(ctx, waiter)
			case "wait":
				pl.apply(ctx, waiter)
			}
		},
		erch)
}

func (pl *Pipeline) apply(ctx *apply.Context, waiter *sync.Waiter) {
	defer waiter.Destroy()
	ctx.Loging = pl.Loging

	for _, p := range pl.profile.Proc {
		err := apply.UpdateContext(ctx, &p)
		if err != nil {
			waiter.Error(err)
			return
		}

		//処理待ち
		if p.IsWait {
			err := waiter.Wait()
			if err != nil {
				waiter.Error(err)
				return
			}
		}

		//呼び出しチェック
		m, err := apply.Match(ctx, p)
		if err != nil {
			waiter.Error(err)
			return
		}

		//アプリ呼び出し
		if m {
			err := apply.Call(ctx, p)
			if err != nil {
				waiter.Error(err)
				return
			}
		}
	}

	//出力
	err := apply.Output(ctx)
	if err != nil {
		waiter.Error(err)
		return
	}

	//通知用アプリ呼び出し
	for _, cmd := range pl.profile.Notify {
		err := apply.Notify(ctx, cmd)
		if err != nil {
			waiter.Error(err)
			return
		}
	}

	waiter.Close()
}
