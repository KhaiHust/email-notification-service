package helper

import (
	"errors"
	"fmt"
	"github.com/go-playground/validator/v10"
	"reflect"
)

const TagValidateMessage = "validateMessage"
const TagErrorCode = "errorCode"

type CustomValidate struct {
	*validator.Validate
	Message *string
}

func (customValidate *CustomValidate) init(validate *validator.Validate) {
	customValidate.Validate = validate
}
func (customValidate *CustomValidate) Struct(current interface{}) error {
	errValidate := customValidate.Validate.Struct(current)
	if errValidate != nil {
		for _, err := range errValidate.(validator.ValidationErrors) {
			t := reflect.TypeOf(current)
			for i := 0; i < t.NumField(); i++ {
				if string(t.Field(i).Name) == err.Field() {
					errMsg := t.Field(i).Tag.Get(TagValidateMessage)
					if len(errMsg) > 0 {
						return errors.New(errMsg)
					}
				}
			}
			return fmt.Errorf("field error: %s", err.Field())
		}
	}
	return nil
}
