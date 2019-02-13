package main

import (
	"fmt"
	"log"
	"net/http"
	"path/filepath"

	"github.com/campus-iot/geo-api/api"
	"github.com/gorilla/mux"
	"github.com/xeipuuv/gojsonschema"
)

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
	routeur.HandleFunc("/tdoa", api.GetTdoa)
	http.Handle("/", routeur)
	print("Server started")
	log.Fatal(http.ListenAndServe(":8081", routeur))
}
