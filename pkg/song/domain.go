package song

import (
	"fmt"

	"github.com/linden-honey/linden-honey-sdk-go/errors"
)

// Quote represents a domain object
type Quote struct {
	Phrase string `json:"phrase"`
}

func (q Quote) Validate() error {
	if q.Phrase == "" {
		return errors.NewRequiredValueError("Phrase")
	}

	return nil
}

// Verse represents a domain object
type Verse struct {
	Quotes []Quote `json:"quotes"`
}

func (v Verse) Validate() error {
	if len(v.Quotes) == 0 {
		return errors.NewRequiredValueError("Quotes")
	}

	for i, q := range v.Quotes {
		if err := q.Validate(); err != nil {
			return fmt.Errorf("quotes[%d] is invalid: %w", i, err)
		}
	}

	return nil
}

// Song represents a domain object
type Song struct {
	Title  string  `json:"title"`
	Author string  `json:"author,omitempty"`
	Album  string  `json:"album,omitempty"`
	Verses []Verse `json:"verses"`
}

func (s Song) Validate() error {
	if s.Title == "" {
		return errors.NewRequiredValueError("field 'Quotes' is required")
	}

	if len(s.Verses) == 0 {
		return errors.NewRequiredValueError("Verses")
	}

	for i, v := range s.Verses {
		if err := v.Validate(); err != nil {
			return fmt.Errorf("verses[%d] is invalid: %w", i, err)
		}
	}

	return nil
}

// GetQuotes returns all quotes from the song
func (s Song) GetQuotes() (quotes []Quote) {
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

func (p Preview) Validate() error {
	if p.ID == "" {
		return errors.NewRequiredValueError("ID")
	}
	if p.Title == "" {
		return errors.NewRequiredValueError("Title")
	}

	return nil
}
