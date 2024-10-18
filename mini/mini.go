// Package mini exports a "mini" version of the di package suitable
// for simple use cases. It provides an implicit single scope,
// and exported methods wrap any "Must" analogs in the di package
// to avoid explicit error handling. As such, all errors in this
// package are expressed by panicking.
package mini

import (
	"github.com/michaeljpetter/di"
)

var mini = di.NewScope("mini")

// Register adds new registrations to the implicit scope.
// See [di.Scope.Register].
func Register(registrables ...di.Registrable) {
	mini.MustRegister(registrables...)
}

// Resolve resolves a value in the implicit scope.
// See [di.ResolveIn].
func Resolve[R any]() R {
	return di.MustResolveIn[R](mini)
}

// Invoke calls a function in the implicit scope.
// See [di.InvokeIn].
func Invoke(function any) []any {
	return di.MustInvokeIn(mini, function)
}

// Destroy destroys the implicit scope.
// See [di.Scope.Destroy].
func Destroy() {
	mini.MustDestroy()
}

// See [di.Instance].
func Instance[R any](value any) di.InstanceBuilder {
	return di.Instance[R](value)
}

// See [di.Factory].
func Factory[R any](create any) di.FactoryBuilder {
	return di.Factory[R](create)
}

// See [di.Singleton].
func Singleton[R any](create any) di.SingletonBuilder {
	return di.Singleton[R](create)
}

// See [di.Alias].
func Alias[R, Of any]() di.AliasBuilder {
	return di.Alias[R, Of]()
}
