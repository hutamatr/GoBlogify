package helpers

import (
	"strconv"
	"unicode"

	"errors"

	"github.com/go-playground/validator/v10"
)

var Validate = validator.New(validator.WithRequiredStructEnabled())

func ValidationPassword(field validator.FieldLevel) bool {

	length, err := strconv.Atoi(field.Param())

	PanicError(err, "Failed convert param")

	var (
		hasMinLen  = false
		hasUpper   = false
		hasLower   = false
		hasNumber  = false
		hasSpecial = false
	)

	value := field.Field().String()

	if len(value) >= length {
		hasMinLen = true
	}

	for _, char := range value {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsNumber(char):
			hasNumber = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}

	return hasMinLen && hasUpper && hasLower && hasNumber && hasSpecial
}

func ValidationConfirmPassword(field validator.FieldLevel) bool {
	value, _, _, ok := field.GetStructFieldOK2()

	if !ok {
		PanicError(errors.New("failed get confirm password field"), "failed get confirm password field")
	}

	value2 := field.Field().String()

	return value.String() == value2
}

func CustomValidation() error {
	if err := Validate.RegisterValidation("password", ValidationPassword); err != nil {
		return err
	}

	if err := Validate.RegisterValidation("confirm_password", ValidationConfirmPassword); err != nil {
		return err
	}
	return nil
}
