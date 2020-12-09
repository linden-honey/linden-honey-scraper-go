package validator

import (
	"github.com/go-playground/validator/v10"
)

// Validator represents the validator implementation
type Validator struct {
}

// NewValidator returns a pointer to the new instance of Validator
func NewValidator() *Validator {
	return &Validator{}
}

// Validate returns true if structure is valid
func (v *Validator) Validate(s interface{}) bool {
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
