package crawl

import (
	"log"
)

// FetchConsumption gets the consumption for a given NIS.
func FetchConsumption(nis string) (string, string, error) {
	log.Printf("Fetching power consumption for %s", nis)
	var consumption, amount string
	consumption = "400"
	amount = "350000"
	return consumption, amount, nil
}
