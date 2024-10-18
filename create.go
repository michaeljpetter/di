package di

import (
	"reflect"
)

func validateCreate(r reflect.Type, create reflect.Value) (reflect.Type, error) {
	if !create.IsValid() {
		return nil, newErrNil("create")
	}
	if create.Kind() != reflect.Func {
		return nil, newErrNotFunc("create", create)
	}
	if create.IsNil() {
		return nil, newErrNil("create")
	}

	c := create.Type()

	if c.NumOut() != 1 &&
		(c.NumOut() != 2 || c.Out(1) != reflect.TypeFor[error]()) {
		return nil, newErrInvalidFunc("create", create)
	}

	v := c.Out(0)

	if !v.ConvertibleTo(r) {
		return nil, newErrNotConvertible(v, r)
	}

	return v, nil
}

// IsValidCreate checks whether a create function is valid for the given resolved type R.
// This check is performed internally when registering create functions (e.g., with [Factory]),
// thus this method does not typically need to be called explicitly.
//
// In order for create to be valid for resolved type R, it must be:
//   - a non-nil function returning (T) or (T, error), where T is convertible to R
//
// There are no restrictions on its input parameters. They will be resolved as dependencies
// when create is called to produce a value.
func IsValidCreate[R any](create any) bool {
	_, err := validateCreate(reflect.TypeFor[R](), reflect.ValueOf(create))
	return err == nil
}
