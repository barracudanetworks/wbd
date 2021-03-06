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
	app.Name = "wbd"
	app.Usage = "take back control from your televisions"
	app.Version = "0.2.0"
	app.Author = "John Maguire <jmaguire@barracuda.com>"

	app.Commands = []cli.Command{
		{
			Name:    "run",
			Aliases: []string{"r"},
			Usage:   "run the webserver",

			Action: handleRun,

			Flags: []cli.Flag{
				cli.StringFlag{
					Name:   "url,w",
					Value:  "",
					Usage:  "location urls will be relative to (e.g. \"/wbd\")",
					EnvVar: "WBD_URL",
				},
				cli.StringFlag{
					Name:   "listen,l",
					Value:  "0.0.0.0",
					Usage:  "ip address to listen on",
					EnvVar: "WBD_LISTEN",
				},
				cli.IntFlag{
					Name:   "port,p",
					Value:  80,
					Usage:  "port to listen on",
					EnvVar: "WBD_PORT",
				},
				cli.StringFlag{
					Name:   "database,D",
					Value:  "wbd.db",
					Usage:  "sqlite database location",
					EnvVar: "WBD_DATABASE",
				},
			},
		},
		{
			Name:    "url",
			Aliases: []string{"u"},
			Usage:   "add, remove, or list urls in rotation",

			Action: handleUrl,

			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "add,a",
					Usage: "add specified url to rotation",
				},
				cli.StringFlag{
					Name:  "delete,d",
					Usage: "remove specified url from rotation",
				},
				cli.BoolFlag{
					Name:  "list,l",
					Usage: "list urls in rotation (can be combined with --delete or --add)",
				},
				cli.StringFlag{
					Name:   "database,D",
					Value:  "wbd.db",
					Usage:  "sqlite database location",
					EnvVar: "WBD_DATABASE",
				},
			},
		},
		{
			Name:    "list",
			Aliases: []string{"l"},
			Usage:   "add, remove, or list url lists",

			Action: handleList,

			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "add,a",
					Usage: "create a new list",
				},
				cli.StringFlag{
					Name:  "delete,d",
					Usage: "remove an existing list",
				},
				cli.BoolFlag{
					Name:  "list,l",
					Usage: "list url lists in database (can be combined with --delete or --add)",
				},
				cli.StringFlag{
					Name:   "database,D",
					Value:  "wbd.db",
					Usage:  "sqlite database location",
					EnvVar: "WBD_DATABASE",
				},
			},
		},
		{
			Name:    "client",
			Aliases: []string{"c"},
			Usage:   "alias, remove, or list clients",

			Action: handleClient,

			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "alias,a",
					Usage: "alias a client's identifier to something human",
				},
				cli.StringFlag{
					Name:  "to,t",
					Usage: "name to alias client as",
				},
				cli.StringFlag{
					Name:  "delete,d",
					Usage: "remove specified client from database",
				},
				cli.BoolFlag{
					Name:  "list,l",
					Usage: "list known clients",
				},
				cli.StringFlag{
					Name:   "database,D",
					Value:  "wbd.db",
					Usage:  "sqlite database location",
					EnvVar: "WBD_DATABASE",
				},
			},
		},
		{
			Name:    "assign",
			Aliases: []string{"a"},
			Usage:   "assign a client or url to a list",

			Action: handleAssign,

			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "list,l",
					Usage: "target list",
				},
				cli.StringFlag{
					Name:  "url,u",
					Usage: "url to assign to list",
				},
				cli.StringFlag{
					Name:  "client,c",
					Usage: "client to assign to list",
				},
				cli.BoolFlag{
					Name:  "delete,d",
					Usage: "remove association between a list and a client or url",
				},
				cli.StringFlag{
					Name:   "database,D",
					Value:  "wbd.db",
					Usage:  "sqlite database location",
					EnvVar: "WBD_DATABASE",
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
					Name:   "database,D",
					Value:  "wbd.db",
					Usage:  "sqlite database location",
					EnvVar: "WBD_DATABASE",
				},
			},
		},
		{
			Name:  "clean",
			Usage: "delete the database (WARNING: very destructive)",

			Action: handleClean,

			Flags: []cli.Flag{
				cli.StringFlag{
					Name:   "database,D",
					Value:  "wbd.db",
					Usage:  "sqlite database location",
					EnvVar: "WBD_DATABASE",
				},
			},
		},
	}

	app.Run(os.Args)
}
