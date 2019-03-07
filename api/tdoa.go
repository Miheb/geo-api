package api

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	//"path/filepath"

	"github.com/campus-iot/geo-api/models"
	"github.com/campus-iot/geo-api/utils"
	//"github.com/xeipuuv/gojsonschema"
)

func GetTdoa(w http.ResponseWriter, r *http.Request) {

	// Load JSON Schema
	//pathSchema, _ := filepath.Abs("schema/geo-schema.json")
	//schemaLoader := gojsonschema.NewReferenceLoader("file://" + pathSchema)
	body, _ := ioutil.ReadAll(r.Body)
	log.Println("Beginning GetTdoa"+string(body))
	//documentLoader := gojsonschema.NewStringLoader(string(body))


	// Validate JSON
	//result, err := gojsonschema.Validate(schemaLoader, documentLoader)
	/*if err != nil {
		http.Error(w, "Incorrect request : "+err.Error(), http.StatusBadRequest)
		log.Println("Error validate - " + err.Error())
		return
	}
	if !result.Valid() {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		log.Println("JSON Invalid - " + err.Error())
		return
	}*/

	// Unmarshal JSON
	var tdoaRequest models.TdoaRequest
	errUnmarshal := json.Unmarshal(body, &tdoaRequest)
	if errUnmarshal != nil {
		http.Error(w, "Internal error : "+errUnmarshal.Error(), http.StatusInternalServerError)
		log.Println("Error JSON Unmarshal - " + errUnmarshal.Error())
		return
	}

	// Verify the number of gateways
	if len(tdoaRequest.LoRaWAN) < 3 {
		http.Error(w, "Not enough gateways to locate, must be at least 3", http.StatusBadRequest)
		log.Println("Error not enough gateways to locate, must be at least 3")
		return
	}

	// Trilateration
	location := utils.Inter3(tdoaRequest.LoRaWAN[0], tdoaRequest.LoRaWAN[1], tdoaRequest.LoRaWAN[2])
	response := models.LocalizationResponse{
		Result: &location,
	}

	// Marshal JSON
	jsonResponse, errMarshal := json.Marshal(&response)
	if errMarshal != nil {
		http.Error(w, "Internal error", http.StatusInternalServerError)
		log.Println("Error JSON Marshal - " + errMarshal.Error())
		return
	}

	log.Println("Ending GetTdoa")

	// Send response
	io.WriteString(w, string(jsonResponse))
}
