package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"path/filepath"

	"github.com/campus-iot/geo-API/swagger"
	"github.com/gorilla/mux"
	"github.com/xeipuuv/gojsonschema"
)

func Test(w http.ResponseWriter, r *http.Request) {
	// A very simple health check.
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	//b, err := ioutil.ReadAll(r.Body)

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

	pathSchema, _ := filepath.Abs("schema/geo-schema.json")
	pathDoc, _ := filepath.Abs("test/data.json")

	schemaLoader := gojsonschema.NewReferenceLoader("file:///" + pathSchema)

	documentLoader := gojsonschema.NewReferenceLoader("file:///" + pathDoc)

	result, err := gojsonschema.Validate(schemaLoader, documentLoader)

	if err != nil {
		panic(err.Error())
	}

	if result.Valid() {
		fmt.Printf("The document is valid\n")
	} else {
		fmt.Printf("The document is not valid. see errors :\n")
		for _, desc := range result.Errors() {
			fmt.Printf("- %s\n", desc)
		}
	}

	routeur := mux.NewRouter()
	routeur.HandleFunc("/CalculTriloc", Test)
	http.Handle("/", routeur)
	print("Server started")
	log.Fatal(http.ListenAndServe(":8081", routeur))
}
