package di_test

import (
	"github.com/michaeljpetter/di"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestInstance(t *testing.T) {
	type intT int

	t.Run("Minimal", func(t *testing.T) {
		s := di.NewScope("test")
		s.MustRegister(di.Instance[intT](33.7))

		assert.EqualValues(t, 33, di.MustResolveIn[intT](s))

		s.MustDestroy()
	})

	t.Run("Full", func(t *testing.T) {
		var destroyed any

		s := di.NewScope("test")
		s.MustRegister(
			di.Instance[intT](33.7).
				Destroy(func(v float64) { destroyed = v }))

		assert.EqualValues(t, 33, di.MustResolveIn[intT](s))

		s.MustDestroy()

		assert.Equal(t, 33.7, destroyed)
	})

	t.Run("InvalidValue", func(t *testing.T) {
		s := di.NewScope("test")

		err := s.Register(
			di.Instance[int](nil))

		assert.ErrorIs(t, err, di.ErrNotConvertible)

		s.MustDestroy()
	})

	t.Run("InvalidDestroy", func(t *testing.T) {
		s := di.NewScope("test")

		err := s.Register(
			di.Instance[intT](33.7).
				Destroy(func() {}))

		assert.ErrorIs(t, err, di.ErrInvalidFunc)

		s.MustDestroy()
	})
}
