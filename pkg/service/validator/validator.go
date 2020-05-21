package validator

import (
	"github.com/go-playground/validator/v10"
	log "github.com/sirupsen/logrus"
)

// Validator represents the validator interface
type Validator interface {
	Validate(s interface{}) bool
}

// defaultValidator represents the default validator implementation
type defaultValidator struct {
	logger *log.Logger
}

// NewDefaultValidator returns a pointer to the new instance of defaultValidator
func NewDefaultValidator(logger *log.Logger) *defaultValidator {
	return &defaultValidator{
		logger: logger,
	}
}

// Validate returns true if structure is valid and logs errors
func (v *defaultValidator) Validate(s interface{}) bool {
	validate := validator.New()
	err := validate.Struct(s)
	isValid := err == nil
	if !isValid {
		//TODO return error with readable string representation
		validationErrors := err.(validator.ValidationErrors)
		v.logger.Printf("Validation failed for %v", s)
		v.logger.Println(validationErrors)
	}
	return isValid
}
