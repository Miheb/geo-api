package models

type GatewayReceptionTdoa struct {
	GatewayId string `json:"gatewayId"`
	AntennaId int `json:"antennaId"`
	Rssi int `json:"rssi"`
	Snr float64 `json:"snr"`
	Toa int `json:"toa"`
	EncryptedToa string `json:"encryptedToa"`

	AntennaLocation *AntennaLocation `json:"antennaLocation,omitempty"`
}
