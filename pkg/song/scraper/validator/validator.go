package validator

import (
	extvalidator "github.com/go-playground/validator/v10"
)

// Validator represents the default validator implementation
type Validator struct {
}

// NewValidator returns a pointer to the new instance of Validator
func NewValidator() (*Validator, error) {
	return &Validator{}, nil
}

// Validate returns true if structure is valid
func (v *Validator) Validate(s interface{}) bool {
	//TODO change to interface filed of validator or rewrite validator at all
	validate := extvalidator.New()

	err := validate.Struct(s)
	isValid := err == nil
	if !isValid {
		//TODO return error with readable string representation
		//validationErrors := err.(validator.ValidationErrors)
	}

	return isValid
}
