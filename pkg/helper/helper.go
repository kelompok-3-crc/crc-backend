package helper

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

func MapUnmarshalErrors(err error) map[string]string {
	errors := make(map[string]string)

	// Check if it's a JSON syntax error
	if syntaxErr, ok := err.(*json.SyntaxError); ok {
		errors["json"] = fmt.Sprintf("JSON syntax error at position %d: %v", syntaxErr.Offset, syntaxErr.Error())
		return errors
	}

	if typeErr, ok := err.(*json.UnmarshalTypeError); ok {
		errors[typeErr.Field] = fmt.Sprintf(
			"Field harus bertipe %v, bukan %v",
			typeErr.Type.String(),
			typeErr.Value,
		)
		return errors
	}

	// Handle invalid field names
	if invalidErr := strings.Split(err.Error(), "json: unknown field "); len(invalidErr) > 1 {
		fieldName := strings.Trim(invalidErr[1], "\"")
		errors[fieldName] = "Field tidak dikenali"
		return errors
	}

	// Default error
	errors["json"] = err.Error()
	return errors
}
func GetProjectRoot() (string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	return wd, nil
}
