package utils

import (
	"regexp"

	"github.com/go-playground/validator/v10"
)

var Validate *validator.Validate

func RegisterAllValidators() {
	Validate = validator.New()
	Validate.RegisterValidation("PasswordValidation", PasswordValidation)
}

func PasswordValidation(fl validator.FieldLevel) bool {
	regex := regexp.MustCompile(`^(?=.*[A-Za-z])(?=.*\d)[A-Za-z\d]{8,}$`)
	value := fl.Field().String()
	return regex.MatchString(value)
}
