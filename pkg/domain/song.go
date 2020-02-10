package domain

// Song represents a domain object
type Song struct {
	Title  string   `validate:"required" json:"title"`
	Author string   `json:"author,omitempty"`
	Album  string   `json:"album,omitempty"`
	Verses []*Verse `validate:"required" json:"verses"`
}

// GetQuotes returns all quotes from the song
func (s Song) GetQuotes() (quotes []*Quote) {
	for _, verse := range s.Verses {
		quotes = append(quotes, verse.Quotes...)
	}
	return quotes
}
