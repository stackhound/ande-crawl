package main

import (
	"github.com/stackhound/ande-crawl/crawl"
	"github.com/stackhound/ande-crawl/db"
	"log"
	"strconv"
)

func main() {
	records, err := db.GetAvailableNIS()
	if err != nil {
		log.Fatal("Couldn't get the data:", err)
	}

	for _, nis := range records {
		t := strconv.FormatInt(nis.NIS, 10)
		log.Println(t)
		consumption, amount, err := crawl.FetchConsumption(string(t))

		if err != nil {
			log.Fatal("Couldn't fetch consumption!")

			// Skip this:
			continue
		}

		// No errors, store the record:
		record := db.ConsumptionRecord{}
		record.Consumption = consumption
		record.Amount = amount
		db.StoreConsumptionRecord(&record)
	}
}
