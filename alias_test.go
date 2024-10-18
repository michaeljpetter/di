package di_test

import (
	"errors"
	"github.com/michaeljpetter/di"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAlias(t *testing.T) {
	t.Run("Minimal", func(t *testing.T) {
		s := di.NewScope("test")
		s.MustRegister(
			di.Factory[float64](rotate(3.4, 5.6)),
			di.Alias[int, float64]())

		assert.Equal(t, 3, di.MustResolveIn[int](s))
		assert.Equal(t, 5, di.MustResolveIn[int](s))

		s.MustDestroy()
	})

	t.Run("NotConvertible", func(t *testing.T) {
		s := di.NewScope("test")

		err := s.Register(
			di.Alias[int, string]())

		assert.ErrorIs(t, err, di.ErrNotConvertible)

		s.MustDestroy()
	})

	t.Run("NotRegistered", func(t *testing.T) {
		s := di.NewScope("test")
		s.MustRegister(
			di.Alias[int, float64]())

		value, err := di.ResolveIn[int](s)
		assert.Zero(t, value)
		assert.ErrorIs(t, err, di.ErrNotRegistered)

		s.MustDestroy()
	})

	t.Run("Error", func(t *testing.T) {
		errs := rotate(errors.New("whoops"), errors.New("floops"))

		s := di.NewScope("test")
		s.MustRegister(
			di.Factory[float64](func() (float64, error) { return 77.1, errs() }),
			di.Alias[int, float64]())

		value, err := di.ResolveIn[int](s)
		assert.Zero(t, value)
		assert.ErrorContains(t, err, "whoops")

		_, err = di.ResolveIn[int](s)
		assert.ErrorContains(t, err, "floops")

		s.MustDestroy()
	})
}
