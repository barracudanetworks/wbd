package main

import (
	"fmt"
	"log"
	"os"

	"github.com/barracudanetworks/wbd/config"
	"github.com/barracudanetworks/wbd/database"
	"github.com/barracudanetworks/wbd/web"

	"github.com/codegangsta/cli"
	"github.com/howeyc/gopass"
)

func handleRun(c *cli.Context) error {
	conf := &config.Configuration{
		ListenAddress: c.String("listen"),
		ListenPort:    c.Int("port"),
		WebAddress:    c.String("url"),
		Database:      c.String("database"),
	}

	if _, err := os.Stat(conf.Database); err != nil {
		log.Fatal("database does not exist")
	}
	log.Printf("Using database %s", conf.Database)

	if conf.ListenPort == 0 {
		conf.ListenPort = 80
	}
	if conf.ListenAddress == "" {
		conf.ListenAddress = "0.0.0.0"
	}

	web.Start(conf)

	return nil
}

func handleUrl(c *cli.Context) error {
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

	return nil
}

func handleClient(c *cli.Context) error {
	if _, err := os.Stat(c.String("database")); err != nil {
		log.Fatal("database does not exist")
	}
	log.Printf("Using database %s", c.String("database"))

	aliasClient, toAlias := c.String("alias"), c.String("to")
	deleteClient := c.String("delete")

	if aliasClient != "" && deleteClient != "" {
		log.Fatal("Can't both remove and alias a client")
	}

	db, err := database.Connect(c.String("database"))
	defer db.Close()
	if err != nil {
		log.Fatal(err)
	}

	if aliasClient != "" {
		if toAlias == "" {
			log.Fatal("No alias specified (use --to)")
		}

		log.Printf("Aliasing client '%s' to '%s'", aliasClient, toAlias)
		if err := db.SetClientAlias(aliasClient, toAlias); err != nil {
			log.Fatal(err)
		}
	}

	if deleteClient != "" {
		log.Printf("Removing client '%s' from the database", deleteClient)
		if err := db.DeleteClient(deleteClient); err != nil {
			log.Fatal(err)
		}
	}

	if c.Bool("list") {
		log.Print("Known clients:")
		clients, err := db.FetchClients()
		if err != nil {
			log.Fatal(err)
		}

		for _, client := range clients {
			if client.Alias == "" {
				log.Printf("  %s (%s) - Last active %s", client.Identifier, client.IpAddress, client.LastPing)
			} else {
				log.Printf("  %s [%s] (%s) - Last active %s", client.Alias, client.Identifier, client.IpAddress, client.LastPing)
			}
		}
	}

	return nil
}

func handleAssign(c *cli.Context) error {
	if _, err := os.Stat(c.String("database")); err != nil {
		log.Fatal("database does not exist")
	}
	log.Printf("Using database %s", c.String("database"))

	deleteFlag := c.Bool("delete")
	assignList := c.String("list")
	assignUrl, assignClient := c.String("url"), c.String("client")

	if (assignList == "" && (!deleteFlag || assignClient == "")) || (assignClient == "" && assignUrl == "") {
		log.Fatal("Must specify a list, and a client or URL to assign to it")
	}

	db, err := database.Connect(c.String("database"))
	defer db.Close()
	if err != nil {
		log.Fatal(err)
	}

	if assignUrl != "" {
		// delete association if delete flag is true
		if deleteFlag {
			if err := db.RemoveUrlFromList(assignList, assignUrl); err != nil {
				log.Fatal(err)
			}
			log.Printf("Removed URL %s from list %s", assignUrl, assignList)
		} else {
			if err := db.AssignUrlToList(assignList, assignUrl); err != nil {
				log.Fatal(err)
			}
			log.Printf("Assigned URL %s to list %s", assignUrl, assignList)
		}
	}
	if assignClient != "" {
		// delete association if delete flag is true
		if deleteFlag {
			if err := db.RemoveClientFromList(assignClient); err != nil {
				log.Fatal(err)
			}
			log.Printf("Assigned client %s back to the Default list", assignClient)
		} else {
			if err := db.AssignClientToList(assignList, assignClient); err != nil {
				log.Fatal(err)
			}
			log.Printf("Assigned client %s to list %s", assignClient, assignList)
		}
	}

	return nil
}

func handleList(c *cli.Context) error {
	if _, err := os.Stat(c.String("database")); err != nil {
		log.Fatal("database does not exist")
	}
	log.Printf("Using database %s", c.String("database"))

	addList, deleteList := c.String("add"), c.String("delete")
	if addList != "" && deleteList != "" {
		log.Fatal("Can't both remove and add a list")
	}

	db, err := database.Connect(c.String("database"))
	defer db.Close()
	if err != nil {
		log.Fatal(err)
	}

	if addList != "" {
		log.Printf("Creating list %s", addList)
		if err := db.InsertList(addList); err != nil {
			log.Fatal(err)
		}
	}

	if deleteList != "" {
		log.Printf("Deleting list %s", deleteList)
		if err := db.DeleteList(deleteList); err != nil {
			log.Fatal(err)
		}
	}

	if c.Bool("list") {
		log.Print("URL lists defined:")
		lists, err := db.FetchLists()
		if err != nil {
			log.Fatal(err)
		}

		for _, list := range lists {
			urls, err := db.FetchListUrlsByName(list)
			if err != nil {
				log.Fatal(err)
			}

			log.Print("  ", list)

			for _, url := range urls {
				log.Print("    ", url)
			}
		}
	}

	return nil
}

func handleInstall(c *cli.Context) error {
	log.Print("Starting installation")

	var (
		path     string
		password []byte
	)

	path = c.String("database")

	// Don't overwrite db if one already exists
	if _, err := os.Stat(path); err == nil {
		log.Fatal("database already exists")
	}

	if resp := confirmDefault("Would you like to set a password?", true); resp == true {
		fmt.Printf("Password: ")

		var err error
		password, err = gopass.GetPasswd()
		if err != nil {
			log.Fatal(err)
		}
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
		tx.Rollback()
		log.Fatal(err)
	}

	// Insert password if one was given
	if len(password) == 0 {
		if err = db.InsertConfig("password", string(password)); err != nil {
			tx.Rollback()
			log.Fatal(err)
		}
	}

	tx.Commit()
	log.Print("Database created")

	return nil
}

func handleClean(c *cli.Context) error {
	database := c.String("database")
	log.Printf("Removing database at %s", database)

	if _, err := os.Stat(database); err != nil {
		log.Fatal("database does not exist")
	}

	if err := os.Remove(database); err != nil {
		log.Fatal(err)
	}

	log.Print("Database removed")

	return nil
}
