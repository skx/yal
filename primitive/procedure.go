package primitive

import "github.com/skx/yal/env"

// Procedure holds a user-defined function.
type Procedure struct {
	Args []Symbol
	Body Primitive
	Env  *env.Environment
	F    func(args []Primitive) Primitive
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
	return "procedure"
}
