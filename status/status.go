package status

import (
	"encoding/json"
	"net/http"
)

// StatusHandler (just commented to avoid annoying vscode underlines)
type statusHandler struct {
	Iterations int64 `json:"iterations"`
}

//S This is what we are going to be exporting
var S statusHandler

func (h statusHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(S)
}

// Listen (just commented to avoid annoying vscode underlines)
func Listen() {
	var h http.Handler
	h = statusHandler{}
	http.ListenAndServe(":5000", h)
}
