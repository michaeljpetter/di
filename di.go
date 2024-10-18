// Package di implements dependency injection.
//
// A [Scope] serves as the container for registrations and resolution.
// Begin by creating a new scope:
//
//	root := di.NewScope("root")
//
// Then, add registrations to the scope via any of the builders
// available in this package, such as [Singleton]:
//
//	root.Register(
//		Singleton[MyInterface](mypackage.NewService).
//			Destroy(mypackage.Service.Shutdown))
//
// The "create" function, mypackage.NewService, may have further dependencies in its
// own arguments, which will be resolved by the scope when MyInterface is
// resolved:
//
//	my, err := di.ResolveIn[MyInterface](root)
package di

import (
	"reflect"
)

type provider[R any] func(*Scope) (R, error)

// Registrable is the base interface implemented by all
// registration builders in the di package. The following builders are available:
//   - [Instance]
//   - [Factory]
//   - [Singleton]
//   - [Alias]
type Registrable interface {
	register(*Scope) error
}

func typeName(t reflect.Type) string {
	if t == nil {
		return "<nil>"
	}
	if t == reflect.TypeFor[any]() {
		return "any"
	}
	return t.String()
}

func valueType(value reflect.Value) reflect.Type {
	if value.IsValid() {
		return value.Type()
	}
	return nil
}

func valueIsNil(value reflect.Value) bool {
	defer func() { recover() }()
	return value.IsNil()
}
