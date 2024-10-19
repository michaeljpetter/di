package di

import (
	"errors"
	"fmt"
	"reflect"
)

const errSeparator = "\n └> "

// Err is the root error wrapped by all errors in the di package.
var Err = errors.New("di")

func newErr(format string, values ...any) error {
	return fmt.Errorf("%w: %s", Err, fmt.Sprintf(format, values...))
}

// ErrNil indicates the invalid use of a nil argument.
var ErrNil = fmt.Errorf("%w: must be non-nil", Err)

func newErrNil(name string) error {
	return fmt.Errorf("%w: %s", ErrNil, name)
}

// ErrNotFunc indicates the invalid use of a non-function argument.
var ErrNotFunc = fmt.Errorf("%w: must be a function", Err)

func newErrNotFunc(name string, f reflect.Value) error {
	return fmt.Errorf("%w: %s (received %s)", ErrNotFunc, name, typeName(valueType(f)))
}

// ErrInvalidFunc indicates the use of a function with an invalid signature.
var ErrInvalidFunc = fmt.Errorf("%w: invalid function", Err)

func newErrInvalidFunc(name string, f reflect.Value) error {
	return fmt.Errorf("%w: %s as %s", ErrInvalidFunc, name, typeName(valueType(f)))
}

// ErrNotConvertible indicates the use of incompatible types where convertibility is required.
var ErrNotConvertible = fmt.Errorf("%w: not convertible", Err)

func newErrNotConvertible(from reflect.Type, to reflect.Type) error {
	return fmt.Errorf("%w: %s to %s", ErrNotConvertible, typeName(from), typeName(to))
}

// ErrNotAssignable indicates the use of incompatible types where assignability is required.
var ErrNotAssignable = fmt.Errorf("%w: not assignable", Err)

func newErrNotAssignable(from reflect.Type, to reflect.Type) error {
	return fmt.Errorf("%w: %s to %s", ErrNotAssignable, typeName(from), typeName(to))
}

// ErrNotRegistered indicates that no registration was found for a resolved type.
var ErrNotRegistered = fmt.Errorf("%w: not registered", Err)

func newErrNotRegistered(r reflect.Type) error {
	return fmt.Errorf("%w: %s", ErrNotRegistered, typeName(r))
}

// ErrCycle indicates that a cycle was detected during resolution.
var ErrCycle = fmt.Errorf("%w: cycle detected", Err)

func newErrCycle(t trace) error {
	return fmt.Errorf("%w: %v", ErrCycle, t)
}

// ErrRegister indicates that an error occurred during registration, and wraps the error detail.
var ErrRegister = fmt.Errorf("%w: register", Err)

func newErrRegister(s *Scope, r Registrable, err error) error {
	return fmt.Errorf("%w: %v <- %v%s%w", ErrRegister, s, r, errSeparator, err)
}

// ErrResolve indicates that an error occurred during resolution, and wraps the error detail.
var ErrResolve = fmt.Errorf("%w: resolve", Err)

func newErrResolve(s *Scope, r reflect.Type, err error) error {
	return fmt.Errorf("%w: %v -> %v%s%w", ErrResolve, s, r, errSeparator, err)
}

// ErrInvoke indicates that an error occurred during invocation, and wraps the error detail.
var ErrInvoke = fmt.Errorf("%w: invoke", Err)

func newErrInvoke(s *Scope, f reflect.Value, err error) error {
	return fmt.Errorf("%w: %v <- %v%s%w", ErrInvoke, s, typeName(valueType(f)), errSeparator, err)
}

// ErrDestroy indicates that an error occurred during destruction, and wraps the error detail.
var ErrDestroy = fmt.Errorf("%w: destroy", Err)

func newErrDestroy(s *Scope, d destroyer, err error) error {
	return fmt.Errorf("%w: %v -> %v%s%w", ErrDestroy, s, d, errSeparator, err)
}
