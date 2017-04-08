package db

import (
	"errors"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"log"
	"time"
)

var (
	databaseName   = "ande"
	collectionName = "users"
)

// ConsumptionRecord represents a data structure for the JSON document to be stored.
type ConsumptionRecord struct {
	NIS          string    `json:"nis" bson:"nis"`
	Consumption  int64     `json:"consumption" bson:"consumption"`
	Amount       int64     `json:"amount" bson:"amount"`
	Expiration   time.Time `json:"expiration" bson:"expiration"`
	PendingBills int64     `json:"pending_bills" bson:"pending_bills"`
}

// User represents the User document.
type User struct {
	NIS      int64 `json:"nis" bson:"nis"`
	Category int   `json:"category" bson:"category"`
}

// getSession defines cluster and starts the connection
func getSession() (session *mgo.Session, err error) {
	uri := "mongodb://joel:12345678@ds155150.mlab.com:55150/ande" //os.Getenv("MONGO_URL")
	if uri == "" {
		return nil, errors.New("No connection string found")
	}
	return mgo.Dial(uri)
}

// GetAvailableNIS returns an array of NIS records from db
func GetAvailableNIS() (users []User, err error) {
	log.Println("Fetching NIS records.")

	var session *mgo.Session
	session, err = getSession()
	defer session.Close()
	if err != nil {
		return users, err
	}

	c := session.DB(databaseName).C(collectionName)
	err = c.Find(bson.M{}).All(&users)

	return users, err
}

// StoreConsumptionRecord stores the consumption record.
func StoreConsumptionRecord(record *ConsumptionRecord) {
	log.Println("Storing consumption record: ", record)
}
