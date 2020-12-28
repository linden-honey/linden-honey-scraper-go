package song

// Quote represents a domain object
type Quote struct {
	Phrase string `json:"phrase"`
}

// Verse represents a domain object
type Verse struct {
	Quotes []Quote `json:"quotes"`
}

// Song represents a domain object
type Song struct {
	Title  string  `json:"title"`
	Author string  `json:"author,omitempty"`
	Album  string  `json:"album,omitempty"`
	Verses []Verse `json:"verses"`
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
	ID    string `json:"id"`
	Title string `json:"title"`
}
