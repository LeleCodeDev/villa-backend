package appError

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"reflect"
	"strings"
)

func typeName(t reflect.Type) string {
	switch t.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return "an integer"
	case reflect.Float32, reflect.Float64:
		return "a number"
	case reflect.String:
		return "a string"
	case reflect.Bool:
		return "a boolean"
	case reflect.Slice, reflect.Array:
		return "an array"
	case reflect.Struct, reflect.Map:
		return "an object"
	}
	return t.String()
}

func GetBindError(err error) (map[string]any, bool) {
	var typeErr *json.UnmarshalTypeError
	var syntaxErr *json.SyntaxError

	switch {
	case errors.As(err, &typeErr):
		if typeErr.Field == "" {
			return map[string]any{"body": "request body has the wrong JSON type"}, true
		}
		result := make(map[string]any)
		path := strings.Split(typeErr.Field, ".")
		field := path[len(path)-1]
		setNested(result, path, fmt.Sprintf("%s must be %s", field, typeName(typeErr.Type)))
		return result, true

	case errors.As(err, &syntaxErr):
		return map[string]any{"body": "request body is not valid JSON"}, true

	case errors.Is(err, io.EOF), errors.Is(err, io.ErrUnexpectedEOF):
		return map[string]any{"body": "request body is empty or incomplete"}, true
	}

	return nil, false
}
