// Package env contains our key=value store, which is used to implement
// the environment.
//
// We need to avoid circular references, so this package will store "any"
// values rather than "Primitive" values which is actually what we interact
// with.
package env

// Environment holds our state
type Environment struct {

	// parent contains the parent scope, if any.
	parent *Environment

	// values holds the actual values
	values map[string]any
}

// New creates a new environment, with no parent.
func New() *Environment {
	return &Environment{
		values: map[string]any{},
	}
}

// NewEnvironment creates a new environment, which will use the specified
// parent environment for values in a higher scope.
func NewEnvironment(parent *Environment) *Environment {
	return &Environment{
		parent: parent,
		values: map[string]any{},
	}
}

// Get retrieves a value from the environment
func (env *Environment) Get(key string) (any, bool) {
	if v, ok := env.values[key]; ok {
		return v, ok
	}
	if env.parent == nil {
		return nil, false
	}
	return env.parent.Get(key)
}

// Set updates the environment.
func (env *Environment) Set(key string, value any) {
	env.values[key] = value
}

// SetOuter sets the variable in the parent scope, if not present in this
func (env *Environment) SetOuter(key string, value any) {
	if _, ok := env.values[key]; ok {
		env.values[key] = value
		return
	}
	if env.parent != nil {
		env.parent.SetOuter(key, value)
	}
}
