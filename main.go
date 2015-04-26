package main

import (
	"log"
	"os"

	"github.com/codegangsta/cli"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)
	log.SetPrefix("(root) ")

	app := cli.NewApp()
	app.Name = "wbc"
	app.Usage = "take back control from your televisions"
	app.Version = "0.1.0"
	app.Author = "John Maguire <contact@johnmaguire.me>"

	app.Commands = []cli.Command{
		{
			Name:    "run",
			Aliases: []string{"r"},
			Usage:   "run the webserver",

			Action: handleRun,

			Flags: []cli.Flag{
				cli.StringFlag{
					Name:   "address,H",
					Value:  "0.0.0.0",
					Usage:  "address to listen on",
					EnvVar: "WBC_ADDRESS",
				},
				cli.IntFlag{
					Name:   "port,p",
					Value:  80,
					Usage:  "port to listen on",
					EnvVar: "WBC_PORT",
				},
				cli.StringFlag{
					Name:   "database,d",
					Value:  "wbc.db",
					Usage:  "sqlite database location",
					EnvVar: "WBC_DATABASE",
				},
			},
		},
		{
			Name:    "install",
			Aliases: []string{"i"},
			Usage:   "install the database",

			Action: handleInstall,

			Flags: []cli.Flag{
				cli.StringFlag{
					Name:   "database,d",
					Value:  "wbc.db",
					Usage:  "sqlite database location",
					EnvVar: "WBC_DATABASE",
				},
			},
		},
		{
			Name:  "clean",
			Usage: "delete the database (WARNING: very destructive)",

			Action: handleClean,

			Flags: []cli.Flag{
				cli.StringFlag{
					Name:   "database,d",
					Value:  "wbc.db",
					Usage:  "sqlite database location",
					EnvVar: "WBC_DATABASE",
				},
			},
		},
	}

	app.Run(os.Args)
}
