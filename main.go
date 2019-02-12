package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"

	"./swagger"
	"github.com/gorilla/mux"
)

func Test(w http.ResponseWriter, r *http.Request) {
	// A very simple health check.
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	b, err := ioutil.ReadAll(r.Body)

	// METTRE LE CODE DE WIWI QUI CALCULE
	var response swagger.LocalizationResponse
	tmp := swagger.LocationEstimate{
		Latitude:  1.12345,
		Longitude: 1.22345,
		Altitude:  1.32345,
		Accuracy:  4.5,
	}

	response.Result = &tmp

	b, err := json.Marshal(&response)
	if err != nil {
		fmt.Println(err)
		return
	}

	// In the future we could report back on the status of our DB, or our cache
	// (e.g. Redis) by performing a simple PING, and include them in the response.
	io.WriteString(w, string(b))
}

func main() {
	routeur := mux.NewRouter()
	routeur.HandleFunc("/CalculTriloc", Test)
	http.Handle("/", routeur)
	print("Server started")
	log.Fatal(http.ListenAndServe("localhost:8080", routeur))
}
