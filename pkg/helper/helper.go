package helper

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
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
	// First try environment variable
	if root := os.Getenv("APP_ROOT"); root != "" {
		return root, nil
	}

	// Fallback to working directory
	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	// Navigate up until we find the go.mod file
	for {
		if _, err := os.Stat(filepath.Join(wd, "go.mod")); err == nil {
			return wd, nil
		}
		parent := filepath.Dir(wd)
		if parent == wd {
			return "", fmt.Errorf("could not find project root")
		}
		wd = parent
	}
}
