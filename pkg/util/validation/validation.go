package validation

import (
	"log"

	validator "github.com/go-playground/validator/v10"
)

// Validate returns true if structure is valid
func Validate(s interface{}) bool {
	validate := validator.New()
	err := validate.Struct(s)
	isValid := err == nil
	if !isValid {
		validationErrors := err.(validator.ValidationErrors)
		log.Printf("Validation failed for %v", s)
		log.Println(validationErrors)
	}
	return isValid
}
