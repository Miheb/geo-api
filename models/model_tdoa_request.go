package models

type TdoaRequest struct {
	LoRaWAN []GatewayReceptionTdoa `json:"lorawan"`
}