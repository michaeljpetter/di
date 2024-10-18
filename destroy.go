package di

import (
	"fmt"
	"reflect"
)

func validateDestroy(v reflect.Type, destroy reflect.Value) error {
	if !destroy.IsValid() {
		return nil
	}
	if destroy.Kind() != reflect.Func {
		return newErrNotFunc("destroy", destroy)
	}
	if destroy.IsNil() {
		return nil
	}

	d := destroy.Type()

	if d.NumIn() != 1 ||
		(d.NumOut() != 0 &&
			(d.NumOut() != 1 || d.Out(0) != reflect.TypeFor[error]())) {
		return newErrInvalidFunc("destroy", destroy)
	}

	v0 := d.In(0)

	if !v.AssignableTo(v0) {
		return newErrNotAssignable(v, v0)
	}

	return nil
}

// IsValidDestroy checks whether a destroy function is valid for the given value type V.
// This check is performed internally when registering destroy functions (e.g., with [Factory]),
// thus this method does not typically need to be called explicitly.
//
// In order for destroy to be valid for value type V, it must be one of the following:
//   - an untyped nil
//   - a nil function
//   - a non-nil function of the form func(T) or func(T) error, where V is assignable to T
func IsValidDestroy[V any](destroy any) bool {
	err := validateDestroy(reflect.TypeFor[V](), reflect.ValueOf(destroy))
	return err == nil
}

type destroyer struct {
	value, destroy reflect.Value
}

func (d destroyer) String() string {
	return fmt.Sprintf("[%s] %v", typeName(valueType(d.value)), d.value.Interface())
}

func (d destroyer) Destroy() error {
	out := d.destroy.Call([]reflect.Value{d.value})
	var err error
	if 0 < len(out) {
		err, _ = out[0].Interface().(error)
	}
	return err
}
