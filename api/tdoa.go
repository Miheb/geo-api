package api

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"path/filepath"

	"github.com/campus-iot/geo-api/models"
	"github.com/campus-iot/geo-api/utils"
	log "github.com/sirupsen/logrus"
	"github.com/xeipuuv/gojsonschema"
)

//GetTdoa getting tdoa
func GetTdoa(w http.ResponseWriter, r *http.Request) {

	pathSchema, _ := filepath.Abs("schema/geo-schema.json")
	schemaLoader := gojsonschema.NewReferenceLoader("file://" + pathSchema)

	body, _ := ioutil.ReadAll(r.Body)
	documentLoader := gojsonschema.NewStringLoader(string(body))
	log.WithFields(log.Fields{"======BEGINNING GetTDOA": "testé"}).Info("========")
	result, err := gojsonschema.Validate(schemaLoader, documentLoader)
	if err != nil {
		log.WithFields(log.Fields{"ERROR validate": err}).Info("========")
		http.Error(w, "Incorrect request", http.StatusBadRequest)
		return
	}

	if !result.Valid() {
		log.WithFields(log.Fields{"ERROR badrequest": err}).Info("========")
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	var gateways []models.GatewayReceptionTdoa
	err2 := json.Unmarshal(body, &gateways)
	if err2 != nil {
		fmt.Println(err2)
		http.Error(w, "Internal erro", http.StatusInternalServerError)
		return
	}

	if len(gateways) < 3 {
		http.Error(w, "Not enough gateways to locate, must be at least 3", http.StatusBadRequest)
		return
	}

	location := utils.Inter3(gateways[0], gateways[1], gateways[2])

	response := models.LocalizationResponse{
		Result: &location,
	}

	jsonResponse, err := json.Marshal(&response)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Internal erro", http.StatusInternalServerError)
		return
	}

	log.WithFields(log.Fields{"ENDING GETTDOA": jsonResponse}).Info("========")

	io.WriteString(w, string(jsonResponse))
}
