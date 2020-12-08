package validator

import (
	"github.com/go-playground/validator/v10"
)

// defaultValidator represents the default validator implementation
type defaultValidator struct {
}

// NewDefaultValidator returns a pointer to the new instance of defaultValidator
func NewDefaultValidator() *defaultValidator {
	return &defaultValidator{}
}

// Validate returns true if structure is valid and logs errors
func (v *defaultValidator) Validate(s interface{}) bool {
	//TODO change to interface filed of validator or rewrite validator at all
	validate := validator.New()

	err := validate.Struct(s)
	isValid := err == nil
	if !isValid {
		//TODO return error with readable string representation
		//validationErrors := err.(validator.ValidationErrors)
	}

	return isValid
}
