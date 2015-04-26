package database

import "database/sql"

var sqlCreateTables string = `
CREATE TABLE config (
	identifier TEXT,
	value TEXT
);
CREATE TABLE clients (
	identifier TEXT,
	ip_address TEXT,
	last_ping  INTEGER
);
CREATE TABLE urls (
	id INTEGER PRIMARY KEY,
	url TEXT
);
`
var sqlInsertUrl string = "INSERT INTO urls(url) VALUES(?);"
var sqlInsertConfig string = "INSERT INTO config(identifier, value) VALUES(?, ?);"

type Database struct {
	Conn *sql.DB
}

func (db *Database) Close() (err error) {
	err = db.Conn.Close()
	return
}

func (db *Database) InsertUrl(url string) (err error) {
	_, err = db.Conn.Exec(sqlInsertUrl, url)
	return
}

func (db *Database) InsertConfig(identifier string, value string) (err error) {
	_, err = db.Conn.Exec(sqlInsertConfig, identifier, value)
	return
}

func (db *Database) CreateTables() (err error) {
	// Run SQL to create necessary schema
	_, err = db.Conn.Exec(sqlCreateTables)
	return
}

func Connect(database string) (db *Database, err error) {
	c, err := sql.Open("sqlite3", database)
	if err == nil {
		db = &Database{c}
	}

	return
}
