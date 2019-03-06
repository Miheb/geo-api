package api

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"

	"github.com/campus-iot/geo-api/models"
	"github.com/campus-iot/geo-api/utils"
	"github.com/xeipuuv/gojsonschema"
)

//GetTdoa getting tdoa
func GetTdoa(w http.ResponseWriter, r *http.Request) {

	// Load JSON Schema
	pathSchema, _ := filepath.Abs("schema/geo-schema.json")
	schemaLoader := gojsonschema.NewReferenceLoader("file://" + pathSchema)
	body, _ := ioutil.ReadAll(r.Body)
	documentLoader := gojsonschema.NewStringLoader(string(body))

	log.Println("Beginning GetTdoa")

	// Validate JSON
	result, err := gojsonschema.Validate(schemaLoader, documentLoader)
	if err != nil {
		http.Error(w, "Incorrect request : "+err.Error(), http.StatusBadRequest)
		log.Println("Error validate - " + err.Error())
		return
	}
	if !result.Valid() {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		log.Println("JSON Invalid - " + err.Error())
		return
	}

	var gateways []models.GatewayReceptionTdoa
	err = json.Unmarshal(body, &gateways)
	if err != nil {
		http.Error(w, "Internal error : "+err.Error(), http.StatusInternalServerError)
		log.Println("Error JSON Unmarshal - " + err.Error())
		return
	}

	if len(gateways) < 3 {
		http.Error(w, "Not enough gateways to locate, must be at least 3", http.StatusBadRequest)
		log.Println("Error not enough gateways to locate, must be at least 3")
		return
	}

	location := utils.Inter3(gateways[0], gateways[1], gateways[2])

	response := models.LocalizationResponse{
		Result: &location,
	}

	jsonResponse, err := json.Marshal(&response)
	if err != nil {
		http.Error(w, "Internal error", http.StatusInternalServerError)
		log.Println("Error JSON Marshal - " + err.Error())
		return
	}

	log.Println("Ending GetTdoa")

	io.WriteString(w, string(jsonResponse))
}
