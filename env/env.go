// Package env contains the key=value store, which is used to implement
// the environment.
//
// We need to avoid circular references, so this package will store "any"
// values rather than "Primitive" values which is actually what we interact
// with.
//
// Typically you'd create an Environment with New, but to allow scopes,
// or call-frames, you can create a nested environment via NewEnvironment.
package env

import (
	"github.com/skx/yal/config"
)

// Environment holds our state
type Environment struct {

	// parent contains the parent scope, if any.
	parent *Environment

	// values holds the actual values
	values map[string]any

	// ioconfig holds the interface to the outside world,
	// which is used for I/O
	ioconfig *config.Config
}

// Get retrieves a value from the environment.
//
// If the value isn't found in the current scope, and a parent is present,
// then that parent will be used.
func (env *Environment) Get(key string) (any, bool) {
	if v, ok := env.values[key]; ok {
		return v, ok
	}
	if env.parent == nil {
		return nil, false
	}
	return env.parent.Get(key)
}

// Items returns all the items contained within our environment.
func (env *Environment) Items() map[string]any {

	// The return value
	x := make(map[string]any)

	// If we have a parent scope then set the values from that.
	if env.parent != nil {
		for pk, pv := range env.parent.Items() {
			x[pk] = pv
		}
	}

	// Add the items in our scope after those of the parent,
	// in case we have a shadowed/more-specific value.
	for k, v := range env.values {
		x[k] = v
	}

	// all done
	return x
}

// New creates a new environment, with no parent.
func New() *Environment {
	return &Environment{
		values:   make(map[string]any),
		ioconfig: config.New(),
	}
}

// NewEnvironment creates a new environment, which will use the specified
// parent environment for values in a higher level.
func NewEnvironment(parent *Environment) *Environment {
	return &Environment{
		parent:   parent,
		values:   make(map[string]any),
		ioconfig: parent.ioconfig,
	}
}

// Set updates the contents of the current environment.
func (env *Environment) Set(key string, value any) {
	env.values[key] = value
}

// SetInDefinition sets the variable where it is defined, and returns true.
// If the value is not defined anywhere then we return false.
func (env *Environment) SetInDefinition(key string, value any) bool {

	// Is it set in this scope?
	if _, ok := env.values[key]; ok {

		// Then update and return success
		env.values[key] = value
		return true
	}
	if env.parent != nil {
		if env.parent.SetInDefinition(key, value) {
			return true
		}
	}
	return false
}

// SetIOConfig updates the configuration object which is stored
// in our environment
func (env *Environment) SetIOConfig(cfg *config.Config) {
	env.ioconfig = cfg
}

// GetIOConfig returns the configuration object which is stored in
// our environment.
func (env *Environment) GetIOConfig() *config.Config {
	return env.ioconfig
}
