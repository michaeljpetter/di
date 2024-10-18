package di

import (
	"fmt"
	"reflect"
	"sync"
)

func singleton(s *Scope, r reflect.Type, provider reflect.Type, create reflect.Value, destroy reflect.Value) error {
	value, err := validateCreate(r, create)
	if err != nil {
		return err
	}

	if err = validateDestroy(value, destroy); err != nil {
		return err
	}

	var once sync.Once
	result := []reflect.Value{reflect.Zero(r), reflect.Zero(reflect.TypeFor[error]())}

	s.registerProvider(r, reflect.MakeFunc(
		provider,
		func([]reflect.Value) []reflect.Value {
			once.Do(func() {
				if out, err := s.invoke(create); err != nil {
					result[1] = reflect.ValueOf(err)
				} else if 1 < len(out) && !out[1].IsNil() {
					result[1] = out[1]
				} else {
					result[0] = out[0].Convert(r)
					s.registerDestroyer(out[0], destroy)
				}
			})

			return result
		},
	))

	return nil
}

// SingletonBuilder provides configuration of a [Singleton].
type SingletonBuilder interface {
	Registrable
	// Destroy configures a destroy function for the value created by this singleton.
	// See IsValidDestroy for details.
	Destroy(destroy any) SingletonBuilder
}

// Singleton defines a one-time value creator (such as a "New" function).
// The type parameter R defines the resolved type for the created value.
// See [IsValidCreate] for details.
//
// Singleton creates a new value the first time it is resolved, and returns
// the same cached value every time thereafter.
//   - Dependencies are resolved at the time of value creation.
//   - Dependencies are resolved from the scope in which the singleton was registered.
func Singleton[R any](create any) SingletonBuilder {
	return &singletonBuilder{
		r:        reflect.TypeFor[R](),
		provider: reflect.TypeFor[provider[R]](),
		create:   reflect.ValueOf(create),
	}
}

type singletonBuilder struct {
	r, provider     reflect.Type
	create, destroy reflect.Value
}

func (b *singletonBuilder) String() string {
	return fmt.Sprintf("Singleton[%s]", typeName(b.r))
}

func (b *singletonBuilder) Destroy(destroy any) SingletonBuilder {
	b.destroy = reflect.ValueOf(destroy)
	return b
}

func (b *singletonBuilder) register(s *Scope) error {
	return singleton(s, b.r, b.provider, b.create, b.destroy)
}
