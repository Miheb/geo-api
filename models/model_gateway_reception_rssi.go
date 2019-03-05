package models

type GatewayReceptionRssi struct {
	GatewayId string `json:"gatewayId"`

	AntennaId int32 `json:"antennaId"`

	Rssi float64 `json:"rssi"`

	Snr float64 `json:"snr"`

	AntennaLocation *AntennaLocation `json:"antennaLocation,omitempty"`
}
