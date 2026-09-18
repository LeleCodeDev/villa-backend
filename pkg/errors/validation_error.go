package errors

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

func GetValidationError(err validator.ValidationErrors) map[string]string {
	errors := make(map[string]string)

	for _, e := range err {
		field := e.Field()

		switch e.Tag() {
		case "required":
			errors[field] = fmt.Sprintf("%s is required", field)
		case "email":
			errors[field] = fmt.Sprintf("%s must be a valid email", field)
		case "min":
			errors[field] = fmt.Sprintf("%s must contain at least %s characters", field, e.Param())
		case "max":
			errors[field] = fmt.Sprintf("%s must contain at most %s characters", field, e.Param())
		case "gt":
			errors[field] = fmt.Sprintf("%s must be greater than %s", field, e.Param())
		case "gte":
			errors[field] = fmt.Sprintf("%s must be greater than or equal %s", field, e.Param())
		case "lt":
			errors[field] = fmt.Sprintf("%s must be less than %s", field, e.Param())
		case "lte":
			errors[field] = fmt.Sprintf("%s must be less than or equal %s", field, e.Param())
		case "oneof":
			values := strings.Join(strings.Split(e.Param(), " "), ", ")
			errors[field] = fmt.Sprintf("%s must be one of these values: %s", field, values)
		default:
			errors[field] = fmt.Sprintf("%v is invalid", field)
		}
	}

	return errors
}
