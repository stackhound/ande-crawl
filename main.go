package main

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/jasonlvhit/gocron"
	"github.com/stackhound/ande-crawl/crawl"
	"github.com/stackhound/ande-crawl/db"
	"github.com/stackhound/ande-crawl/status"
)

func main() {
	// Con esto metemos el servidor web en una goroutine
	go status.Listen()

	gocron.Every(1).Day().Do(dayCheck)
	// remove, clear and next_run
	_, time := gocron.NextRun()
	fmt.Println(time)
	// function Start start all the pending jobs
	<-gocron.Start()

}

func dayCheck() {
	log.Println("Another day under the sun!")
	//This is updated daily so we check if it is our desired date every day
	currentTime := time.Now().Local()
	day := currentTime.Day()
	if day != 28 {
		return
	}
	doCrawl()
}

func doCrawl() {
	records, err := db.GetAvailableNIS()
	if err != nil {
		log.Fatal("Couldn't get the data:", err)
	}

	for _, nis := range records {
		t := strconv.FormatInt(nis.NIS, 10)
		consumption, amount, pendingBills, expirationDate, err := crawl.FetchConsumption(t)

		if err != nil {
			log.Fatal("Couldn't fetch consumption!")

			// Skip this:
			continue
		}

		// No errors, store the record:
		record := db.ConsumptionRecord{}
		record.NIS = t
		record.Consumption = consumption
		record.Amount = amount
		record.PendingBills = pendingBills
		record.Expiration = expirationDate
		err = db.StoreConsumptionRecord(&record)
		if err != nil {
			log.Println("Couldn't insert data:", err)
		}
		log.Println("Done with this one")
	}
	status.S.Iterations++
	log.Println("Goodbye")
}
