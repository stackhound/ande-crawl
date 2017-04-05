package db

import (
	//"encoding/json"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"log"
	"os"
)

var (
	mgoSession   *mgo.Session
	databaseName = "stackhound"
)

// ConsumptionRecord represents a data structure for the JSON document to be stored.
type ConsumptionRecord struct {
	NIS         string `json:"nis"`
	Consumption int64  `json:"consumption"`
	Amount      int64  `json:"amount"`
	//CreatedAt   int64  `json:"created_at"`
}

type nis struct {
	NIS string `json:"nis"`
}

type nises []nis //array of multiple nis

// getSession defines cluster and starts the connection
func getSession() *mgo.Session {
	uri := "mongodb://joel:12345678@ds153400.mlab.com:53400/stackhound" //os.Getenv("MONGO_URL")
	if uri == "" {
		log.Println("No se pudo encontrar la url ")
		os.Exit(1)
	}

	if mgoSession == nil {
		var err error
		mgoSession, err = mgo.Dial(uri)
		if err != nil {
			panic(err) // no, not really
		}
	}
	return mgoSession.Clone()
}

// withCollection lets you pass collection name and query
func withCollection(collection string, s func(*mgo.Collection) error) error {
	session := getSession()
	defer session.Close()
	c := session.DB(databaseName).C(collection)
	return s(c)
}

// GetAvailableNIS returns an array of NIS records from db
func GetAvailableNIS() []string {
	log.Println("Fetching NIS records.")
	var records []string
	var result nises
	query := func(c *mgo.Collection) error {
		err := c.Find(bson.M{}).All(&result)
		if err != nil {
			log.Fatal("Error al buscar el dato: ", err)
			//return
		}
		return err
	}
	withCollection("nis", query)

	for _, onenis := range result {

		records = append(records, onenis.NIS) //add this nis to the array
	}

	/*str, _ := json.MarshalIndent(result, "", " ")
	log.Printf("%s\n", str)*/

	//log.Println("NIS : ", result.NIS)
	//log.Println(availablenis)
	// Sample values:
	records = append(records, "noesnis")

	return records
}

// StoreConsumptionRecord stores the consumption record.
func StoreConsumptionRecord(record *ConsumptionRecord) {
	log.Println("Storing consumption record: ", record)
}
