package main

import (
	"embed"
	"fmt"
	"os"

	"github.com/cxfksword/go-docker-skeleton/pkg/app"
	"github.com/cxfksword/go-docker-skeleton/pkg/log"
	"github.com/cxfksword/go-docker-skeleton/pkg/mode"
	"github.com/cxfksword/go-docker-skeleton/router"
	"github.com/urfave/cli/v2"
)

var (
	// service name
	AppName = "go-docker-skeleton"
	// service description
	AppDesc = "Skeleton for run go service in docker"

	/*********Will auto update by ci build *********/
	Version   = "unknown"
	Commit    = "unknown"
	BuildDate = "unknown"
	Mode      = mode.Dev
	/*********Will auto update by ci build *********/

	//go:embed view/dist/*
	f embed.FS
)

func main() {
	log.Init()
	version := app.VersionInfo{Version: Version, BuildDate: BuildDate, Commit: Commit}
	app := app.New(AppName, AppDesc, Mode, version)

	cmd := &cli.App{
		Name:  AppName,
		Usage: AppDesc,
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:        "debug",
				Aliases:     []string{"vv"},
				Usage:       "Change to debug log level",
				Value:       false,
				Destination: &app.DebugLogLevel,
			},
			&cli.BoolFlag{
				Name:        "verbose",
				Aliases:     []string{"vvv"},
				Usage:       "Change to verbose log level",
				Value:       false,
				Destination: &app.VerboseLogLevel,
			},
		},
		Commands: []*cli.Command{
			{
				Name:  "run",
				Usage: "Run web admin",
				Flags: []cli.Flag{
					&cli.IntFlag{
						Name:        "port",
						Aliases:     []string{"p"},
						Usage:       "web server port",
						Destination: &app.Port,
					},
					&cli.StringFlag{
						Name:        "config",
						Aliases:     []string{"c"},
						Usage:       "Load configuration from `FILE`",
						Destination: &app.ConfigFilePath,
					},
				},
				Action: func(c *cli.Context) error {
					router := router.Create(&f)
					app.Run(router)
					return nil
				},
			},
		},
	}

	fmt.Printf("Starting %s version: %s\n", AppName, Version+"@"+BuildDate+"@"+Mode)
	err := cmd.Run(os.Args)
	if err != nil {
		panic(err)
	}
}
