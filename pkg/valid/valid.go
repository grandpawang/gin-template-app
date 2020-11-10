package valid

import "github.com/go-playground/validator/v10"

// validate object
var validate = validator.New()

// Valid form
func Valid(data interface{}) error {
	return validate.Struct(data)
}
