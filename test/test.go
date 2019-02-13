package test

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/campus-iot/geo-api/models"
)

func Test(w http.ResponseWriter, r *http.Request) {
	// A very simple health check.
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	//b, err := ioutil.ReadAll(r.Body)

	// METTRE LE CODE DE WIWI QUI CALCULE
	var response models.LocalizationResponse
	tmp := models.LocationEstimate{
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

	var json3Sat = []byte(`[
	{
		"gatewayId": "string",
    "antennaId": 0,
    "rssi": 2,
    "snr": 0,
    "toa": 0,
    "encryptedToa": "",
    "antennaLocation": {
      "latitude": 4,
      "longitude": 0,
      "altitude": 0
		}
	},
	{
		"gatewayId": "string",
    "antennaId": 0,
    "rssi": 2,
    "snr": 0,
    "toa": 0,
    "encryptedToa": "",
    "antennaLocation": {
      "latitude": 2,
      "longitude": 2,
      "altitude": 0
		}
	},{
		"gatewayId": "string",
    "antennaId": 0,
    "rssi": 2,
    "snr": 0,
    "toa": 0,
    "encryptedToa": "",
    "antennaLocation": {
      "latitude": 4,
      "longitude": 4,
      "altitude": 0
		}
	}
	]`)
	var sats []models.GatewayReceptionTdoa
	err := json.Unmarshal(json3Sat, &sats)
	if err != nil {
		fmt.Println("error:", err)
	}
	fmt.Printf("%v", sats)
	if len(sats) < 3 {
		fmt.Println("error: you should provide at least 3 gateway data")
	} else {
		fmt.Println(inter3(sats[0], sats[1], sats[2]))
	}

}
