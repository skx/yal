package primitive

import "github.com/skx/yal/env"

// Procedure holds a user-defined function.
type Procedure struct {

	// Arguments to this procedure
	Args []Symbol

	// Body is the body to execute, in the case where F is nil.
	//
	// In this case the primitive is written in 100% pure lisp.
	Body Primitive

	// Env contains the environment within which this procedure is executed.
	Env *env.Environment

	// F contains a pointer to the golang implementation of this procedure,
	// if it is a native one.
	F func(args []Primitive) Primitive
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
	return "(lambda " + args.ToString() + " " + p.Body.ToString() + ")"
}

// Type returns the type of this primitive object.
func (p *Procedure) Type() string {
	if p.F != nil {
		return "procedure(golang)"
	}
	return "procedure(lisp)"
}
