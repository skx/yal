package primitive

import "github.com/skx/yal/env"

// GolangPrimitiveFn is the type which represents a function signature for
// a lisp-usable function implemented in golang.
type GolangPrimitiveFn func(e *env.Environment, args []Primitive) Primitive

// Procedure holds a user-defined function.
//
// This structure is used to hold both the built-in functions, implemented in
// golang and those which are written in lisp - either as functions or macros.
type Procedure struct {

	// Arguments to this procedure.
	Args []Symbol

	// Body is the body to execute, in the case where F is nil.
	Body Primitive

	// Env contains the environment within which this procedure is executed.
	Env *env.Environment

	// F contains a pointer to the golang implementation of this procedure,
	// if it is a native one.
	F GolangPrimitiveFn

	// Help contains some function-specific help text, ideally with
	// an example usage of the function.
	Help string

	// Macro is true is this function should have arguments passed literally, and
	// not evaluated.
	Macro bool
}

// ToInterface converts this object to a golang value
func (p *Procedure) ToInterface() any {
	return p.ToString()
}

// ToString converts this object to a string.
func (p *Procedure) ToString() string {
	if p.F != nil {
		return "#built-in-function"
	}
	args := List{}
	for _, x := range p.Args {
		args = append(args, x)
	}

	// might be a macro
	first := "lambda"
	if p.Macro {
		first = "macro"
	}

	return "(" + first + " " + args.ToString() + " " + p.Body.ToString() + ")"
}

// Type returns the type of this primitive object.
func (p *Procedure) Type() string {
	if p.Macro {
		return "macro"
	}
	if p.F != nil {
		return "procedure(golang)"
	}
	return "procedure(lisp)"
}
