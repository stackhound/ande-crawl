package db

import (
	"log"
)

// ConsumptionRecord represents a data structure for the JSON document to be stored.
type ConsumptionRecord struct {
	NIS         string `json:"nis"`
	Consumption int64  `json:"consumption"`
	Amount      int64  `json:"amount"`
	//CreatedAt   int64  `json:"created_at"` //
}

// GetAvailableNIS returns an array of NIS records.
func GetAvailableNIS() []string {
	log.Println("Fetching NIS records.")

	var records []string
	// Sample values:
	records = append(records, "1000", "1001")

	return records
}

// StoreConsumptionRecord stores the consumption record.
func StoreConsumptionRecord(record *ConsumptionRecord) {
	log.Println("Storing consumption record: ", record)
}
