package main

import (
	"log"
	"obfpl/-packages/exec"
	llog "obfpl/-packages/log"
	"obfpl/app"
	"os"
	"path/filepath"

	"github.com/urfave/cli/v2"
)

func main() {
	edir, err := exec.GetExecDir()
	if err != nil {
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
				Value:   filepath.Join(edir, "profile.yml"),
				Usage:   "プロファイルのパス",
			},
		},
		Action: Action,
	}

	err = app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func Action(c *cli.Context) error {
	a := app.App{
		Args: app.Args{
			Src:     c.Path("src"),
			Dst:     c.Path("dst"),
			Profile: c.Path("profile"),
		},
		Logger: llog.NewLogger(func(msg string) {
			println(msg)
		}),
	}

	return a.Run()
}
