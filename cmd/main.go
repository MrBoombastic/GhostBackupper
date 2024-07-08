package main

import (
	"fmt"
	"github.com/MrBoombastic/GhostBackupper/pkg/backup"
	"github.com/urfave/cli/v2"
	"log"
	"os"
)

const Version = "0.0.6"

func main() {
	cli.VersionPrinter = func(cCtx *cli.Context) {
		fmt.Printf("GhostBackupper version %s\n", cCtx.App.Version)
	}
	cli.VersionFlag = &cli.BoolFlag{
		Name:    "version",
		Aliases: []string{"v"},
		Usage:   "Shows GhostBackupper version",
	}
	app := &cli.App{
		Name:      "ghostbackupper",
		Usage:     "Simple GhostCMS backup tool",
		UsageText: "ghostbackupper <command> <options...> - run 'ghostbackupper <command> --help' for more",
		Version:   Version,
		Commands: []*cli.Command{
			{
				Name:        "backup",
				Usage:       "Backs up whole Ghost - database and files",
				Description: "Backs up whole Ghost - database and files",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:  "db_host",
						Usage: "Your MySQL server address",
						Value: "localhost",
					},
					&cli.StringFlag{
						Name:     "db_user",
						Usage:    "Your MySQL username",
						Required: true,
					},
					&cli.StringFlag{
						Name:     "db_password",
						Usage:    "Your MySQL server password",
						Required: true,
					},
					&cli.StringFlag{
						Name:     "db_database",
						Usage:    "Your MySQL database name",
						Required: true,
					},
					&cli.StringFlag{
						Name:     "mega_login",
						Usage:    "Your Mega.nz login (only if you want to upload!)",
						Required: false,
					},
					&cli.StringFlag{
						Name:     "mega_password",
						Usage:    "Your Mega.nz password (only if you want to upload!)",
						Required: false,
					},
					&cli.UintFlag{
						Name:  "db_port",
						Value: 3306,
						Usage: "Your MySQL server port",
					},
					&cli.StringFlag{
						Name:      "content",
						Usage:     "Ghost's 'content' directory path",
						TakesFile: true,
						Required:  true,
					},
					&cli.StringFlag{
						Name:    "output",
						Aliases: []string{"file", "f"},
						Usage:   "Output filename (not path!). Unix epoch will be appended to it at the beginning.",
						Value:   "backup.tar.gz",
					},
				},
				Action: func(context *cli.Context) error {
					return backup.Create(context)
				},
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
