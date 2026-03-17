package validator

import (
	"strings"

	playvalidator "github.com/go-playground/validator/v10"
)

type Validator struct {
	validate *playvalidator.Validate
}

func New() *Validator {
	return &Validator{validate: playvalidator.New()}
}

func (v *Validator) Struct(payload any) map[string]string {
	if err := v.validate.Struct(payload); err != nil {
		if validationErrors, ok := err.(playvalidator.ValidationErrors); ok {
			errors := make(map[string]string, len(validationErrors))
			for _, item := range validationErrors {
				errors[strings.ToLower(item.Field())] = item.Tag()
			}
			return errors
		}
		return map[string]string{"request": err.Error()}
	}
	return nil
}
