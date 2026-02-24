package validation

import "fmt"

type FieldError struct {
	Field   string `json:"field"`
	Tag     string `json:"tag"`
	Message string `json:"message"`
}

type ErrorResponse struct {
	Errors []FieldError `json:"errors"`
}

func (e ErrorResponse) Error() string {
	return fmt.Sprintf("validation failed: %d errors", len(e.Errors))
}
