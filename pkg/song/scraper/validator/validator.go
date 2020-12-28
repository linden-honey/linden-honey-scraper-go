package validator

import (
	"fmt"

	"github.com/linden-honey/linden-honey-scraper-go/pkg/song"
)

// Validator represents the default validator implementation
type Validator struct {
}

// NewValidator returns a pointer to the new instance of Validator
func NewValidator() (*Validator, error) {
	return &Validator{}, nil
}

func (v *Validator) ValidateSong(s song.Song) error {
	if s.Title == "" {
		return NewMissingRequiredFieldError("title")
	}

	if len(s.Verses) == 0 {
		return NewMissingRequiredFieldError("verses")
	}

	for i, v := range s.Verses {
		if err := validateVerse(v); err != nil {
			return fmt.Errorf("verses[%d] is invalid: %w", i, err)
		}
	}

	return nil
}

func validateVerse(v song.Verse) error {
	if len(v.Quotes) == 0 {
		return NewMissingRequiredFieldError("quotes")
	}

	for i, q := range v.Quotes {
		if err := validateQuote(q); err != nil {
			return fmt.Errorf("quotes[%d] is invalid: %w", i, err)
		}
	}

	return nil
}

func validateQuote(q song.Quote) error {
	if q.Phrase == "" {
		return NewMissingRequiredFieldError("phrase")
	}

	return nil
}

func (v *Validator) ValidatePreview(p song.Preview) error {
	if p.ID == "" {
		return NewMissingRequiredFieldError("id")
	}
	if p.Title == "" {
		return NewMissingRequiredFieldError("title")
	}

	return nil
}
