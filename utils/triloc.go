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

func LatLonToXY(lat, lon float64) (float64, float64) {
	radius := 6371.0
	var x, y float64
	x = radius * lon * math.Cos(1.)
	y = radius * lat

	return x, y
}

//Origin is set at the intersection of greenwitch and the equator, perhaps it might be changed
func XYToLatLon(x, y float64) (float64, float64) {
	radius := 6371.0
	var lat, lon float64
	lat = y / radius
	lon = x / (radius * math.Cos(1.))

	return lat, lon
}

func convertReceptionTdoa(g models.GatewayReceptionTdoa) (float64, float64) {
	return LatLonToXY(g.AntennaLocation.Latitude, g.AntennaLocation.Longitude)
}

func convertResult(X, Y float64) (float64, float64) {
	return XYToLatLon(X, Y)
}

func Inter3(g1, g2, g3 models.GatewayReceptionTdoa) models.LocationEstimate {

	G1x, G1y := convertReceptionTdoa(g1)
	G2x, G2y := convertReceptionTdoa(g2)
	G3x, G3y := convertReceptionTdoa(g3)

	// CX2 := 2 * (g2.AntennaLocation.Latitude - g1.AntennaLocation.Latitude)
	// CX3 := 2 * (g3.AntennaLocation.Latitude - g1.AntennaLocation.Latitude)
	// CY2 := 2 * (g2.AntennaLocation.Longitude - g1.AntennaLocation.Longitude)
	// CY3 := 2 * (g3.AntennaLocation.Longitude - g1.AntennaLocation.Longitude)
	// CR2 := math.Pow(g1.Rssi, 2) - math.Pow(g2.Rssi, 2) + (math.Pow(g2.AntennaLocation.Latitude, 2) + math.Pow(g2.AntennaLocation.Longitude, 2)) - (math.Pow(g1.AntennaLocation.Latitude, 2) + math.Pow(g1.AntennaLocation.Longitude, 2))
	// CR3 := math.Pow(g1.Rssi, 2) - math.Pow(g3.Rssi, 2) + (math.Pow(g3.AntennaLocation.Latitude, 2) + math.Pow(g3.AntennaLocation.Longitude, 2)) - (math.Pow(g1.AntennaLocation.Latitude, 2) + math.Pow(g1.AntennaLocation.Longitude, 2))

	CX2 := 2 * (G2y - G1y)
	CX3 := 2 * (G3y - G1y)
	CY2 := 2 * (G2x - G1x)
	CY3 := 2 * (G3x - G1x)
	CR2 := math.Pow(float64(g1.Rssi), 2) - math.Pow(float64(g2.Rssi), 2) + (math.Pow(G2y, 2) + math.Pow(G2x, 2)) - (math.Pow(G1y, 2) + math.Pow(G1x, 2))
	CR3 := math.Pow(float64(g1.Rssi), 2) - math.Pow(float64(g3.Rssi), 2) + (math.Pow(G3y, 2) + math.Pow(G3x, 2)) - (math.Pow(G1y, 2) + math.Pow(G1x, 2))

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
	resultlat, resultlon := convertResult(CX, CY)
	return models.LocationEstimate{resultlat, resultlon, 0, 0}
}
