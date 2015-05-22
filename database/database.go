package database

import "database/sql"

const (
	// schema
	sqlCreateTables string = `
CREATE TABLE config (
	identifier TEXT NOT NULL,
	value TEXT
);

CREATE TABLE clients (
	identifier  TEXT NOT NULL,
	ip_address  TEXT NOT NULL,
	last_ping   TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
	url_list_id INTEGER,
);

CREATE TABLE url_lists (
    id INTEGER PRIMARY KEY,
    name TEXT
);

CREATE TABLE urls (
	id INTEGER PRIMARY KEY,
	url TEXT
);

CREATE TABLE url_list_url (
	id INTEGER PRIMARY KEY,
	url_id INTEGER,
	url_list_id INTEGER
);

`

	// clients table
	sqlInsertClient string = `
	INSERT INTO clients (identifier, ip_address)
	VALUES(?, ?);
	`
	sqlGetClient string = `
	SELECT identifier, timestamp, ip_address, url_list_id
	FROM clients WHERE identifier = ?;
	`
	sqlSetClientList      string = "UPDATE clients SET url_list_id = ? WHERE identifer = ?;"
	sqlSetClientIpAddress string = "UPDATE clients SET ip_address = ? WHERE identifier = ?;"
	sqlTouchClient        string = "UPDATE clients SET last_ping = CURRENT_TIMESTAMP;"

	// urls table
	sqlFindUrlId string = "SELECT id FROM urls WHERE url = ?;"
	sqlInsertUrl string = "INSERT INTO urls(url) VALUES(?);"
	sqlFetchUrls string = "SELECT url FROM urls;"
	sqlDeleteUrl string = "DELETE FROM urls WHERE url = ?;"

	// url_lists table
	sqlFindListId string = "SELECT id FROM url_lists WHERE name = ?;"
	sqlInsertList string = "INSERT INTO url_lists(name) VALUES(?);"
	sqlFetchLists string = "SELECT name FROM url_lists;"
	sqlDeleteList string = "DELETE FROM url_lists WHERE name = ?;"

	// url_lists_url table
	sqlInsertListUrl string = "INSERT INTO url_list_url(url_list_id, url_id) VALUES(?, ?);"
	sqlDeleteListUrl string = "DELETE FROM url_list_url WHERE id = ?;"
	sqlFetchListUrls string = `
	SELECT url FROM urls
	INNER JOIN url_list_url ON url_list_url.url_id = urls.id
	WHERE url_list_id = ?;
	`

	sqlInsertConfig string = "INSERT INTO config(identifier, value) VALUES(?, ?);"
)

type Database struct {
	Conn *sql.DB
}

type Client struct {
	identifier  string
	timestamp   string
	ip_address  string
	url_list_id int
}

func (db *Database) Close() (err error) {
	err = db.Conn.Close()
	return
}

func (db *Database) InsertClient(identifier string, ip_address string) (err error) {
	_, err = db.Conn.Exec(sqlInsertClient, identifier, ip_address)
	return
}

func (db *Database) SetClientIpAddress(identifier string, ip_address string) (err error) {
	_, err = db.Conn.Exec(sqlSetClientIpAddress, identifier, ip_address)
	return
}

func (db *Database) TouchClient(identifier string) (err error) {
	_, err = db.Conn.Exec(sqlTouchClient, identifier)
	return
}

func (db *Database) AssignClientToList(name string, client_id string) (err error) {
	list_id, err := db.FindListId(name)
	if err != nil {
		return
	}

	_, err = db.Conn.Exec(sqlSetClientList, list_id, client_id)
	return
}

func (db *Database) GetClient(identifier string) (client Client, err error) {
	err = db.Conn.QueryRow(sqlGetClient, identifier).Scan(
		&client.identifier,
		&client.timestamp,
		&client.ip_address,
		&client.url_list_id)

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

func (db *Database) InsertList(name string) (err error) {
	_, err = db.Conn.Exec(sqlInsertList, name)
	return
}

func (db *Database) DeleteList(name string) (err error) {
	_, err = db.Conn.Exec(sqlDeleteList, name)
	return
}

func (db *Database) InsertConfig(identifier string, value string) (err error) {
	_, err = db.Conn.Exec(sqlInsertConfig, identifier, value)
	return
}

func (db *Database) FindListId(name string) (id int, err error) {
	row := db.Conn.QueryRow(sqlFindListId, name)

	err = row.Scan(&id)

	return
}

func (db *Database) FindUrlId(url string) (id int, err error) {
	row := db.Conn.QueryRow(sqlFindUrlId, url)

	err = row.Scan(&id)

	return
}

func (db *Database) FetchUrls() (urls []string, err error) {
	rows, err := db.Conn.Query(sqlFetchUrls)
	if err != nil {
		return
	}
	defer rows.Close()

	for rows.Next() {
		var url string

		err = rows.Scan(&url)
		if err != nil {
			return
		}

		urls = append(urls, url)
	}

	err = rows.Err()

	return
}

func (db *Database) FetchLists() (lists []string, err error) {
	rows, err := db.Conn.Query(sqlFetchLists)
	if err != nil {
		return
	}
	defer rows.Close()

	for rows.Next() {
		var name string

		err = rows.Scan(&name)
		if err != nil {
			return
		}

		lists = append(lists, name)
	}

	err = rows.Err()

	return
}

func (db *Database) FetchListUrls(name string) (urls []string, err error) {
	list_id, err := db.FindListId(name)
	if err != nil {
		return
	}

	rows, err := db.Conn.Query(sqlFetchListUrls, list_id)
	if err != nil {
		return
	}
	defer rows.Close()

	for rows.Next() {
		var url string

		err = rows.Scan(&url)
		if err != nil {
			return
		}

		urls = append(urls, url)
	}

	err = rows.Err()

	return
}

func (db *Database) AssignUrlToList(name string, url string) (err error) {
	list_id, err := db.FindListId(name)
	if err != nil {
		return
	}

	url_id, err := db.FindUrlId(url)
	if err != nil {
		return
	}

	_, err = db.Conn.Exec(sqlInsertListUrl, list_id, url_id)
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
