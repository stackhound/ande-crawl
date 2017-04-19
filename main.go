package main

import (
	"log"
	"strconv"

	"github.com/stackhound/ande-crawl/crawl"
	"github.com/stackhound/ande-crawl/db"
	"github.com/stackhound/ande-crawl/status"
)

func main() {
	// Con esto metemos el servidor web en una goroutine
	go status.Listen()

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
