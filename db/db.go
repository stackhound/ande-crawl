package db

import (
	"errors"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"log"
	"time"
)

var (
	databaseName          = "ande"
	userCollection        = "users"
	consumptionCollection = "consumptions"
)

// ConsumptionRecord represents a data structure for the JSON document to be stored.
type ConsumptionRecord struct {
	NIS          string    `json:"nis" bson:"nis"`
	Consumption  int64     `json:"consumption" bson:"consumption"`
	Amount       int64     `json:"amount" bson:"amount"`
	PendingBills int64     `json:"pendingBills" bson:"pendingBills"`
	Expiration   time.Time `json:"expiration" bson:"expiration"`
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

	c := session.DB(databaseName).C(userCollection)
	err = c.Find(bson.M{}).All(&users)

	return users, err
}

// StoreConsumptionRecord stores the consumption record.
func StoreConsumptionRecord(record *ConsumptionRecord) (err error) {
	log.Println("Storing consumption record: ", record)
	var session *mgo.Session
	session, err = getSession()
	defer session.Close()
	if err != nil {
		log.Println("It seems like there was an error: ", err)
	}
	c := session.DB(databaseName).C(consumptionCollection)
	err = c.Insert(&ConsumptionRecord{record.NIS, record.Consumption, record.Amount, record.PendingBills, record.Expiration})
	return err
}
