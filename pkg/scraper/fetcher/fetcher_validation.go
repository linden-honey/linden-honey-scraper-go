package fetcher

import (
	"errors"

	sdkerrors "github.com/linden-honey/linden-honey-sdk-go/errors"
)

// Validate validates a [Fetcher] and returns an error if validation is failed.
func (f Fetcher) Validate() error {
	errs := make([]error, 0)

	if f.baseURL.String() == "" {
		errs = append(errs, sdkerrors.NewRequiredValueError("baseURL"))
	}

	if f.client == nil {
		errs = append(errs, sdkerrors.NewRequiredValueError("client"))
	}

	return errors.Join(errs...)
}
