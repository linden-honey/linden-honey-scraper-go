package fetcher

import (
	sdkerrors "github.com/linden-honey/linden-honey-sdk-go/errors"
)

// Validate a Fetcher and returns an error if validation is failed
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

func (cfg RetryConfig) Validate() error {
	if cfg.Attempts <= 0 {
		return sdkerrors.NewInvalidValueError("Attempts", sdkerrors.ErrNonPositiveNumber)
	}

	if cfg.Factor <= 0 {
		return sdkerrors.NewInvalidValueError("Factor", sdkerrors.ErrNonPositiveNumber)
	}

	return nil
}
