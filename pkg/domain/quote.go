package domain

// Quote represents a domain object
type Quote struct {
	Phrase string `validate:"required",json:"phrase"`
}
