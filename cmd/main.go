package main

import (
	"fmt"
	"github.com/MrBoombastic/GhostBackupper/pkg/backup"
	"github.com/urfave/cli/v2"
	"log"
	"os"
)

const Version = "1.0.0"

func main() {
	app := &cli.App{
		Name:      "ghostbackupper",
		Usage:     "Simple GhostCMS backup tool",
		UsageText: "ghostbackupper <command> <options...> - run 'ghostbackupper <command> --help' for more",
		Commands: []*cli.Command{
			{
				Name:        "version",
				Description: "Shows GhostBackupper version",
				Aliases:     []string{"v", "ver"},
				Action: func(context *cli.Context) error {
					fmt.Printf("GhostBackupper version %v. Choo-choo!", Version)
					return nil
				},
			},
			{
				Name:        "backup",
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
						Name:     "db_pass",
						Usage:    "Your MySQL server password",
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
					&cli.StringFlag{
						Name:     "db_database",
						Usage:    "Your MySQL database name",
						Required: true,
					},
					&cli.UintFlag{
						Name:  "db_port",
						Value: 3306,
						Usage: "Your MySQL server port",
					},
					&cli.StringFlag{
						Name:     "content",
						Usage:    "Ghost's 'content' directory path",
						Required: true,
					},
					&cli.StringFlag{
						Name:    "output",
						Aliases: []string{"file", "f"},
						Usage:   "Output filename (not path!)",
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
