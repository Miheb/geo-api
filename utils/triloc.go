package utils

import (
	"fmt"
	"math"

	"github.com/campus-iot/geo-API/swagger"
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

func LatLonToXY(lat, lon, LatOrigin, LonOrigin float64) (float64, float64) {
	radius := 6371.0
	var x, y float64

	x = 2. * math.Pi * radius * math.Cos(((lat+LatOrigin)/2.)*(math.Pi/180.)) * ((lon - LonOrigin) / 360.)
	y = 2. * math.Pi * radius * ((lat - LatOrigin) / 360.)

	return x, y
}

//Origin is set at the intersection of greenwitch and the equator, perhaps it might be changed
func XYToLatLon(x, y, LatOrigin, LonOrigin float64) (float64, float64) {
	radius := 6371.0
	var LatResult, LonResult float64
	LatResult = LatOrigin + ((360. * y) / (2. * math.Pi * radius))
	LonResult = LonOrigin + ((360. * x) / (2. * math.Pi * radius * math.Cos(((LatResult+LatOrigin)/2.)*(math.Pi/180.))))

	return LatResult, LonResult
}

func convertReceptionTdoa(g1, g2, g3 models.GatewayReceptionTdoa) (float64, float64, float64, float64, float64, float64) {
	var Xg1, Yg1 float64 = 0., 0.

	var Xg2, Yg2 float64 = LatLonToXY(g2.AntennaLocation.Latitude, g2.AntennaLocation.Longitude, g1.AntennaLocation.Latitude, g1.AntennaLocation.Longitude)

	var Xg3, Yg3 float64 = LatLonToXY(g3.AntennaLocation.Latitude, g3.AntennaLocation.Longitude, g1.AntennaLocation.Latitude, g1.AntennaLocation.Longitude)

	return Xg1, Yg1, Xg2, Yg2, Xg3, Yg3
}

func convertResult(X, Y, LatOrigin, LongOrigin float64) (float64, float64) {
	return XYToLatLon(X, Y, LatOrigin, LongOrigin)
}

func Inter3(g1, g2, g3 models.GatewayReceptionTdoa) models.LocationEstimate {

	G1x, G1y, G2x, G2y, G3x, G3y := convertReceptionTdoa(g1, g2, g3)

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
	CR2 := math.Pow(g1.Rssi, 2) - math.Pow(g2.Rssi, 2) + (math.Pow(G2y, 2) + math.Pow(G2x, 2)) - (math.Pow(G1y, 2) + math.Pow(G1x, 2))
	CR3 := math.Pow(g1.Rssi, 2) - math.Pow(g3.Rssi, 2) + (math.Pow(G3y, 2) + math.Pow(G3x, 2)) - (math.Pow(G1y, 2) + math.Pow(G1x, 2))

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
	resultlat, resultlon := convertResult(CX, CY, g1.AntennaLocation.Latitude, g1.AntennaLocation.Longitude)
	return models.LocationEstimate{resultlat, resultlon, 0, 0}
}

func sq(x float64) float64 {
	return math.Pow(x, 2)
}

func tdoa(g1, g2, g3 swagger.GatewayReceptionTdoa) swagger.LocationEstimate {

	v := 50.0

	cX1 := -2 * (g2.AntennaLocation.Latitude - g1.AntennaLocation.Latitude)
	cX2 := -2 * (g3.AntennaLocation.Latitude - g2.AntennaLocation.Latitude)
	cX3 := -2 * (g1.AntennaLocation.Latitude - g3.AntennaLocation.Latitude)
	cY1 := -2 * (g2.AntennaLocation.Longitude - g1.AntennaLocation.Longitude)
	cY2 := -2 * (g3.AntennaLocation.Longitude - g2.AntennaLocation.Longitude)
	cY3 := -2 * (g1.AntennaLocation.Longitude - g3.AntennaLocation.Longitude)
	cT1 := 2 * sq(v) * float64(g2.Toa-g1.Toa)
	cT2 := 2 * sq(v) * float64(g3.Toa-g2.Toa)
	cT3 := 2 * sq(v) * float64(g1.Toa-g3.Toa)
	cR1 := -sq(g2.AntennaLocation.Latitude) + sq(g1.AntennaLocation.Latitude) - sq(g2.AntennaLocation.Longitude) + sq(g1.AntennaLocation.Longitude) + sq(v)*(sq(float64(g2.Toa))-sq(float64(g1.Toa)))
	cR2 := -sq(g3.AntennaLocation.Latitude) + sq(g2.AntennaLocation.Latitude) - sq(g3.AntennaLocation.Longitude) + sq(g2.AntennaLocation.Longitude) + sq(v)*(sq(float64(g3.Toa))-sq(float64(g2.Toa)))
	cR3 := -sq(g1.AntennaLocation.Latitude) + sq(g3.AntennaLocation.Latitude) - sq(g1.AntennaLocation.Longitude) + sq(g3.AntennaLocation.Longitude) + sq(v)*(sq(float64(g1.Toa))-sq(float64(g3.Toa)))

	cXX1 := cX1 - cT1*cX2/cT2
	cXX2 := cX2 - cT2*cX3/cT3
	cYY1 := cY1 - cT1*cY2/cT2
	cYY2 := cY2 - cT2*cY3/cT3
	cRR1 := cR1 - cT1*cR2/cT2
	cRR2 := cR2 - cT2*cR3/cT3

	x := (cRR1 - cYY1*cRR2/cYY2) / (cXX1 - cYY1*cXX2/cYY2)
	y := (cRR2 - x*cXX2) / cYY2
	//	t := (cR3 - y*cY3 - x*cX3) / cT3

	return swagger.LocationEstimate{x, y, 0, 0}

}

// func main() {
// 	x, y := LatLonToXY(45.1845498, 5.7525638, 45.1777139, 5.7306643)
// 	lat, lon := XYToLatLon(1.7164330526576157, 0.7601173990493916, 45.1777139, 5.7306643)
// 	fmt.Print(x, y)
// 	fmt.Print("\n")
// 	fmt.Print(lat, lon)
// }
