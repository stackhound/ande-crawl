package crawl

import (
	"fmt"
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

type Result struct {
	InvoiceCount   int64
	Amount         int64
	ExpirationDate string
}

// FetchConsumption gets the consumption for a given NIS.
func FetchConsumption(nis string) (string, string, error) {
	log.Printf("Fetching power consumption for %s", nis)
	var consumption, amount string
	consumption = "400"
	amount = "350000"
	return consumption, amount, nil
}

func query(nis string) (result Result, err error) {
	values := url.Values{
		"name": {nis},
	}
	var resp *http.Response
	resp, err = http.PostForm(endpointUrl, values)
	fmt.Println(resp, err)
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

func main() {
	// var nis = "1427205"
	var nis = "1427216"
	fmt.Println("Consultando NIS:", nis)
	result, err := query(nis)

	if result.Amount > 0 && err == nil {
		if result.InvoiceCount == 0 {
			fmt.Println("Factura al dia, monto del ultimo ciclo: Gs.", result.Amount)
		} else {
			fmt.Println("Factura con boletas acumuladas, boletas:", result.InvoiceCount)
			fmt.Println("Monto total es:", result.Amount)
		}
		fmt.Println("Vence el:", result.ExpirationDate)
	}

	// fmt.Println(result, err)
}
