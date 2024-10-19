package di

import (
	"errors"
	"reflect"
	"slices"
	"sync"
)

// Scope defines a container for registrations and resolution.
type Scope struct {
	name string

	providers     map[reflect.Type]reflect.Value
	providersLock *sync.RWMutex

	destroyers     []destroyer
	destroyersLock *sync.RWMutex
}

// NewScope creates a new [Scope] with the given name.
func NewScope(name string) *Scope {
	return &Scope{
		name: name,

		providers:     make(map[reflect.Type]reflect.Value),
		providersLock: new(sync.RWMutex),

		destroyers:     make([]destroyer, 0),
		destroyersLock: new(sync.RWMutex),
	}
}

// String returns the name of the scope.
func (s *Scope) String() string {
	return s.name
}

func (s *Scope) registerProvider(r reflect.Type, provider reflect.Value) {
	s.providersLock.Lock()
	s.providers[r] = provider
	s.providersLock.Unlock()
}

func (s *Scope) registerDestroyer(value reflect.Value, destroy reflect.Value) {
	if !destroy.IsValid() || destroy.IsNil() {
		return
	}

	s.destroyersLock.Lock()
	s.destroyers = append(s.destroyers, destroyer{value, destroy})
	s.destroyersLock.Unlock()
}

func (s *Scope) resolve(r reflect.Type, trace trace) (reflect.Value, error) {
	s.providersLock.RLock()
	provider, ok := s.providers[r]
	s.providersLock.RUnlock()

	if !ok {
		return reflect.Zero(r), newErrResolve(s, r, newErrNotRegistered(r))
	}

	if cycle := slices.Index(trace, r); 0 <= cycle {
		return reflect.Zero(r), newErrResolve(s, r, newErrCycle(append(trace[cycle:], r)))
	}

	out := provider.Call([]reflect.Value{
		reflect.ValueOf(s), reflect.ValueOf(append(trace, r)),
	})

	value := out[0]
	err, _ := out[1].Interface().(error)

	if err != nil {
		err = newErrResolve(s, r, err)
	}

	return value, err
}

func (s *Scope) invoke(function reflect.Value, trace trace) ([]reflect.Value, error) {
	f := function.Type()
	args := make([]reflect.Value, f.NumIn())

	for i := range f.NumIn() {
		arg, err := s.resolve(f.In(i), trace)
		if err != nil {
			return nil, newErrInvoke(s, function, err)
		}
		args[i] = arg
	}

	return function.Call(args), nil
}

// Register adds new registrations to the scope.
// Any of the builders in this package may be used to configure a registrable entry (see [Registrable]).
// Registrations are validated at the time of registration,
// and [ErrRegister] is returned for any errors encountered.
func (s *Scope) Register(registrables ...Registrable) error {
	errs := make([]error, len(registrables))

	for i, r := range registrables {
		if err := r.register(s); err != nil {
			errs[i] = newErrRegister(s, r, err)
		}
	}

	return errors.Join(errs...)
}

// MustRegister is like [Scope.Register] but panics on error.
func (s *Scope) MustRegister(registrables ...Registrable) {
	if err := s.Register(registrables...); err != nil {
		panic(err)
	}
}

// Destroy finalizes the scope, and returns [ErrDestroy] for any errors encountered.
// The scope must not be used after it has been destroyed.
//   - All registered destroy functions are called.
func (s *Scope) Destroy() error {
	s.destroyersLock.RLock()
	errs := make([]error, len(s.destroyers))

	for i, d := range slices.Backward(s.destroyers) {
		if err := d.Destroy(); err != nil {
			errs[i] = newErrDestroy(s, d, err)
		}
	}

	s.destroyersLock.RUnlock()
	return errors.Join(errs...)
}

// MustDestroy is like [Scope.Destroy] but panics on error.
func (s *Scope) MustDestroy() {
	if err := s.Destroy(); err != nil {
		panic(err)
	}
}

// ResolveIn resolves a value for the given type R within the given scope.
// [ErrResolve] is returned if resolution fails.
func ResolveIn[R any](s *Scope) (R, error) {
	value, err := s.resolve(reflect.TypeFor[R](), nil)
	iface, _ := value.Interface().(R)
	return iface, err
}

// MustResolveIn is like [ResolveIn] but panics on error.
func MustResolveIn[R any](s *Scope) R {
	iface, err := ResolveIn[R](s)
	if err != nil {
		panic(err)
	}
	return iface
}

// InvokeIn calls the given function after resolving any input parameters
// as dependencies with the given scope, and returns its result.
// [ErrInvoke] is returned if invocation fails.
func InvokeIn(s *Scope, function any) ([]any, error) {
	f := reflect.ValueOf(function)

	if !f.IsValid() {
		return nil, newErrNil("function")
	}
	if f.Kind() != reflect.Func {
		return nil, newErrNotFunc("function", f)
	}
	if f.IsNil() {
		return nil, newErrNil("function")
	}

	out, err := s.invoke(f, nil)
	if err != nil {
		return nil, err
	}

	ifaces := make([]any, len(out))
	for i, value := range out {
		ifaces[i] = value.Interface()
	}

	return ifaces, nil
}

// MustInvokeIn is like [InvokeIn] but panics on error.
func MustInvokeIn(s *Scope, function any) []any {
	ifaces, err := InvokeIn(s, function)
	if err != nil {
		panic(err)
	}
	return ifaces
}
