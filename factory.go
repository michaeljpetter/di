package di

import (
	"fmt"
	"reflect"
)

func factory(s *Scope, r reflect.Type, provider reflect.Type, create reflect.Value, destroy reflect.Value) error {
	v, err := validateCreate(r, create)
	if err != nil {
		return err
	}

	if err = validateDestroy(v, destroy); err != nil {
		return err
	}

	s.registerProvider(r, reflect.MakeFunc(
		provider,
		func(args []reflect.Value) []reflect.Value {
			resolver := args[0].Interface().(*Scope)
			trace := args[1].Interface().(trace)

			result := []reflect.Value{reflect.Zero(r), reflect.Zero(reflect.TypeFor[error]())}

			if out, err := resolver.invoke(create, trace); err != nil {
				result[1] = reflect.ValueOf(err)
			} else if 1 < len(out) && !out[1].IsNil() {
				result[1] = out[1]
			} else {
				result[0] = out[0].Convert(r)
				resolver.registerDestroyer(out[0], destroy)
			}

			return result
		},
	))

	return nil
}

// FactoryBuilder provides configuration of a [Factory].
type FactoryBuilder interface {
	Registrable
	// Destroy configures a destroy function for the values created by this factory.
	// See IsValidDestroy for details.
	Destroy(destroy any) FactoryBuilder
}

// Factory defines a value creator (such as a "New" function).
// The type parameter R defines the resolved type for created values.
// See [IsValidCreate] for details.
//
// Factory creates a new value each time it is resolved.
//   - Dependencies are resolved at the time of value creation.
//   - Dependencies are resolved from the scope in which the factory is being resolved.
func Factory[R any](create any) FactoryBuilder {
	return &factoryBuilder{
		r:        reflect.TypeFor[R](),
		provider: reflect.TypeFor[provider[R]](),
		create:   reflect.ValueOf(create),
	}
}

type factoryBuilder struct {
	r, provider     reflect.Type
	create, destroy reflect.Value
}

func (b *factoryBuilder) String() string {
	return fmt.Sprintf("Factory[%s]", typeName(b.r))
}

func (b *factoryBuilder) Destroy(destroy any) FactoryBuilder {
	b.destroy = reflect.ValueOf(destroy)
	return b
}

func (b *factoryBuilder) register(s *Scope) error {
	return factory(s, b.r, b.provider, b.create, b.destroy)
}
