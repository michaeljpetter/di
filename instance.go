package di

import (
	"fmt"
	"reflect"
)

func instance(s *Scope, r reflect.Type, provider reflect.Type, value reflect.Value, destroy reflect.Value) error {
	value, err := validateValue(r, value)
	if err != nil {
		return err
	}

	if err = validateDestroy(value.Type(), destroy); err != nil {
		return err
	}

	result := []reflect.Value{value.Convert(r), reflect.Zero(reflect.TypeFor[error]())}

	s.registerProvider(r, reflect.MakeFunc(
		provider,
		func([]reflect.Value) []reflect.Value {
			return result
		},
	))

	s.registerDestroyer(value, destroy)

	return nil
}

// InstanceBuilder provides configuration of an [Instance].
type InstanceBuilder interface {
	Registrable
	// Destroy configures a destroy function for the value associated with this instance.
	// See IsValidDestroy for details.
	Destroy(destroy any) InstanceBuilder
}

// Instance defines an externally created value.
// The type parameter R defines the resolved type for the value.
//
// Instance returns its fixed value each time it is resolved.
func Instance[R any](value any) InstanceBuilder {
	return &instanceBuilder{
		r:        reflect.TypeFor[R](),
		provider: reflect.TypeFor[provider[R]](),
		value:    reflect.ValueOf(value),
	}
}

type instanceBuilder struct {
	r, provider    reflect.Type
	value, destroy reflect.Value
}

func (b *instanceBuilder) String() string {
	return fmt.Sprintf("Instance[%s]", typeName(b.r))
}

func (b *instanceBuilder) Destroy(destroy any) InstanceBuilder {
	b.destroy = reflect.ValueOf(destroy)
	return b
}

func (b *instanceBuilder) register(s *Scope) error {
	return instance(s, b.r, b.provider, b.value, b.destroy)
}
