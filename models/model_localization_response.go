package models

type LocalizationResponse struct {
	Result *LocationEstimate `json:"result,omitempty"`

	Warnings []string `json:"warnings,omitempty"`

	Errors []string `json:"errors,omitempty"`

	CorrelationId string `json:"correlationId,omitempty"`
}
