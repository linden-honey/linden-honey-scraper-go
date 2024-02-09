package flow

import (
	"errors"
	"strings"

	sdkerrors "github.com/linden-honey/linden-honey-sdk-go/errors"
)

func (dto RunSimpleFlowRequest) Validate() error {
	errs := make([]error, 0)

	if strings.TrimSpace(dto.ArtifactName) == "" {
		errs = append(errs, sdkerrors.NewRequiredValueError("ArtifactName"))
	}

	return errors.Join(errs...)
}
