package di_test

import (
	"github.com/michaeljpetter/di"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestIsValidValue(t *testing.T) {
	type interfaceT interface{ M() }
	type structT struct{ v int }
	type intT int

	t.Run("OfType", func(t *testing.T) {
		assert.True(t, di.IsValidValue[int](33))
	})

	t.Run("Convertible", func(t *testing.T) {
		assert.True(t, di.IsValidValue[intT](33))
	})

	t.Run("NotConvertible", func(t *testing.T) {
		assert.False(t, di.IsValidValue[interfaceT](new(structT)))
	})

	t.Run("TypedNil", func(t *testing.T) {
		assert.True(t, di.IsValidValue[*structT]((*structT)(nil)))
	})

	t.Run("Nil", func(t *testing.T) {
		assert.True(t, di.IsValidValue[*structT](nil))
	})

	t.Run("NotNilable", func(t *testing.T) {
		assert.False(t, di.IsValidValue[structT](nil))
	})
}
