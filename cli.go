package main

import (
	"fmt"
	"log"

	"github.com/codegangsta/cli"
	"github.com/howeyc/gopass"
	"github.com/johnmaguire/wbc/web"
)

func handleInstall(c *cli.Context) {
	var password string

	if resp := confirmDefault("Would you like to setup a password?", true); resp == true {
		fmt.Printf("Password: ")
		password = string(gopass.GetPasswd())

		log.Print("Installing SQLite database")
	}

	log.Print("Admin password: ", password)
}

func handleRun(c *cli.Context) {
	address := fmt.Sprintf("%s:%d", c.String("address"), c.Int("port"))
	web.Start(address)
}

func handleReset(c *cli.Context) {
	log.Print("Wiping database")
}
