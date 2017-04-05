package main

import (
	"log"

	"github.com/stackhound/ande-crawl/crawl"
	"github.com/stackhound/ande-crawl/db"
)

func main() {
	records := db.GetAvailableNIS()

	for _, nis := range records {
		consumption, amount, err := crawl.FetchConsumption(nis)

		if err != nil {
			log.Fatal("Couldn't fetch consumption!")

			// Skip this:
			continue
		}

		// No errors, store the record:
		record := db.ConsumptionRecord{}
		//record.Consumption = consumption
		//record.Amount = amount
		db.StoreConsumptionRecord(&record)
	}
}
