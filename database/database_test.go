package database

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestConnection(t *testing.T) {
	// temporary in-memory database
	db, err := Connect(":memory:")
	assert.Nil(t, err, "It should connect and create the DB.")
	assert.IsType(t, new(Database), db, "It should return a Database instance.")

	err = db.Close()
	assert.Nil(t, err, "It should close the database.")
}

func TestTableCreation(t *testing.T) {
	db, _ := Connect(":memory:")
	defer db.Close()

	err := db.CreateTables()
	assert.Nil(t, err, "It should create the tables.")
}

func TestClient(t *testing.T) {
	assert := assert.New(t)

	db, _ := Connect(":memory:")
	defer db.Close()

	db.CreateTables()

	// test client creation
	err := db.InsertClient("TestClient", "0.0.0.0")
	assert.Nil(err)

	// test retrieving a client and that creation was successful
	client, err := db.GetClient("TestClient")
	assert.Nil(err)
	assert.IsType(*new(Client), client, "It should return a Client instance.")

	assert.Equal("0.0.0.0", client.IpAddress)
	assert.Equal("TestClient", client.Identifier)
	assert.NotEqual("", client.LastPing, "It should be an ISO-8601 formatted date")
	assert.Equal(0, client.UrlListId)

	// update the last ping
	time.Sleep(1000 * time.Millisecond)
	db.TouchClient(client.Identifier)

	// pull updated info
	updated_client, _ := db.GetClient(client.Identifier)

	// verify the ping updated
	assert.NotEqual(client.LastPing, updated_client.LastPing)

	// verify the other data wasn't tampered
	assert.Equal(client.Identifier, updated_client.Identifier)
	assert.Equal(client.IpAddress, updated_client.IpAddress)
	assert.Equal(client.UrlListId, updated_client.UrlListId)
}
