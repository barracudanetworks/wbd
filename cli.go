package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/johnmaguire/wbc/web"

	"github.com/codegangsta/cli"
	"github.com/howeyc/gopass"
	_ "github.com/mattn/go-sqlite3"
)

var sqlCreateTables string = `
CREATE TABLE config (
	identifier text,
	value text
);
CREATE TABLE clients (
	identifier text,
	ip_address text,
	last_ping  integer
);
`
var sqlInsertPassword string = "INSERT INTO config(identifier, value) VALUES(?, ?);"

func createDatabase(database string, password string) {
	log.Printf("Creating database at %s", database)

	db, err := sql.Open("sqlite3", database)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}

	if _, err := db.Exec(sqlCreateTables); err != nil {
		log.Fatal(err)
	}

	if password != "" {
		if _, err = db.Exec(sqlInsertPassword, "password", password); err != nil {
			log.Fatal(err)
		}
	}

	tx.Commit()

	log.Print("Database created")
}

func handleInstall(c *cli.Context) {
	log.Print("Starting installation")

	var (
		database string
		password string
	)

	database = c.String("database")

	// Don't overwrite db if one already exists
	if _, err := os.Stat(database); err == nil {
		log.Fatal("database already exists")
	}

	if resp := confirmDefault("Would you like to set a password?", true); resp == true {
		fmt.Printf("Password: ")
		password = string(gopass.GetPasswd())
	}

	createDatabase(database, password)
}

func handleRun(c *cli.Context) {
	address := fmt.Sprintf("%s:%d", c.String("address"), c.Int("port"))
	web.Start(address)
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
