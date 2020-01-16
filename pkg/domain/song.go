package domain

// Song represents a domain object
type Song struct {
	Title  string
	Author string
	Album  string
	Verses []Verse
}

// GetQuotes returns all quotes from the song
func (s Song) GetQuotes() []Quote {
	quotes := make([]Quote, 0)
	for _, verse := range s.Verses {
		quotes = append(quotes, verse.Quotes...)
	}
	return quotes
}
