package common

import "github.com/go-playground/validator/v10"

func APIVersionValidator() validator.Func {
	return func(fl validator.FieldLevel) bool {
		version, ok := fl.Field().Interface().(string)
		if !ok {
			return false
		}
		return version == "v0" || version == "v1" || version == "latest"
	}
}
