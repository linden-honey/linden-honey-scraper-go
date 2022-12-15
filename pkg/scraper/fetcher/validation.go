package fetcher

import (
	"errors"

	sdkerrors "github.com/linden-honey/linden-honey-sdk-go/errors"
)

// Validate validates a [Fetcher] and returns an error if validation is failed.
func (f Fetcher) Validate() error {
	if f.baseURL == nil {
		return sdkerrors.NewRequiredValueError("baseURL")
	}

	if f.encoding == nil {
		return sdkerrors.NewRequiredValueError("encoding")
	}

	if f.client == nil {
		return sdkerrors.NewRequiredValueError("client")
	}

	if f.retry != nil {
		if err := f.retry.Validate(); err != nil {
			return sdkerrors.NewInvalidValueError("retry", err)
		}
	}

	return nil
}

// Validate validates a [RetryConfig] and returns an error if validation is failed.
func (cfg RetryConfig) Validate() error {
	if cfg.Attempts <= 0 {
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
