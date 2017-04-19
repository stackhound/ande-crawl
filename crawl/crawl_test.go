package crawl

import (
	"net/http"
	"testing"
	"time"
)

type testHandler struct {
	//
}

// Esta funcion va a ser llamada en cada request:
func (h testHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	/* Necesitamos responder algo como esto:
	<br/><p><center>El NIS 1427216 cuenta con 2 facturas pendientes de pago. <br />  Total Gs.: 521.000 comisión incluida. <br />Fecha de vencimiento 2017-04-17</p></center><br /><br />
	*/
	s := "<br/><p>testing</p></center><br /><br />"

	w.Write([]byte(s))
}

func init() {
	// Con esto hacemos override del valor original de endpointUrl:
	endpointUrl = "http://localhost:8000/consulta/consulta_02.php"

	// Corremos un servidor web en el 8000, usando un handler llamado TestHandler.
	var handler http.Handler
	handler = testHandler{}
	go http.ListenAndServe(":8000", handler)
	time.Sleep(1 * time.Second)
}

func TestFetchConsumption(t *testing.T) {
	nis := "123"
	consumption, amount, pendingBills, _, _ := FetchConsumption(nis)
	if consumption != 1669 {
		t.Fatal("Consumption does not match")
	}
	if amount != 521000 {
		t.Fatal("Amount does not match")
	}
	if pendingBills != 2 {
		t.Fatal("Pending Bills do not match")
	}
}
