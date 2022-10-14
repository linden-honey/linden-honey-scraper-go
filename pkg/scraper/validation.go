package scraper

import (
	sdkerrors "github.com/linden-honey/linden-honey-sdk-go/errors"
)

// Validate a Scraper and returns an error if validation is failed
func (scr Scraper) Validate() error {
	if scr.fetcher == nil {
		return sdkerrors.NewRequiredValueError("fetcher")
	}

	if scr.parser == nil {
		return sdkerrors.NewRequiredValueError("parser")
	}

	return nil
}
