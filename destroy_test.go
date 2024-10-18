package di_test

import (
	"github.com/michaeljpetter/di"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestIsValidDestroy(t *testing.T) {
	type interfaceT interface{ M() }
	type structT struct{ v int }
	type intT int

	noCall := func() { t.Error("should not call destroy function") }

	t.Run("ArgNoReturns", func(t *testing.T) {
		assert.True(t, di.IsValidDestroy[*int](
			func(*int) { noCall() },
		))
	})

	t.Run("ArgReturnsError", func(t *testing.T) {
		assert.True(t, di.IsValidDestroy[*structT](
			func(*structT) error { noCall(); return nil },
		))
	})

	t.Run("NoArgs", func(t *testing.T) {
		assert.False(t, di.IsValidDestroy[*int](
			func() { noCall() },
		))
	})

	t.Run("TooManyArgs", func(t *testing.T) {
		assert.False(t, di.IsValidDestroy[*structT](
			func(*structT, int) { noCall() },
		))
	})

	t.Run("TooManyReturns", func(t *testing.T) {
		assert.False(t, di.IsValidDestroy[interfaceT](
			func(interfaceT) (error, int) { noCall(); return nil, 0 },
		))
	})

	t.Run("ReturnNotError", func(t *testing.T) {
		assert.False(t, di.IsValidDestroy[*int](
			func(*int) bool { noCall(); return true },
		))
	})

	t.Run("Assignable", func(t *testing.T) {
		assert.True(t, di.IsValidDestroy[interfaceT](
			func(any) { noCall() },
		))
	})

	t.Run("NotAssignable", func(t *testing.T) {
		assert.False(t, di.IsValidDestroy[intT](
			func(int) { noCall() },
		))
	})

	t.Run("NilFunc", func(t *testing.T) {
		assert.True(t, di.IsValidDestroy[interfaceT](
			(func(interfaceT))(nil),
		))
	})

	t.Run("Nil", func(t *testing.T) {
		assert.True(t, di.IsValidDestroy[*structT](nil))
	})

	t.Run("NotFunc", func(t *testing.T) {
		assert.False(t, di.IsValidDestroy[int](123))
	})
}
