package validation

import (
	"fmt"
	"strings"

	govalidator "github.com/go-playground/validator/v10"
)

type Validator struct {
	v *govalidator.Validate
}

func New() *Validator {
	return &Validator{v: govalidator.New()}
}

func (v *Validator) Struct(input any) error {
	if err := v.v.Struct(input); err != nil {
		verrs, ok := err.(govalidator.ValidationErrors)
		if !ok {
			return err
		}
		resp := ErrorResponse{Errors: make([]FieldError, 0, len(verrs))}
		for _, ferr := range verrs {
			resp.Errors = append(resp.Errors, FieldError{
				Field:   strings.ToLower(ferr.Field()),
				Tag:     ferr.Tag(),
				Message: fmt.Sprintf("%s failed on %s", ferr.Field(), ferr.Tag()),
			})
		}
		return resp
	}
	return nil
}
