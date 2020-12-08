package domain

// Verse represents a domain object
type Verse struct {
	Quotes []Quote `validate:"required" json:"quotes"`
}
