package validators

import (
	"github.com/go-playground/validator/v10"
)

// CustomValidator struct to wrap the validator instance
type CustomValidator struct {
	validator *validator.Validate
}

// Validate the struct
func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}

// NewCustomValidator creates a new CustomValidator instance
func NewCustomValidator() *CustomValidator {
	return &CustomValidator{validator: validator.New()}
}
