package main

import (
	"fmt"
	"log"
	"os"

	"github.com/codegangsta/cli"
	"github.com/johnmaguire/wbc/web"
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
			// TODO: Implement server
			Action: func(c *cli.Context) {
				address := fmt.Sprintf("%s:%d", c.String("address"), c.Int("port"))
				log.Print("Listening on ", address)

				web.Start(address)
			},
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
			// TODO: Implement installation
			Action: func(c *cli.Context) {
				log.Print("Installing SQLite database")
			},
		},
		{
			Name:  "reset",
			Usage: "reset the database (WARNING: very destructive)",
			// TODO: Implement wipe
			Action: func(c *cli.Context) {
				log.Print("Wiping database")
			},
		},
	}

	app.Run(os.Args)
}
