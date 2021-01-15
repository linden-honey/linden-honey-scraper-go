package song

import (
	"fmt"

	"github.com/linden-honey/linden-honey-sdk-go/errors"
)

// Validate validates a Quote and returns an error if validation is failed
func (q Quote) Validate() error {
	if q.Phrase == "" {
		return errors.NewRequiredValueError("Phrase")
	}

	return nil
}

// Validate validates a Verse and returns an error if validation is failed
func (v Verse) Validate() error {
	if len(v.Quotes) == 0 {
		return errors.NewRequiredValueError("Quotes")
	}

	for i, q := range v.Quotes {
		if err := q.Validate(); err != nil {
			return fmt.Errorf("'Quotes[%d]' is invalid: %w", i, err)
		}
	}

	return nil
}

// Validate validates a Song and returns an error if validation is failed
func (s Song) Validate() error {
	if s.Title == "" {
		return errors.NewRequiredValueError("Title")
	}

	if len(s.Verses) == 0 {
		return errors.NewRequiredValueError("Verses")
	}

	for i, v := range s.Verses {
		if err := v.Validate(); err != nil {
			return fmt.Errorf("'Verses[%d]' is invalid: %w", i, err)
		}
	}

	return nil
}

// Validate validates a Preview and returns an error if validation is failed
func (p Preview) Validate() error {
	if p.ID == "" {
		return errors.NewRequiredValueError("ID")
	}
	if p.Title == "" {
		return errors.NewRequiredValueError("Title")
	}

	return nil
}
