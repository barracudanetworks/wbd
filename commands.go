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
	if _, err := os.Stat(c.String("database")); err != nil {
		log.Fatal("database does not exist")
	}
	log.Printf("Using database %s", c.String("database"))

	log.Print("Attempting to start web server")

	address := fmt.Sprintf("%s:%d", c.String("listen"), c.Int("port"))
	web.Start(address, c.String("database"))
}

func handleUrl(c *cli.Context) {
	if _, err := os.Stat(c.String("database")); err != nil {
		log.Fatal("database does not exist")
	}
	log.Printf("Using database %s", c.String("database"))

	addUrl, deleteUrl := c.String("add"), c.String("delete")
	if addUrl != "" && deleteUrl != "" {
		log.Fatal("Can't both remove and add a URL")
	}

	db, err := database.Connect(c.String("database"))
	defer db.Close()
	if err != nil {
		log.Fatal(err)
	}

	if addUrl != "" {
		log.Printf("Adding url %s to rotation", addUrl)
		if err := db.InsertUrl(addUrl); err != nil {
			log.Fatal(err)
		}
	}

	if deleteUrl != "" {
		log.Printf("Removing url %s from rotation", deleteUrl)
		if err := db.DeleteUrl(deleteUrl); err != nil {
			log.Fatal(err)
		}
	}

	if c.Bool("list") {
		log.Print("URLs in rotation:")
		urls, err := db.FetchUrls()
		if err != nil {
			log.Fatal(err)
		}

		for _, url := range urls {
			log.Print("  ", url)
		}
	}
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
	defer db.Close()
	if err != nil {
		log.Fatal(err)
	}

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
