package utils

import (
	"fmt"
	"math"

	"github.com/campus-iot/geo-api/models"
)

func isEqualLat(gw1, gw2 models.GatewayReceptionTdoa) bool {
	return gw1.AntennaLocation.Latitude == gw2.AntennaLocation.Latitude
}

func isEqualLong(gw1, gw2 models.GatewayReceptionTdoa) bool {
	return gw1.AntennaLocation.Longitude == gw2.AntennaLocation.Longitude
}

func isEqualPlace(gw1, gw2 models.GatewayReceptionTdoa) bool {
	return isEqualLat(gw1, gw2) && isEqualLong(gw1, gw2)
}

func inter3(g1, g2, g3 models.GatewayReceptionTdoa) models.LocationEstimate {

	CX2 := 2 * (g2.AntennaLocation.Latitude - g1.AntennaLocation.Latitude)
	CX3 := 2 * (g3.AntennaLocation.Latitude - g1.AntennaLocation.Latitude)
	CY2 := 2 * (g2.AntennaLocation.Longitude - g1.AntennaLocation.Longitude)
	CY3 := 2 * (g3.AntennaLocation.Longitude - g1.AntennaLocation.Longitude)
	CR2 := math.Pow(g1.Rssi, 2) - math.Pow(g2.Rssi, 2) + (math.Pow(g2.AntennaLocation.Latitude, 2) + math.Pow(g2.AntennaLocation.Longitude, 2)) - (math.Pow(g1.AntennaLocation.Latitude, 2) + math.Pow(g1.AntennaLocation.Longitude, 2))
	CR3 := math.Pow(g1.Rssi, 2) - math.Pow(g3.Rssi, 2) + (math.Pow(g3.AntennaLocation.Latitude, 2) + math.Pow(g3.AntennaLocation.Longitude, 2)) - (math.Pow(g1.AntennaLocation.Latitude, 2) + math.Pow(g1.AntennaLocation.Longitude, 2))

	var CX float64
	var CY float64

	if isEqualPlace(g1, g2) || isEqualPlace(g1, g3) {
		fmt.Println("The three gateways are not distinct")
		return models.LocationEstimate{0, 0, 0, 0}
	} else if isEqualLat(g1, g2) {
		// 1 et 2 même x
		CY = CR2 / CY2
		CX = (CR3 - CY*CY3) / CX3
	} else if isEqualLat(g1, g3) {
		// 1 et 3 même x
		CY = CR3 / CY3
		CX = (CR2 - CY*CY2) / CX2
	} else if isEqualLong(g1, g2) {
		// 1 et 2 même y
		CX = CR2 / CX2
		CY = (CR3 - CX*CX3) / CY3
	} else if isEqualLong(g1, g3) {
		// 1 et 3 même y
		CX = CR3 / CX3
		CY = (CR2 - CX*CX2) / CY2
	} else {
		CYnum := CR2 - CX2*CR3/CX3
		CYden := CY2 - CX2*CY3/CX3

		CY = CYnum / CYden
		CX = (CR3 - CY*CY3) / CX3
	}

	return models.LocationEstimate{CX, CY, 0, 0}
}
