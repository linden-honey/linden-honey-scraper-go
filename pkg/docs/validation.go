package docs

import (
	"github.com/linden-honey/linden-honey-sdk-go/errors"
)

func (s Spec) Validate() error {
	if len(s) == 0 {
		return errors.NewRequiredValueError("Spec")
	}

	return nil
}
