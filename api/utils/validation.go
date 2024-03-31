package utils

import (
	"net/mail"
	"reflect"
	"strings"

	"github.com/coderero/erochat-server/types"
	"github.com/go-playground/validator/v10"
)

func ConvertValidationErrors(err error) []types.Error {
	var errors []types.Error
	validationErrors := err.(validator.ValidationErrors)
	for _, e := range validationErrors {
		errors = append(errors, types.Error{
			Field:  strings.ToLower(e.Field()),
			Reason: msgForTag(e.Tag()),
		})
	}
	return errors

}

// Code Reference:
// https://github.com/go-playground/validator/issues/258#issuecomment-257281334
func ValidatorTagFunc(fld reflect.StructField) string {
	field := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
	if field == "-" {
		return ""
	}

	return field
}

func msgForTag(tag string) string {
	switch tag {
	case "required":
		return "field is required"
	case "email":
		return "email is invalid"
	case "min":
		return "request doesn't satisfy minimum length"
	case "max":
		return "field is too long to process"
	case "gt":
		return "field should be greater than as mentioned"
	case "lt":
		return "field should be less than as mentioned"
	case "alphanum":
		return "field should be alphanumeric"
	case "alpha":
		return "field should only contains alphabets"
	case "numeric":
		return "field should only contain numeric values"
	}
	return "unknown error occured"
}

func IsValidEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}
