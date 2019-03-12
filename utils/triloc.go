package utils

import (
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

func latLontoXY(lat, lon, LatOrigin, LonOrigin float64) Point {
	radius := 6371.0
	var x, y float64

	x = 2. * math.Pi * radius * math.Cos(((lat+LatOrigin)/2.)*(math.Pi/180.)) * ((lon - LonOrigin) / 360.)
	y = 2. * math.Pi * radius * ((lat - LatOrigin) / 360.)

	return Point{
		X: x,
		Y: y,
		Z: 0.,
	}
}

func xyToLatLon(p Point, LatOrigin, LonOrigin float64) (float64, float64) {
	radius := 6371.0
	var LatResult, LonResult float64
	LatResult = LatOrigin + ((360. * p.Y) / (2. * math.Pi * radius))
	LonResult = LonOrigin + ((360. * p.X) / (2. * math.Pi * radius * math.Cos(((LatResult+LatOrigin)/2.)*(math.Pi/180.))))

	return LatResult, LonResult
}

func convertReceptionTdoa(g1, g2, g3 models.GatewayReceptionTdoa) (Point, Point, Point) {
	p1 := Point{
		X: 0.,
		Y: 0.,
	}

	p2 := latLontoXY(g2.AntennaLocation.Latitude, g2.AntennaLocation.Longitude, g1.AntennaLocation.Latitude, g1.AntennaLocation.Longitude)
	p3 := latLontoXY(g3.AntennaLocation.Latitude, g3.AntennaLocation.Longitude, g1.AntennaLocation.Latitude, g1.AntennaLocation.Longitude)

	return p1, p2, p3

}

func convertResult(p Point, LatOrigin, LongOrigin float64) (float64, float64) {
	return xyToLatLon(p, LatOrigin, LongOrigin)
}

//Need to be modified, to match with our real Ptx which is the transmit power
func convertRssi(rssi float64) float64 {
	n := 2.8
	Ptx := -59.
	return math.Pow(10., ((Ptx - rssi) / (10. * n)))
}

func Inter3(g1, g2, g3 models.GatewayReceptionTdoa) models.LocationEstimate {

	p1, p2, p3 := convertReceptionTdoa(g1, g2, g3)

	p1.R = convertRssi(g1.Rssi)
	p2.R = convertRssi(g2.Rssi)
	p3.R = convertRssi(g3.Rssi)

	solution, _ := Solve(p1, p2, p3)

	resultlat, resultlon := convertResult(solution.First(), g1.AntennaLocation.Latitude, g1.AntennaLocation.Longitude)
	return models.LocationEstimate{resultlat, resultlon, 0, 0}
}

func sq(x float64) float64 {
	return math.Pow(x, 2)
}

// func main() {
// 	x, y := LatLonToXY(45.197416, 5.776886, 45.205206, 5.700256)
// 	x1, y1 := LatLonToXY(45.159076, 5.732446, 45.205206, 5.700256)
// 	x2, y2 := LatLonToXY(45.177715, 5.730671, 45.205206, 5.700256)
// 	lat, lon := XYToLatLon(2.3836356826726677, -3.0585678038019584, 45.205206, 5.700256)
// 	fmt.Println(x, y)
// 	fmt.Println(x1, y1)
// 	fmt.Println(x2, y2)
// 	fmt.Print("\n")
// 	fmt.Print(lat, lon)
// 	fmt.Print("\n")
//tm := time.Unix(1518328047, 0)
//fmt.Println(tm)

// var json3Sat = []byte(`[
// 	{
// 	  "gatewayId": "string",
// 	  "antennaId": 0,
// 	  "rssi": 3.886,
// 	  "snr": 0,
// 	  "toa": 212962,
// 	  "encryptedToa": "",
// 	  "antennaLocation": {
// 		"latitude": 45.205206,
// 		"longitude": 5.700256,
// 		"altitude": 0
// 	  }
// 	},
// 	{
// 	  "gatewayId": "string",
// 	  "antennaId": 0,
// 	  "rssi": 4.24,
// 	  "snr": 0,
// 	  "toa": 214143,
// 	  "encryptedToa": "",
// 	  "antennaLocation": {
// 		"latitude": 45.197416,
// 		"longitude": 5.776886,
// 		"altitude": 0
// 	  }
// 	},
// 	{
// 	  "gatewayId": "string",
// 	  "antennaId": 70,
// 	  "rssi": 2.091,
// 	  "snr": 0,
// 	  "toa": 206975,
// 	  "encryptedToa": "",
// 	  "antennaLocation": {
// 		"latitude": 45.159076,
// 		"longitude": 5.732446,
// 		"altitude": 0
// 	  }
// 	}
//   ]`)
// var sats []models.GatewayReceptionTdoa
// err := json.Unmarshal(json3Sat, &sats)
// if err != nil {
// 	fmt.Println("error:", err)
// }
// if len(sats) < 3 {
// 	fmt.Println("error: you should provide at least 3 gateways data")
// } else {
// 	fmt.Println(Inter3(sats[0], sats[1], sats[2]))
// }

// }
