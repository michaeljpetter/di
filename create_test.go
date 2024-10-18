package di_test

import (
	"github.com/michaeljpetter/di"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestIsValidCreate(t *testing.T) {
	type interfaceT interface{ M() }
	type structT struct{ v int }
	type intT int

	noCall := func() { t.Error("should not call create function") }

	t.Run("NoArgsReturns1", func(t *testing.T) {
		assert.True(t, di.IsValidCreate[*int](
			func() *int { noCall(); return nil },
		))
	})

	t.Run("NoArgsReturns2", func(t *testing.T) {
		assert.True(t, di.IsValidCreate[*int](
			func() (*int, error) { noCall(); return nil, nil },
		))
	})

	t.Run("ArgsReturns1", func(t *testing.T) {
		assert.True(t, di.IsValidCreate[interfaceT](
			func(string, *int, *structT) interfaceT { noCall(); return nil },
		))
	})

	t.Run("ArgsReturns2", func(t *testing.T) {
		assert.True(t, di.IsValidCreate[interfaceT](
			func(intT, []string, any) (interfaceT, error) { noCall(); return nil, nil },
		))
	})

	t.Run("NoReturns", func(t *testing.T) {
		assert.False(t, di.IsValidCreate[*int](
			func() { noCall() },
		))
	})

	t.Run("TooManyReturns", func(t *testing.T) {
		assert.False(t, di.IsValidCreate[*int](
			func() (*int, error, string) { noCall(); return nil, nil, "" },
		))
	})

	t.Run("SecondReturnNotError", func(t *testing.T) {
		assert.False(t, di.IsValidCreate[*int](
			func() (*int, string) { noCall(); return nil, "" },
		))
	})

	t.Run("Returns1Convertible", func(t *testing.T) {
		assert.True(t, di.IsValidCreate[intT](
			func() int64 { noCall(); return 0 },
		))
	})

	t.Run("Returns2Convertible", func(t *testing.T) {
		assert.True(t, di.IsValidCreate[int32](
			func() (float32, error) { noCall(); return 0, nil },
		))
	})

	t.Run("Returns1NotConvertible", func(t *testing.T) {
		assert.False(t, di.IsValidCreate[interfaceT](
			func() *structT { noCall(); return nil },
		))
	})

	t.Run("Returns2NotConvertible", func(t *testing.T) {
		assert.False(t, di.IsValidCreate[*int](
			func() (int, error) { noCall(); return 0, nil },
		))
	})

	t.Run("NilFunc", func(t *testing.T) {
		assert.False(t, di.IsValidCreate[interfaceT](
			(func() interfaceT)(nil),
		))
	})

	t.Run("Nil", func(t *testing.T) {
		assert.False(t, di.IsValidCreate[*structT](nil))
	})

	t.Run("NotFunc", func(t *testing.T) {
		assert.False(t, di.IsValidCreate[int](123))
	})
}
