package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/campus-iot/geo-api/models"
)

func GetTdoa(w http.ResponseWriter, r *http.Request) {
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
