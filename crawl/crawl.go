package crawl

import (
	"log"
)

// FetchConsumption gets the consumption for a given NIS.
func FetchConsumption(nis string) (int64, int64, error) {
	log.Printf("Fetching power consumption for %s", nis)
	var consumption, amount int64
	consumption = 400
	amount = 350000
	return consumption, amount, nil
}
