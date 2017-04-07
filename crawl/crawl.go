package crawl

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
)

const (
	endpointUrl = "http://201.217.43.238:9080/consulta/consulta_02.php"
)

var (
	/*Following values taken from http://www.ande.gov.py/calcule_consumo.php */
	type1Ammount int64 = 312 //real one:  311.55
	type2Ammount int64 = 350 //real one: 349.89
	type3Ammount int64 = 365 //real one: 365.45
)

type Result struct {
	InvoiceCount   int64
	Amount         int64
	ExpirationDate string
}

// FetchConsumption gets the consumption for a given NIS.
func FetchConsumption(nis string) (int64, int64, error) {
	log.Printf("Fetching power consumption for %s", nis)
	var consumption, amount int64
	result, err := query(nis)

	if result.Amount > 0 && err == nil {
		amount = result.Amount
		consumption = amount / type1Ammount
		/*	if result.InvoiceCount == 0 {
				//fmt.Println("Factura al dia, monto del ultimo ciclo: Gs.", result.Amount)
			} else {
				//fmt.Println("Factura con boletas acumuladas, boletas:", result.InvoiceCount)
				//fmt.Println("Monto total es:", result.Amount)
			}*/
		//fmt.Println("Vence el:", result.ExpirationDate)
	}
	return consumption, amount, nil

}

func query(nis string) (result Result, err error) {
	values := url.Values{
		"name": {nis},
	}
	var resp *http.Response
	resp, err = http.PostForm(endpointUrl, values)
	//mt.Println(resp, err)
	if err == nil {
		var body []byte
		body, err = ioutil.ReadAll(resp.Body)
		bodyStr := string(body)

		expirationDateExpr := regexp.MustCompile("Fecha de vencimiento (.{10})")
		expirationDateMatch := expirationDateExpr.FindStringSubmatch(bodyStr)
		if len(expirationDateMatch) >= 1 {
			result.ExpirationDate = expirationDateMatch[1]
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

	}
	return result, err
}
