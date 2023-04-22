package app

import (
	"obfpl/app/obs"
	"obfpl/app/pipeline"
	"os"
	"path/filepath"
)

type Args struct {
	Src     string
	Dst     string
	Profile string
}

type App struct {
	Args   Args
	Logger *Logger
}

func (app App) Run() error {
	//絶対パス変換
	profile, err := filepath.Abs(app.Args.Profile)
	if err != nil {
		return err
	}
	app.Logger.LogSubmsg("プロファイル", profile)

	dst, err := GetValidPath(app.Args.Dst)
	if err != nil {
		return err
	}
	app.Logger.LogSubmsg("出力先", dst)

	src, err := GetValidPath(app.Args.Src)
	if err != nil {
		return err
	}
	app.Logger.LogSubmsg("監視先", src)

	//プロファイル読み込み
	pipe, err := pipeline.NewPipeline(profile, dst)
	if err != nil {
		app.Logger.Log("プロファイルの読み込みに失敗しました")
		return err
	}
	app.Logger.LogSubmsg("プロファイルを読み込みました", "プロファイルの形式 : "+pipe.Type)

	dirobs, err := obs.NewDirObserver(src)
	if err != nil {
		return err
	}

	app.Logger.Log("\nフォルダの監視を開始します...")
	go dirobs.Observe(pipe.CreateChecking(dirobs.Msgch))
	for msg := range dirobs.Msgch {
		app.Logger.Log(msg)
	}

	return nil
}

func GetValidPath(path string) (string, error) {
	r, err := filepath.Abs(path)
	if err != nil {
		return "", err
	}

	if err := os.MkdirAll(r, 0777); err != nil {
		return "", err
	}

	return r, nil
}
