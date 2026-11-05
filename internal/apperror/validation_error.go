package appError

import (
	"strings"

	"github.com/go-playground/validator/v10"
)

var messages = map[string]string{
	"required": "{field} is required",
	"email":    "{field} must be a valid email",
	"min":      "{field} must contain at least {param} characters",
	"max":      "{field} must contain at most {param} characters",
	"gt":       "{field} must be greater than {param}",
	"gte":      "{field} must be greater than or equal {param}",
	"lt":       "{field} must be less than {param}",
	"lte":      "{field} must be less than or equal {param}",
}

var (
	indexReplacer = strings.NewReplacer("[", ".", "]", "")
	dateReplacer  = strings.NewReplacer("2006", "YYYY", "01", "MM", "02", "DD", "15", "HH", "04", "mm", "05", "ss")
)

func GetValidationError(errs validator.ValidationErrors) map[string]any {
	result := make(map[string]any)

	for _, e := range errs {
		path := strings.Split(indexReplacer.Replace(e.Namespace()), ".")[1:]
		setNested(result, path, buildMessage(e))
	}

	return result
}

func setNested(node map[string]any, path []string, msg string) {
	for _, key := range path[:len(path)-1] {
		child, ok := node[key].(map[string]any)
		if !ok {
			child = make(map[string]any)
			node[key] = child
		}
		node = child
	}
	node[path[len(path)-1]] = msg
}

func buildMessage(e validator.FieldError) string {
	param := e.Param()

	switch e.Tag() {
	case "oneof":
		param = strings.ReplaceAll(param, " ", ", ")
		return strings.NewReplacer("{field}", e.Field(), "{param}", param).
			Replace("{field} must be one of these values: {param}")
	case "datetime":
		param = dateReplacer.Replace(param)
		return strings.NewReplacer("{field}", e.Field(), "{param}", param).
			Replace("{field} must be a valid date in format {param}")
	}

	tmpl, ok := messages[e.Tag()]
	if !ok {
		tmpl = "{field} is invalid"
	}
	return strings.NewReplacer("{field}", e.Field(), "{param}", param).Replace(tmpl)
}
