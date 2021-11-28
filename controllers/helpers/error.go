package helpers

import (
	"errors"

	"github.com/gookit/validate"
)

func Vaildate(value interface{}, prefix string) error {
	v := validate.Struct(value)
	v = v.WithMessages(map[string]string{
		"uint":      prefix + ".non_integer_{field}",
		"state":     prefix + ".invalid_{field}",
		"role":      prefix + ".invalid_{field}",
		"email":     prefix + ".invalid_{field}",
		"password":  prefix + ".invalid_{field}",
		"required":  prefix + ".invalid_{field}",
		"userState": prefix + ".invalid_{field}",
		"userRole":  prefix + ".invalid_{field}",
	})

	if !v.Validate() {
		for _, errs := range v.Errors.All() {
			for _, err := range errs {
				return errors.New(err)
			}
		}
	}

	return nil
}
