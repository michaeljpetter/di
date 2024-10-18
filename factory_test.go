package di_test

import (
	"errors"
	"github.com/michaeljpetter/di"
	"github.com/stretchr/testify/assert"
	"testing"
)

func rotate[T any](values ...T) func() T {
	return func() T {
		defer func() { values = append(values[1:], values[0]) }()
		return values[0]
	}
}

func TestFactory(t *testing.T) {
	t.Run("Minimal", func(t *testing.T) {
		s := di.NewScope("test")
		s.MustRegister(
			di.Factory[int](rotate(3.4, 5.6)))

		assert.Equal(t, 3, di.MustResolveIn[int](s))
		assert.Equal(t, 5, di.MustResolveIn[int](s))

		s.MustDestroy()
	})

	t.Run("Full", func(t *testing.T) {
		var destroyed []float64

		s := di.NewScope("test")
		s.MustRegister(
			di.Factory[int](rotate(3.4, 5.6)).
				Destroy(func(v float64) { destroyed = append(destroyed, v) }))

		assert.Equal(t, 3, di.MustResolveIn[int](s))
		assert.Equal(t, 5, di.MustResolveIn[int](s))

		s.MustDestroy()

		assert.ElementsMatch(t, []float64{3.4, 5.6}, destroyed)
	})

	t.Run("Error", func(t *testing.T) {
		s := di.NewScope("test")
		s.MustRegister(
			di.Factory[int](func() (int, error) { return 77, errors.New("whoops") }))

		value, err := di.ResolveIn[int](s)
		assert.Zero(t, value)
		assert.ErrorContains(t, err, "whoops")

		s.MustDestroy()
	})

	t.Run("Dependent", func(t *testing.T) {
		type double int
		s := di.NewScope("test")
		s.MustRegister(
			di.Factory[int](rotate(7, 11)),
			di.Factory[double](func(i int) int { return 2 * i }),
		)

		assert.EqualValues(t, 14, di.MustResolveIn[double](s))
		assert.EqualValues(t, 22, di.MustResolveIn[double](s))

		s.MustDestroy()
	})

	t.Run("InvalidCreate", func(t *testing.T) {
		s := di.NewScope("test")

		err := s.Register(
			di.Factory[int](func() {}))

		assert.ErrorIs(t, err, di.ErrInvalidFunc)

		s.MustDestroy()
	})

	t.Run("InvalidDestroy", func(t *testing.T) {
		s := di.NewScope("test")

		err := s.Register(
			di.Factory[int](rotate(0)).
				Destroy(func() {}))

		assert.ErrorIs(t, err, di.ErrInvalidFunc)

		s.MustDestroy()
	})
}
