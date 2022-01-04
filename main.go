package main

import (
	"log"
	"math/rand"
	"obfpl/analyze"
	"obfpl/env"
	"obfpl/pipeline"
	"os"
	"path/filepath"
	"time"

	"github.com/urfave/cli/v2"
)

func main() {
	rand.Seed(time.Now().UnixMilli())
	exedp, err := env.GetExecDir()
	if err != nil {
		log.Fatal(err)
		return
	}

	app := &cli.App{
		Name:            "obfPL",
		Usage:           Version,
		Description:     "このアプリを起動すると指定したのフォルダーを監視します。\n監視中特定のファイルを発見すると、プロファイルに従ってアプリを起動して、生成物を出力します。",
		Version:         Version,
		HideHelpCommand: true,
		Flags: []cli.Flag{
			&cli.PathFlag{
				Name:    "dst",
				Aliases: []string{"d"},
				Value:   "dst",
				Usage:   "処理済みのファイル出力先",
			},
			&cli.PathFlag{
				Name:    "src",
				Aliases: []string{"s"},
				Value:   "src",
				Usage:   "監視するフォルダーのパス",
			},
			&cli.PathFlag{
				Name:    "profile",
				Aliases: []string{"p"},
				Value:   filepath.Join(exedp, "profile.yml"),
				Usage:   "プロファイルのパス",
			},
		},
		Action: appfunc(exedp),
	}

	err = app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func appfunc(exedp string) func(c *cli.Context) error {
	return func(c *cli.Context) error {
		//プロファイル読み込み

		profile, err := filepath.Abs(c.Path("profile"))
		if err != nil {
			return err
		}
		pf, err := analyze.ReadProfile(profile)
		if err != nil {
			return err
		}
		println("プロファイルを読み込みました。\n >", profile)

		//変数展開
		pf = analyze.ExpandVar(*pf, exedp)

		//プロファイルからパイプライン作成
		dstp, err := filepath.Abs(c.Path("dst"))
		if err != nil {
			return err
		}
		pl, err := pipeline.FromPipeline(dstp, pf)
		if err != nil {
			return err
		}
		println("出力先\n >", dstp)

		erch := make(chan error)
		srcp, err := filepath.Abs(c.Path("src"))
		if err != nil {
			return err
		}

		if err := os.MkdirAll(srcp, 0777); err != nil {
			return err
		}
		go pl.ObsFolder(srcp, erch)
		println("監視先\n >", srcp)

		pl.Loging = func(msg string) {
			println(msg)
		}

		println("\nフォルダの監視を開始します...")
		for err := range erch {
			println(err.Error())
		}

		return nil
	}
}
