package validator

import (
	extvalidator "github.com/go-playground/validator/v10"

	"github.com/linden-honey/linden-honey-scraper-go/pkg/song/scraper"
)

// validator represents the scraper.Validator implementation
type validator struct {
}

// NewValidator returns a pointer to the new instance of validator
func NewValidator() scraper.Validator {
	return &validator{}
}

// Validate returns true if structure is valid
func (v *validator) Validate(s interface{}) bool {
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
