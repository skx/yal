// specials.go - Implementation of the special forms.

package eval

import (
	"fmt"

	"github.com/skx/yal/env"
	"github.com/skx/yal/primitive"
)

// evalSpecialForm is invoked to execute one of our special forms.
//
// This is done to centralize the code, and also ensure that eval doesn't
// get too dense.
//
// The return value from this function is "XX, BOOL".  If the boolean result
// is true this function handled the call, otherwise it did not.
//
// This is required because special forms take precedence over other calls.
func (ev *Eval) evalSpecialForm(name string, args []primitive.Primitive, e *env.Environment, expandMacro bool) (primitive.Primitive, bool) {

	switch name {
	case "alias":
		// We need at least one pair.
		if len(args) < 2 {
			return primitive.ArityError(), true
		}

		if len(args)%2 != 0 {
			return primitive.Error(fmt.Sprintf("(alias ..) must have an even length of arguments, got %v", args)), true
		}

		for i := 0; i < len(args); i += 2 {

			// The key/val pair we're working with
			new := args[i]
			orig := args[i+1]

			old, ok := e.Get(orig.ToString())
			if ok {
				e.Set(new.ToString(), old)

				ev.aliases[new.ToString()] = orig.ToString()
			}
		}
		return primitive.Nil{}, true

	// (define
	case "define", "def!":
		if len(args) < 2 {
			return primitive.ArityError(), true
		}
		symb, ok := args[0].(primitive.Symbol)
		if ok {
			val := ev.eval(args[1], e, expandMacro)
			e.Set(string(symb), val)
			return primitive.Nil{}, true
		}
		return primitive.Error(fmt.Sprintf("Expected a symbol, got %v", args[0])), true

	// (defmacro!
	case "defmacro!":
		if len(args) < 2 {
			return primitive.ArityError(), true
		}

		// name of macro
		symb, ok := args[0].(primitive.Symbol)
		if !ok {
			return primitive.Error(fmt.Sprintf("Expected a symbol, got %v", args[0])), true
		}

		// macro body
		val := ev.eval(args[1], e, expandMacro)

		mac, ok2 := val.(*primitive.Procedure)
		if !ok2 {
			return primitive.Error(fmt.Sprintf("expected a function body for (defmacro..), got %v", val)), true
		}

		// this is now a macro
		mac.Macro = true
		e.Set(string(symb), mac)
		return primitive.Nil{}, true

	case "do":
		var ret primitive.Primitive
		for _, x := range args {
			ret = ev.eval(x, e, expandMacro)
		}
		return ret, true

	case "eval":
		if len(args) != 1 {
			return primitive.ArityError(), true
		}

		switch val := args[0].(type) {

		// Evaluate
		case primitive.List:
			// Evaluate the list
			res := ev.eval(args[0], e, expandMacro)

			// Create a new evaluator with
			// the result as a string
			tmp := New(res.ToString())

			// Ensure that we have a suitable
			// child-environment.
			nEnv := env.NewEnvironment(e)

			// Now evaluate it.
			return tmp.Evaluate(nEnv), true

		// symbol solely so we can do env. lookup
		case primitive.Symbol:
			str, ok := e.Get(val.ToString())
			if ok {
				tmp := New(str.(primitive.Primitive).ToString())
				nEnv := env.NewEnvironment(e)
				return tmp.Evaluate(nEnv), true
			}
			return primitive.Nil{}, true

		// string eval
		case primitive.String:
			tmp := New(string(val))
			nEnv := env.NewEnvironment(e)
			return tmp.Evaluate(nEnv), true

		default:
			return primitive.Error(fmt.Sprintf("unexpected type for eval %V.", args[0])), true
		}

	// (lambda
	case "lambda", "fn*":
		// ensure we have arguments
		if len(args) != 2 && len(args) != 3 {
			return primitive.ArityError(), true
		}

		// ensure that our arguments are a list
		argMarkers, ok := args[0].(primitive.List)
		if !ok {
			return primitive.Error(fmt.Sprintf("expected a list for arguments, got %v", args[0])), true
		}

		// Collect arguments
		arguments := []primitive.Symbol{}
		for _, x := range argMarkers {

			xs, ok := x.(primitive.Symbol)
			if !ok {
				return primitive.Error(fmt.Sprintf("expected a symbol for an argument, got %v", x)), true
			}
			arguments = append(arguments, xs)
		}

		body := args[1]
		help := ""

		// If there's an optional help string ..
		if len(args) == 3 {
			help = args[1].ToString()
			body = args[2]
		}

		// This is a procedure, which will default
		// to not being a macro.
		//
		// To make it a macro it should be set with
		// "(defmacro!..)"
		return &primitive.Procedure{
			Args:  arguments,
			Body:  body,
			Env:   e,
			Help:  help,
			Macro: false,
		}, true

	// (quote ..)
	case "quote":
		if len(args) != 1 {
			return primitive.ArityError(), true
		}
		return args[0], true

	case "read":
		if len(args) != 1 {
			return primitive.ArityError(), true
		}

		arg := args[0].ToString()

		// Create a new evaluator with the list
		tmp := New(arg)

		// Read an expression with it.
		//
		// Note here we just _read_ the expression,
		// we don't evaluate it.
		//
		// So we don't need an environment, etc.
		//
		out, err := tmp.readExpression(e)
		if err != nil {
			return primitive.Error(fmt.Sprintf("failed to read %s:%s", arg, err.Error())), true
		}

		// Return it.
		return out, true

	// (set!
	case "set!":
		if len(args) < 2 {
			return primitive.ArityError(), true
		}

		// Get the symbol we're gonna set
		sym, ok := args[0].(primitive.Symbol)
		if !ok {
			return primitive.Error(fmt.Sprintf("tried to set a non-symbol %v", args[0])), true
		}

		// Get the value.
		val := ev.eval(args[1], e, expandMacro)

		// Now set, either locally or in the parent scope.
		if len(args) == 3 {
			e.SetOuter(string(sym), val)
		} else {
			e.Set(string(sym), val)
		}
		return primitive.Nil{}, true

	}
	return primitive.Nil{}, false
}
