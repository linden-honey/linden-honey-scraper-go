package validator

import (
	"fmt"
)

func NewMissingRequiredFieldError(key string) error {
	return fmt.Errorf("missing required field '%s'", key)
}

func NewInvalidFieldValueError(key string, value interface{}, rule string) error {
	return fmt.Errorf("fiild '%s' has invalid value '%v' - %s", key, value, rule)
}
