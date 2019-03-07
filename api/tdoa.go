package api

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/campus-iot/geo-api/models"
	"github.com/campus-iot/geo-api/utils"
)

func GetTdoa(w http.ResponseWriter, r *http.Request) {

	// Load JSON
	body, _ := ioutil.ReadAll(r.Body)
	log.Println("GetTdoa request received : " + string(body))

	// Unmarshal JSON
	var tdoaRequest models.TdoaRequest
	errUnmarshal := json.Unmarshal(body, &tdoaRequest)
	if errUnmarshal != nil {
		http.Error(w, "Internal error : " + errUnmarshal.Error(), http.StatusInternalServerError)
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

	// Send response
	io.WriteString(w, string(jsonResponse))
}
