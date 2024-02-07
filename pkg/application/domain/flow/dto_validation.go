package flow

import (
	"errors"
	"strings"

	sdkerrors "github.com/linden-honey/linden-honey-sdk-go/errors"
)

func (dto SimpleFlowInput) Validate() error {
	errs := make([]error, 0)

	if strings.TrimSpace(dto.OutputFileName) == "" {
		errs = append(errs, sdkerrors.NewRequiredValueError("OutputFileName"))
	}

	return errors.Join(errs...)
}
