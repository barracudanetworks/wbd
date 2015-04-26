package main

import (
	"fmt"
	"log"
	"os"

	"github.com/johnmaguire/wbc/database"
	"github.com/johnmaguire/wbc/web"

	"github.com/codegangsta/cli"
	"github.com/howeyc/gopass"
	_ "github.com/mattn/go-sqlite3"
)

func handleRun(c *cli.Context) {
	address := fmt.Sprintf("%s:%d", c.String("address"), c.Int("port"))
	web.Start(address)
}

func handleInstall(c *cli.Context) {
	log.Print("Starting installation")

	var (
		path     string
		password string
	)

	path = c.String("database")

	// Don't overwrite db if one already exists
	if _, err := os.Stat(path); err == nil {
		log.Fatal("database already exists")
	}

	if resp := confirmDefault("Would you like to set a password?", true); resp == true {
		fmt.Printf("Password: ")
		password = string(gopass.GetPasswd())
	}

	log.Printf("Creating database at %s", path)

	// Create a new connection to the database
	db, err := database.Connect(path)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Start a transaction
	tx, err := db.Conn.Begin()
	if err != nil {
		log.Fatal(err)
	}

	// Create table schema
	if err = db.CreateTables(); err != nil {
		log.Fatal(err)
	}

	// Insert password if one was given
	if password != "" {
		if err = db.InsertConfig("password", password); err != nil {
			log.Fatal(err)
		}
	}

	// Insert a default URL
	if err = db.InsertUrl("https://google.com/"); err != nil {
		log.Fatal(err)
	}

	tx.Commit()
	log.Print("Database created")
}

func handleClean(c *cli.Context) {
	database := c.String("database")
	log.Printf("Removing database at %s", database)

	if _, err := os.Stat(database); err != nil {
		log.Fatal("database does not exist")
	}

	if err := os.Remove(database); err != nil {
		log.Fatal(err)
	}

	log.Print("Database removed")
}
