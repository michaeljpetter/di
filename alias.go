package di

import (
	"fmt"
	"reflect"
)

func alias(s *Scope, r reflect.Type, of reflect.Type, provider reflect.Type) error {
	if !of.ConvertibleTo(r) {
		return newErrNotConvertible(of, r)
	}

	s.registerProvider(r, reflect.MakeFunc(
		provider,
		func(args []reflect.Value) []reflect.Value {
			resolver := args[0].Interface().(*Scope)
			trace := args[1].Interface().(trace)

			result := []reflect.Value{reflect.Zero(r), reflect.Zero(reflect.TypeFor[error]())}

			if out, err := resolver.resolve(of, trace); err != nil {
				result[1] = reflect.ValueOf(err)
			} else {
				result[0] = out.Convert(r)
			}

			return result
		},
	))

	return nil
}

// AliasBuilder provides configuration of an [Alias].
type AliasBuilder interface {
	Registrable
}

// Alias defines a pass-through from one resolved type to another.
// The type parameter R defines the resolved type for the alias.
// The type parameter Of defines the type to be aliased, which must be convertible to type R.
//
// Alias resolves its aliased type each time it is resolved.
//   - The aliased type is resolved from the scope in which the alias is being resolved.
func Alias[R, Of any]() AliasBuilder {
	return &aliasBuilder{
		r:        reflect.TypeFor[R](),
		of:       reflect.TypeFor[Of](),
		provider: reflect.TypeFor[provider[R]](),
	}
}

type aliasBuilder struct {
	r, of, provider reflect.Type
}

func (b *aliasBuilder) String() string {
	return fmt.Sprintf("Alias[%s, %s]", typeName(b.r), typeName(b.of))
}

func (b *aliasBuilder) register(s *Scope) error {
	return alias(s, b.r, b.of, b.provider)
}
