package di

import (
	"reflect"
)

func validateValue(r reflect.Type, value reflect.Value) (reflect.Value, error) {
	if !value.IsValid() {
		zero := reflect.Zero(r)

		if !valueIsNil(zero) {
			return value, newErrNotConvertible(nil, r)
		}

		return zero, nil
	}

	if !value.CanConvert(r) {
		return value, newErrNotConvertible(value.Type(), r)
	}

	return value, nil
}

// IsValidValue checks whether a value is valid for the given resolved type R.
// This check is performed internally when registering values (e.g., with [Instance]),
// thus this method does not typically need to be called explicitly.
//
// In order for value to be valid for resolved type R:
//   - if value is nil, R must be a nilable type
//   - if value is not nil, it must be convertible to R
func IsValidValue[R any](value any) bool {
	_, err := validateValue(reflect.TypeFor[R](), reflect.ValueOf(value))
	return err == nil
}
