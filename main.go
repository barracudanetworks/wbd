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

	app.Commands = []cli.Command{
		{
			Name:    "run",
			Aliases: []string{"r"},
			Usage:   "run the webserver",

			Action: handleRun,

			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "address,H",
					Value: "0.0.0.0",
					Usage: "address to listen on",
				},
				cli.IntFlag{
					Name:   "port,p",
					Value:  80,
					Usage:  "port to listen on",
					EnvVar: "WBC_SERVE_PORT",
				},
			},
		},
		{
			Name:    "install",
			Aliases: []string{"i"},
			Usage:   "install the database",

			Action: handleInstall,
		},
		{
			Name:  "reset",
			Usage: "reset the database (WARNING: very destructive)",

			Action: handleReset,
		},
	}

	app.Run(os.Args)
}
