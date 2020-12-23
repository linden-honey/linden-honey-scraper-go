package song

// Quote represents a domain object
type Quote struct {
	Phrase string `validate:"required" json:"phrase"`
}

// Verse represents a domain object
type Verse struct {
	Quotes []Quote `validate:"required" json:"quotes"`
}

// Song represents a domain object
type Song struct {
	Title  string  `validate:"required" json:"title"`
	Author string  `json:"author,omitempty"`
	Album  string  `json:"album,omitempty"`
	Verses []Verse `validate:"required" json:"verses"`
}

// GetQuotes returns all quotes from the song
func (s *Song) GetQuotes() (quotes []Quote) {
	for _, verse := range s.Verses {
		quotes = append(quotes, verse.Quotes...)
	}
	return quotes
}

// Preview represents a domain object
type Preview struct {
	ID    string `validate:"required" json:"id"`
	Title string `validate:"required" json:"title"`
}
