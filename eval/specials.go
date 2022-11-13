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

	// (env
	case "env":

		// create a new list
		var c primitive.List

		for key, val := range e.Items() {

			v := val.(primitive.Primitive)

			tmp := primitive.NewHash()
			tmp.Set(":name", primitive.String(key))
			tmp.Set(":value", v)

			// Is this a procedure?  If so
			// add the help-text
			proc, ok := v.(*primitive.Procedure)
			if ok {
				if len(proc.Help) > 0 {
					tmp.Set(":help", primitive.String(proc.Help))
				}
			}

			c = append(c, tmp)
		}

		return c, true

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

	// (if
	case "if":
		if len(args) < 2 {
			return primitive.ArityError(), true
		}

		test := ev.eval(args[0], e, expandMacro)

		// If we got an error inside the `if` then we return it
		er, eok := test.(primitive.Error)
		if eok {
			return er, true
		}

		// if the test was false then we return
		// the else-section
		if b, ok := test.(primitive.Bool); (ok && !bool(b)) || primitive.IsNil(test) {
			if len(args) < 3 {
				return primitive.Nil{}, true
			}
			return ev.eval(args[2], e, expandMacro), true
		}

		// otherwise we handle the true-section.
		return ev.eval(args[1], e, expandMacro), true

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

	// let
	case "let":
		if len(args) < 1 {
			return primitive.ArityError(), true
		}

		newEnv := env.NewEnvironment(e)
		bindingsList, ok := args[0].(primitive.List)
		if !ok {
			return primitive.Error(fmt.Sprintf("argument is not a list, got %v", args[0])), true
		}

		for _, binding := range bindingsList {

			// ensure we got a list
			bl, ok := binding.(primitive.List)
			if !ok {
				return primitive.Error(fmt.Sprintf("binding value is not a list, got %v", binding)), true
			}

			if len(bl) < 2 {
				return primitive.ArityError(), true
			}

			// get the value
			bindingVal := ev.eval(bl[1], newEnv, expandMacro)

			// The thing to set
			set, ok2 := bl[0].(primitive.Symbol)
			if !ok2 {
				return primitive.Error(fmt.Sprintf("binding name is not a symbol, got %v", bl[0])), true
			}

			// Finally set the parameter
			newEnv.Set(string(set), bindingVal)
		}

		// Now we've populated the new
		// environment with the pairs we received
		// in the setup phase we can execute
		// the body.
		var ret primitive.Primitive
		for _, x := range args[1:] {
			ret = ev.eval(x, newEnv, expandMacro)
		}
		return ret, true

	// (let*
	case "let*":
		// We need to have at least one argument.
		//
		// Later we'll test for more.  Because we need a multiple of two.
		if len(args) < 1 {
			return primitive.ArityError(), true
		}

		newEnv := env.NewEnvironment(e)
		bindingsList, ok := args[0].(primitive.List)
		if !ok {
			return primitive.Error(fmt.Sprintf("argument is not a list, got %v", args[0])), true
		}

		// Length of binding must be %2
		if len(bindingsList)%2 != 0 {
			return primitive.Error(fmt.Sprintf("list for (len*) must have even length, got %v", bindingsList)), true
		}

		for i := 0; i < len(bindingsList); i += 2 {

			// The key/val pair we're working with
			key := bindingsList[i]
			val := bindingsList[i+1]

			// evaluate the value - use the new environment.
			eVal := ev.eval(val, newEnv, expandMacro)

			// The thing to set
			eKey, ok := key.(primitive.Symbol)
			if !ok {
				return primitive.Error(fmt.Sprintf("binding name is not a symbol, got %v", key)), true
			}

			// Finally set the parameter
			newEnv.Set(string(eKey), eVal)
		}

		// Now we've populated the new
		// environment with the pairs we received
		// in the setup phase we can execute
		// the body.
		var ret primitive.Primitive
		for _, x := range args[1:] {
			ret = ev.eval(x, newEnv, expandMacro)
		}
		return ret, true

	// (macroexpand ..)
	case "macroexpand":
		if len(args) != 1 {
			return primitive.ArityError(), true
		}
		return ev.macroExpand(args[0], e), true

	// (quasiquote ..)
	case "quasiquote":
		if len(args) != 1 {
			return primitive.ArityError(), true
		}
		return ev.eval(ev.quasiquote(args[0]), e, expandMacro), true

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

	// (struct
	case "struct":
		if len(args) <= 1 {
			return primitive.ArityError(), true
		}

		// name of structure
		name := args[0].ToString()

		// the fields it contains
		fields := []string{}

		// convert the fields to strings
		for _, field := range args[1:] {

			f := field.ToString()

			ev.accessors[name+"."+f] = f
			fields = append(fields, f)
		}

		// save the structure as a known-thing
		ev.structs[name] = fields
		return primitive.Nil{}, true

	case "symbol":
		if len(args) != 1 {
			return primitive.ArityError(), true
		}
		return ev.atom(args[0].ToString()), true

	// (try
	case "try":
		if len(args) < 2 {
			return primitive.ArityError(), true
		}

		// first expression is what to execute: a list
		expr := args[0]

		// Cast the argument to a list
		expLst, ok1 := expr.(primitive.List)
		if !ok1 {
			return primitive.Error(fmt.Sprintf("expected a list for argument, got %v", args[0])), true
		}

		// second expression is the catch: a list
		blk := args[1]
		blkLst, ok2 := blk.(primitive.List)
		if !ok2 {
			return primitive.Error(fmt.Sprintf("expected a list for argument, got %v", args[1])), true
		}
		if len(blkLst) != 3 {
			return primitive.Error(fmt.Sprintf("list should have three elements, got %v", blkLst)), true
		}
		if !ev.startsWith(blkLst, "catch") {
			return primitive.Error(fmt.Sprintf("catch list should begin with 'catch', got %v", blkLst)), true
		}

		// Evaluate the expression
		out := ev.eval(expLst, e, expandMacro)

		// Evaluating the expression didn't return an error.
		//
		// Nothing to catch, all OK
		_, ok3 := out.(primitive.Error)
		if !ok3 {
			return out, true
		}

		// The catch statement is blkLst[0] - we tested for that already
		// The variable to bind is blkLst[1]
		// The form to execute with that is blkLst[2]
		tmpEnv := env.NewEnvironment(e)
		tmpEnv.Set(blkLst[1].ToString(), primitive.String(out.ToString()))
		return ev.eval(blkLst[2], tmpEnv, expandMacro), true
	}

	// The input was not handled as a special.
	return primitive.Nil{}, false
}
