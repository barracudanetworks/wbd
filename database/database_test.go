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

func TestLists(t *testing.T) {
	assert := assert.New(t)

	db, _ := Connect(":memory:")
	defer db.Close()

	db.CreateTables()

	list_id, err := db.FindListId("Default")
	assert.Nil(err)
	assert.Equal(DefaultList, list_id)

	list_id, err = db.FindListId("NONEXISTENT")
	assert.NotNil(err)

	err = db.InsertList("test")
	assert.Nil(err)

	list_id, err = db.FindListId("test")
	assert.Nil(err)
	assert.Equal(1, list_id)

	err = db.InsertUrl("http://barracudanetworks.com/")
	assert.Nil(err)

	err = db.AssignUrlToList("test", "http://barracudanetworks.com/")
	assert.Nil(err)

	urls, err := db.FetchListUrlsById(list_id)
	assert.Nil(err)
	assert.Equal(1, len(urls), "There should be a single URL in the list")

	err = db.InsertClient("client", "0.0.0.0")
	assert.Nil(err)

	client, err := db.GetClient("client")
	assert.Nil(err)
	assert.Equal(DefaultList, client.UrlListId, "Client should not be assigned to a list yet")

	err = db.AssignClientToList("test", "client")
	assert.Nil(err)

	client, err = db.GetClient("client")
	assert.Nil(err)
	assert.Equal(list_id, client.UrlListId, "Client should be assigned to the test list")

	err = db.DeleteList("test")
	assert.Nil(err)

	urls, err = db.FetchListUrlsById(list_id)
	assert.Nil(err)
	assert.Equal(0, len(urls), "There should be no URLs left in the deleted list")

	client, err = db.GetClient("client")
	assert.Nil(err)
	assert.Equal(DefaultList, client.UrlListId, "Client should have been removed from the test list")
}
