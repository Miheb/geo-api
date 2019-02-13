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
	"github.com/xeipuuv/gojsonschema"
)

func GetTdoa(w http.ResponseWriter, r *http.Request) {
	pathSchema, _ := filepath.Abs("schema/geo-schema.json")
	schemaLoader := gojsonschema.NewReferenceLoader("file://" + pathSchema)

	body, _ := ioutil.ReadAll(r.Body)
	documentLoader := gojsonschema.NewStringLoader(string(body))

	result, err := gojsonschema.Validate(schemaLoader, documentLoader)
	if err != nil {
		panic(err.Error())
	}

	if !result.Valid() {
		http.Error(w, "Bad Request", http.StatusBadRequest)
	}

	var gateways []models.GatewayReceptionTdoa
	err2 := json.Unmarshal(body, &gateways)
	if err2 != nil {
		fmt.Println(err2)
	}

	if len(gateways) < 3 {
		http.Error(w, "Not enough gateways to locate", http.StatusBadRequest)
	}

	fmt.Println()

	location := utils.Inter3(gateways[0], gateways[1], gateways[2])

	response := models.LocalizationResponse{
		Result: &location,
	}

	jsonResponse, err := json.Marshal(&response)
	if err != nil {
		fmt.Println(err)
		return
	}

	io.WriteString(w, string(jsonResponse))
}
