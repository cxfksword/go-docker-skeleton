package main

import (
	"embed"
	"log"

	"github.com/cxfksword/go-docker-skeleton/app"
	"github.com/cxfksword/go-docker-skeleton/mode"
	"github.com/cxfksword/go-docker-skeleton/model"
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
	f embed.FS //nolint
)

func main() {
	version := model.VersionInfo{Version: Version, BuildDate: BuildDate, Commit: Commit}
	app := app.New(Mode, version)

	cmd := &cli.App{
		Name:  AppName,
		Usage: AppDesc,
		Flags: []cli.Flag{
			&cli.IntFlag{
				Name:        "port",
				Value:       9000,
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
	}

	log.Printf("Starting %s version: %s\n", AppName, Version+"@"+BuildDate+"@"+Mode)
	router := router.Create(&f)
	app.Run(cmd, router)
}
