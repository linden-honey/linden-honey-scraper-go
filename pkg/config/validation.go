package config

import (
	"strings"

	sdkerrors "github.com/linden-honey/linden-honey-sdk-go/errors"
)

// Validate validates a [Config] and returns an error if validation is failed.
func (cfg Config) Validate() error {
	if err := cfg.Server.Validate(); err != nil {
		return sdkerrors.NewInvalidValueError("Server", err)
	}

	if err := cfg.Scrapers.Validate(); err != nil {
		return sdkerrors.NewInvalidValueError("Scrapers", err)
	}

	return nil
}

// Validate validates a [ServerConfig] and returns an error if validation is failed.
func (cfg ServerConfig) Validate() error {
	if strings.TrimSpace(cfg.Host) == "" {
		return sdkerrors.NewInvalidValueError("Host", sdkerrors.ErrEmptyValue)
	}

	if cfg.Port <= 0 {
		return sdkerrors.NewInvalidValueError("Port", sdkerrors.ErrNonPositiveNumber)
	}

	if err := cfg.Health.Validate(); err != nil {
		return sdkerrors.NewInvalidValueError("Health", err)
	}

	if err := cfg.Spec.Validate(); err != nil {
		return sdkerrors.NewInvalidValueError("Spec", err)
	}

	return nil
}

// Validate validates a [HealthConfig] and returns an error if validation is failed.
func (cfg HealthConfig) Validate() error {
	if strings.TrimSpace(cfg.Path) == "" {
		return sdkerrors.NewInvalidValueError("Path", sdkerrors.ErrEmptyValue)
	}

	return nil
}

// Validate validates a [SpecConfig] and returns an error if validation is failed.
func (cfg SpecConfig) Validate() error {
	if strings.TrimSpace(cfg.FilePath) == "" {
		return sdkerrors.NewInvalidValueError("FilePath", sdkerrors.ErrEmptyValue)
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
	if strings.TrimSpace(cfg.BaseURL) == "" {
		return sdkerrors.NewInvalidValueError("BaseURL", sdkerrors.ErrEmptyValue)
	}

	return nil
}
