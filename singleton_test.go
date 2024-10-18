package di_test

import (
	"errors"
	"github.com/michaeljpetter/di"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSingleton(t *testing.T) {
	t.Run("Minimal", func(t *testing.T) {
		s := di.NewScope("test")
		s.MustRegister(
			di.Singleton[int](rotate(3.4, 5.6)))

		for range 3 {
			assert.Equal(t, 3, di.MustResolveIn[int](s))
		}

		s.MustDestroy()
	})

	t.Run("Full", func(t *testing.T) {
		var destroyed []float64

		s := di.NewScope("test")
		s.MustRegister(
			di.Singleton[int](rotate(3.4, 5.6)).
				Destroy(func(v float64) { destroyed = append(destroyed, v) }))

		for range 3 {
			assert.Equal(t, 3, di.MustResolveIn[int](s))
		}

		s.MustDestroy()

		assert.ElementsMatch(t, []float64{3.4}, destroyed)
	})

	t.Run("Error", func(t *testing.T) {
		errs := rotate(errors.New("whoops"), errors.New("floops"))

		s := di.NewScope("test")
		s.MustRegister(
			di.Singleton[int](func() (int, error) { return 77, errs() }))

		for range 3 {
			value, err := di.ResolveIn[int](s)
			assert.Zero(t, value)
			assert.ErrorContains(t, err, "whoops")
		}

		s.MustDestroy()
	})

	t.Run("Dependent", func(t *testing.T) {
		type double int
		s := di.NewScope("test")
		s.MustRegister(
			di.Factory[int](rotate(7, 11, 13)),
			di.Singleton[double](func(i int) int { return 2 * i }),
		)

		for range 3 {
			assert.EqualValues(t, 14, di.MustResolveIn[double](s))
		}

		s.MustDestroy()
	})

	t.Run("InvalidCreate", func(t *testing.T) {
		s := di.NewScope("test")

		err := s.Register(
			di.Singleton[int](func() {}))

		assert.ErrorIs(t, err, di.ErrInvalidFunc)

		s.MustDestroy()
	})

	t.Run("InvalidDestroy", func(t *testing.T) {
		s := di.NewScope("test")

		err := s.Register(
			di.Singleton[int](rotate(0)).
				Destroy(func() {}))

		assert.ErrorIs(t, err, di.ErrInvalidFunc)

		s.MustDestroy()
	})
}
