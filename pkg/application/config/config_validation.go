package config

import (
	"errors"

	sdkerrors "github.com/linden-honey/linden-honey-sdk-go/errors"
)

// Validate validates a [Config] and returns an error if validation is failed.
func (cfg Config) Validate() error {
	if err := cfg.Scrapers.Validate(); err != nil {
		return sdkerrors.NewInvalidValueError("Scrapers", err)
	}

	return nil
}

// Validate validates a [ScrapersConfig] and returns an error if validation is failed.
func (cfg ScrapersConfig) Validate() error {
	if err := cfg.Grob.Validate(); err != nil {
		return sdkerrors.NewInvalidValueError("Grob", err)
	}

	return nil
}

// Validate validates a [ScraperConfig] and returns an error if validation is failed.
func (cfg ScraperConfig) Validate() error {
	if cfg.BaseURL.String() == "" {
		return sdkerrors.NewInvalidValueError("BaseURL", sdkerrors.ErrEmptyValue)
	}

	if err := cfg.Retry.Validate(); err != nil {
		return sdkerrors.NewInvalidValueError("Retry", err)
	}

	return nil
}

// Validate validates a [RetryConfig] and returns an error if validation is failed.
func (cfg RetryConfig) Validate() error {
	if cfg.Attempts == 0 {
		return sdkerrors.NewInvalidValueError("Attempts", sdkerrors.ErrNonPositiveNumber)
	}

	if cfg.MinInterval <= 0 {
		return sdkerrors.NewInvalidValueError("MinInterval", sdkerrors.ErrNonPositiveNumber)
	}

	if cfg.MaxInterval <= 0 {
		return sdkerrors.NewInvalidValueError("MaxInterval", sdkerrors.ErrNonPositiveNumber)
	}

	if cfg.MinInterval > cfg.MaxInterval {
		return sdkerrors.NewInvalidValueError("MinInterval", errors.New("should be less than or equal to MaxInterval"))
	}

	if cfg.Factor <= 0 {
		return sdkerrors.NewInvalidValueError("Factor", sdkerrors.ErrNonPositiveNumber)
	}

	return nil
}
