package crawl

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var (
	endpointUrl = "http://201.217.43.238:9080/consulta/consulta_02.php"
)

var (
	layout = "2006-01-02"
	/*Following values taken from http://www.ande.gov.py/calcule_consumo.php */
	type1Ammount = 312.55 //real one:  311.55
	type2Ammount = 350.89 //real one: 349.89
	type3Ammount = 365.45 //real one: 365.45
)

type Result struct {
	InvoiceCount   int64     `json:"invoiceCount" bson:"invoiceCount"`
	Amount         int64     `json:"amount" bson:"amount"`
	ExpirationDate time.Time `json:"expirationDate" bson:"invoiceCount"`
}

// FetchConsumption gets the consumption for a given NIS.
func FetchConsumption(nis string) (int64, int64, int64, time.Time, error) {
	log.Printf("Fetching power consumption for %s", nis)
	var consumption, amount, pendingBills int64
	var expirationDate time.Time
	result, err := query(nis)

	if result.Amount > 0 && err == nil {
		amount = result.Amount
		consumption = amount / int64(type1Ammount)
		expirationDate = result.ExpirationDate
		pendingBills = result.InvoiceCount
		/*	if result.InvoiceCount == 0 {
				//fmt.Println("Factura al dia, monto del ultimo ciclo: Gs.", result.Amount)
			} else {
				//fmt.Println("Factura con boletas acumuladas, boletas:", result.InvoiceCount)
				//fmt.Println("Monto total es:", result.Amount)
			}*/
		//fmt.Println("Vence el:", result.ExpirationDate)
	}
	return consumption, amount, pendingBills, expirationDate, nil

}

func query(nis string) (result Result, err error) {
	values := url.Values{
		"name": {nis},
	}
	var resp *http.Response
	resp, err = http.PostForm(endpointUrl, values)
	if err == nil {
		var body []byte
		body, err = ioutil.ReadAll(resp.Body)
		bodyStr := string(body)

		expirationDateExpr := regexp.MustCompile("Fecha de vencimiento (.{10})")
		expirationDateMatch := expirationDateExpr.FindStringSubmatch(bodyStr)
		if len(expirationDateMatch) >= 1 {
			exp := expirationDateMatch[1]
			result.ExpirationDate, err = time.Parse(layout, exp)

			if err != nil {
				log.Println(err)

			}
		}

		amountExpr := regexp.MustCompile("Total Gs.: (.*) comisión")
		amountMatch := amountExpr.FindStringSubmatch(bodyStr)
		if len(amountMatch) >= 1 {
			amount := amountMatch[1]
			amount = strings.Replace(amount, ".", "", -1)

			var amountInt int64
			amountInt, err = strconv.ParseInt(amount, 10, 64)
			result.Amount = amountInt
		}

		countExpr := regexp.MustCompile("(.) facturas pendientes")
		countMatch := countExpr.FindStringSubmatch(bodyStr)
		if len(countMatch) >= 1 {
			count := countMatch[1]
			var countInt int64
			countInt, err = strconv.ParseInt(count, 10, 64)
			result.InvoiceCount = countInt
		}

	} else {
		log.Fatal("There was an error:", err)
	}

	return result, err
}
