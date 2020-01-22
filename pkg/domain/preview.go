package domain

// Preview represents a domain object
type Preview struct {
	ID    string `validate:"required",json:"id"`
	Title string `validate:"required",json:"title"`
}
