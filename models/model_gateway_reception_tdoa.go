package models

type GatewayReceptionTdoa struct {
	GatewayId string `json:"gatewayId"`

	AntennaId int32 `json:"antennaId"`

	Rssi float64 `json:"rssi"`

	Snr float64 `json:"snr"`

	Toa int32 `json:"toa"`

	EncryptedToa string `json:"encryptedToa"`

	AntennaLocation *AntennaLocation `json:"antennaLocation,omitempty"`
}
