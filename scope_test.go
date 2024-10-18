package di_test

import (
	"errors"
	"github.com/michaeljpetter/di"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestScope(t *testing.T) {
	t.Run("String", func(t *testing.T) {
		assert.Equal(t, "snappy", di.NewScope("snappy").String())
	})

	t.Run("Register", func(t *testing.T) {
		s := di.NewScope("test")
		assert.NoError(t, s.Register(di.Instance[int](3)))
	})

	t.Run("RegisterEmpty", func(t *testing.T) {
		s := di.NewScope("test")
		assert.NoError(t, s.Register())
	})

	t.Run("RegisterError", func(t *testing.T) {
		s := di.NewScope("test")
		err := s.Register(di.Instance[int]("nope"))
		assert.ErrorIs(t, err, di.ErrRegister)
		assert.ErrorIs(t, err, di.ErrNotConvertible)
	})

	t.Run("MustRegister", func(t *testing.T) {
		s := di.NewScope("test")
		assert.NotPanics(t, func() { s.Register(di.Instance[int](3)) })
	})

	t.Run("MustRegisterEmpty", func(t *testing.T) {
		s := di.NewScope("test")
		assert.NotPanics(t, func() { s.Register() })
	})

	t.Run("MustRegisterError", func(t *testing.T) {
		s := di.NewScope("test")
		assert.Panics(t, func() { s.MustRegister(di.Instance[int]("nope")) })
	})

	t.Run("Destroy", func(t *testing.T) {
		var destroyed []any
		destroy := func(v any) { destroyed = append(destroyed, v) }

		s := di.NewScope("test")
		s.MustRegister(
			di.Factory[int](rotate(3, 7)).Destroy(destroy),
			di.Factory[float64](rotate(4., 6.)).Destroy(destroy))

		di.MustResolveIn[int](s)
		di.MustResolveIn[float64](s)
		s.MustRegister(di.Instance[string]("5").Destroy(destroy))
		di.MustResolveIn[float64](s)
		di.MustResolveIn[int](s)

		assert.NoError(t, s.Destroy())

		assert.Equal(t, []any{7, 6., "5", 4., 3}, destroyed)
	})

	t.Run("DestroyError", func(t *testing.T) {
		errs := rotate(errors.New("whoops"), errors.New("floops"))

		s := di.NewScope("test")
		s.MustRegister(
			di.Factory[int](rotate(1, 2)).
				Destroy(func(int) error { return errs() }))

		di.MustResolveIn[int](s)
		di.MustResolveIn[int](s)

		err := s.Destroy()
		assert.ErrorIs(t, err, di.ErrDestroy)
		assert.ErrorContains(t, err, "whoops")
		assert.ErrorContains(t, err, "floops")
	})

	t.Run("MustDestroy", func(t *testing.T) {
		s := di.NewScope("test")
		s.MustRegister(
			di.Instance[int](7).
				Destroy(func(any) error { return nil }))

		assert.NotPanics(t, func() { s.Destroy() })
	})

	t.Run("MustDestroyError", func(t *testing.T) {
		s := di.NewScope("test")
		s.MustRegister(
			di.Instance[int](7).
				Destroy(func(any) error { return errors.New("whoops") }))

		assert.Panics(t, func() { s.MustDestroy() })
	})

	t.Run("ResolveIn", func(t *testing.T) {
		s := di.NewScope("test")
		s.MustRegister(
			di.Instance[int](55))

		value, err := di.ResolveIn[int](s)
		assert.Equal(t, 55, value)
		assert.NoError(t, err)
	})

	t.Run("ResolveInError", func(t *testing.T) {
		s := di.NewScope("test")
		s.MustRegister(
			di.Factory[int](func() (int, error) { return 77, errors.New("whoops") }))

		_, err := di.ResolveIn[int](s)
		assert.ErrorIs(t, err, di.ErrResolve)
		assert.ErrorContains(t, err, "whoops")
	})

	t.Run("ResolveInNotRegistered", func(t *testing.T) {
		s := di.NewScope("test")

		_, err := di.ResolveIn[int](s)
		assert.ErrorIs(t, err, di.ErrResolve)
		assert.ErrorIs(t, err, di.ErrNotRegistered)
	})

	t.Run("MustResolveIn", func(t *testing.T) {
		s := di.NewScope("test")
		s.MustRegister(
			di.Instance[int](55))

		assert.NotPanics(t, func() { di.MustResolveIn[int](s) })
	})

	t.Run("MustResolveInError", func(t *testing.T) {
		s := di.NewScope("test")
		s.MustRegister(
			di.Factory[int](func() (int, error) { return 77, errors.New("whoops") }))

		assert.Panics(t, func() { di.MustResolveIn[int](s) })
	})

	t.Run("MustResolveInNotRegistered", func(t *testing.T) {
		s := di.NewScope("test")

		assert.Panics(t, func() { di.MustResolveIn[int](s) })
	})

	t.Run("InvokeIn", func(t *testing.T) {
		s := di.NewScope("test")
		s.MustRegister(
			di.Instance[int](6))

		values, err := di.InvokeIn(s, func(i int) int { return i + 1 })
		assert.Equal(t, []any{7}, values)
		assert.NoError(t, err)
	})

	t.Run("InvokeInError", func(t *testing.T) {
		s := di.NewScope("test")
		s.MustRegister(
			di.Factory[int](func() (int, error) { return 77, errors.New("whoops") }))

		_, err := di.InvokeIn(s, func(i int) int { return i + 1 })
		assert.ErrorIs(t, err, di.ErrInvoke)
		assert.ErrorContains(t, err, "whoops")
	})

	t.Run("InvokeInNotRegistered", func(t *testing.T) {
		s := di.NewScope("test")

		_, err := di.InvokeIn(s, func(i int) int { return i + 1 })
		assert.ErrorIs(t, err, di.ErrInvoke)
		assert.ErrorIs(t, err, di.ErrNotRegistered)
	})

	t.Run("InvokeInNil", func(t *testing.T) {
		s := di.NewScope("test")

		_, err := di.InvokeIn(s, nil)
		assert.ErrorIs(t, err, di.ErrNil)
	})

	t.Run("InvokeInNotFunc", func(t *testing.T) {
		s := di.NewScope("test")

		_, err := di.InvokeIn(s, 109)
		assert.ErrorIs(t, err, di.ErrNotFunc)
	})

	t.Run("InvokeInNilFunc", func(t *testing.T) {
		s := di.NewScope("test")

		_, err := di.InvokeIn(s, (func())(nil))
		assert.ErrorIs(t, err, di.ErrNil)
	})

	t.Run("MustInvokeIn", func(t *testing.T) {
		s := di.NewScope("test")
		s.MustRegister(
			di.Instance[int](6))

		assert.NotPanics(t, func() { di.MustInvokeIn(s, func(i int) int { return i + 1 }) })
	})

	t.Run("MustInvokeInError", func(t *testing.T) {
		s := di.NewScope("test")
		s.MustRegister(
			di.Factory[int](func() (int, error) { return 77, errors.New("whoops") }))

		assert.Panics(t, func() { di.MustInvokeIn(s, func(i int) int { return i + 1 }) })
	})

	t.Run("MustInvokeInNotRegistered", func(t *testing.T) {
		s := di.NewScope("test")

		assert.Panics(t, func() { di.MustInvokeIn(s, func(i int) int { return i + 1 }) })
	})

	t.Run("MustInvokeInNil", func(t *testing.T) {
		s := di.NewScope("test")

		assert.Panics(t, func() { di.MustInvokeIn(s, nil) })
	})

	t.Run("MustInvokeInNotFunc", func(t *testing.T) {
		s := di.NewScope("test")

		assert.Panics(t, func() { di.MustInvokeIn(s, 109) })
	})

	t.Run("MustInvokeInNilFunc", func(t *testing.T) {
		s := di.NewScope("test")

		assert.Panics(t, func() { di.MustInvokeIn(s, (func())(nil)) })
	})
}
