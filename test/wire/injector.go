package wire

var injector *Injector

// GetInjector returns the injector.
func GetInjector() *Injector {
	return injector
}

// SetInjector sets the injector.
func SetInjector(i *Injector) {
	injector = i
}
