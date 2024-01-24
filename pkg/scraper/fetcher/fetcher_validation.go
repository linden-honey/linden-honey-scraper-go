package fetcher

import (
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

	return nil
}
