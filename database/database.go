package database

import "database/sql"

const (
	sqlCreateTables string = `
CREATE TABLE config (
	identifier TEXT,
	value TEXT
);

CREATE TABLE clients (
	identifier TEXT,
	ip_address TEXT,
	last_ping  INTEGER
);

CREATE TABLE client_url_list (
    client_id INTEGER,
    url_list_id INTEGER
);
CREATE INDEX client_id_index ON client_url_list (client_id);

CREATE TABLE url_lists (
    id INTEGER PRIMARY KEY,
    url_id INTEGER,
    name TEXT
);

CREATE TABLE urls (
	id INTEGER PRIMARY KEY,
	url TEXT
);

`

	sqlInsertUrl string = "INSERT INTO urls(url) VALUES(?);"
	sqlFetchUrls string = "SELECT url FROM urls;"
	sqlDeleteUrl string = "DELETE FROM urls WHERE url = ?;"

	sqlInsertConfig string = "INSERT INTO config(identifier, value) VALUES(?, ?);"
)

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

func (db *Database) DeleteUrl(url string) (err error) {
	_, err = db.Conn.Exec(sqlDeleteUrl, url)
	return
}

func (db *Database) InsertConfig(identifier string, value string) (err error) {
	_, err = db.Conn.Exec(sqlInsertConfig, identifier, value)
	return
}

func (db *Database) FetchUrls() (urls []string, err error) {
	rows, err := db.Conn.Query(sqlFetchUrls)
	defer rows.Close()
	if err != nil {
		return
	}

	for rows.Next() {
		var url string
		err = rows.Scan(&url)
		urls = append(urls, url)
	}
	err = rows.Err()

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
