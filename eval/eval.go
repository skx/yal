// Package eval contains the core of our lisp interpreter.
//
// We require an environment to execute with, but basically all the
// core logic is here, or in the built-in functions which are added
// by the primitives package.
package eval

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/skx/yal/env"
	"github.com/skx/yal/primitive"
)

// ErrEOF is used to indicate when we've finished parsing
var ErrEOF = errors.New("unexpected EOF")

// Eval holds our program/state
type Eval struct {

	// toks contains the tokenized input, which we'll interpret.
	toks []string

	// offset records where in our list of tokens we're going to
	// read from next.
	offset int

	// context for handling timeout
	context context.Context

	// Recurse keeps track of how many times we've recursed
	recurse int

	// Symbols contains our (interned) symbol atom
	symbols map[string]primitive.Primitive

	// aliases contains any record of aliased functionality
	aliases map[string]string
}

// New constructs a new evaluator.
func New(src string) *Eval {

	// Create with a default context.
	e := &Eval{
		context: context.Background(),
		symbols: make(map[string]primitive.Primitive),
	}

	// Setup the default symbol-table entries

	// true
	t := primitive.Bool(true)
	e.symbols["#t"] = t
	e.symbols["true"] = t

	// false
	f := primitive.Bool(false)
	e.symbols["#f"] = f
	e.symbols["false"] = f

	// nil
	n := primitive.Nil{}
	e.symbols["nil"] = n

	// character literals - escaped characters
	e.symbols["#\\\\a"] = primitive.Character("\a")
	e.symbols["#\\\\b"] = primitive.Character("\b")
	e.symbols["#\\\\f"] = primitive.Character("\f")
	e.symbols["#\\\\n"] = primitive.Character("\n")
	e.symbols["#\\\\r"] = primitive.Character("\r")
	e.symbols["#\\\\t"] = primitive.Character("\t")

	// Record aliases here
	e.aliases = make(map[string]string)

	// tokenize our input program into a series of terms
	e.tokenize(src)

	return e
}

// SetContext allows a context to be passed to the evaluator.
//
// The context allows you to setup a timeout/deadline for the
// execution of user-supplied scripts.
func (ev *Eval) SetContext(ctx context.Context) {
	ev.context = ctx
}

// Aliased returns records of anything that has been aliased with "(alias ..)"
func (ev *Eval) Aliased() map[string]string {
	return ev.aliases
}

// tokenize splits the input string into tokens, via a horrific regular
// expression which I don't understand!
func (ev *Eval) tokenize(str string) {
	re := regexp.MustCompile(`[\s,]*(~@|[\[\]{}()'` + "`" +
		`~^@]|"(?:\\.|[^\\"])*"|;.*|[^\s\[\]{}('"` + "`" +
		`,;)]*)`)

	for _, match := range re.FindAllStringSubmatch(str, -1) {

		// skip empty terms
		if match[1] == "" {
			continue
		}

		// skip comments
		if len(match[1]) > 1 && match[1][0] == ';' {
			continue
		}

		// skip shebang
		if len(match[1]) > 2 && match[1][0] == '#' && match[1][1] == '!' {
			continue
		}

		ev.toks = append(ev.toks, match[1])
	}
}

// atom converts strings into symbols, booleans, etc, as appropriate.
func (ev *Eval) atom(token string) primitive.Primitive {

	// Lookup the contents of the symbol in our
	// symbol-table.
	//
	// This gives us interning for free.
	val, ok := ev.symbols[token]
	if ok {
		return val
	}

	// String
	if token[0] == '"' {
		return primitive.String(strings.ReplaceAll(strings.Trim(token, `"`), `\"`, `"`))
	}

	// Character
	if strings.HasPrefix(token, "#\\") {
		lit := token[2:]

		if len(lit) == 1 {

			// simple case "#\x", for example
			c := primitive.Character(lit)
			ev.symbols[token] = c
			return c
		}

		// Ensure we have an error
		return primitive.Error(fmt.Sprintf("invalid character literal: %s", lit))
	}

	// Is it a number?
	f, err := strconv.ParseFloat(token, 64)
	if err == nil {

		// The value we'll return
		n := primitive.Number(f)

		// If this is an integer then save it in our
		// interned table, for the future.
		if f == float64(int(f)) {

			ev.symbols[token] = n
		}

		return n
	}

	// OK, not something special, not a number, string, or
	// character.  It is a symbol.
	return primitive.Symbol(token)
}

// readExpression uses recursion to read a complete expression from
// our internal array of tokens - as produced by `tokenize`.
func (ev *Eval) readExpression() (primitive.Primitive, error) {

	// Have we walked off the end of the program?
	if ev.offset >= len(ev.toks) {
		return nil, ErrEOF
	}

	// Get the next token, and increase our read-position
	token := ev.toks[ev.offset]
	ev.offset++

	// We'll have different behaviour depending on what we're
	// looking at right now.
	switch token {
	case "'":
		// '... => (quote ...)
		quoted, err := ev.readExpression()
		if err != nil {
			return nil, err
		}
		return primitive.List{ev.atom("quote"), quoted}, nil

	case "`":
		// `... => (quasiquote ...)
		quoted, err := ev.readExpression()
		if err != nil {
			return nil, err
		}
		return primitive.List{ev.atom("quasiquote"), quoted}, nil

	case "~", ",":
		// ~... => (unquote ...)
		quoted, err := ev.readExpression()
		if err != nil {
			return nil, err
		}
		return primitive.List{ev.atom("unquote"), quoted}, nil

	case "~@", "`,", ",@":
		// ~@... => (splice-unquote ...)
		quoted, err := ev.readExpression()
		if err != nil {
			return nil, err
		}
		return primitive.List{ev.atom("splice-unquote"), quoted}, nil

	case "(":
		// ( .. => (list ...)

		// Are we at the end of our program?
		if ev.offset >= len(ev.toks) {
			return nil, ErrEOF
		}

		// Create a list, which we'll populate with items
		// until we reach the matching ")" statement
		list := primitive.List{}

		// Loop until we hit the closing bracket
		for ev.toks[ev.offset] != ")" {

			// Read the sub-expressions, recursively.
			expr, err := ev.readExpression()
			if err != nil {
				return nil, err
			}
			list = append(list, expr)

			// Check again we've not hit the end of the program
			if ev.offset >= len(ev.toks) {
				return nil, ErrEOF
			}
		}

		// We bump the current read-position one more here,
		// which means we skip over the closing ")" character.
		ev.offset++

		return list, nil

	case "{":
		// { .. => (hash ...)

		// Are we at the end of our program?
		if ev.offset >= len(ev.toks) {
			return nil, ErrEOF
		}

		// Create a hash, which we'll populate with items
		// until we reach the matching ")" statement
		hash := primitive.NewHash()

		// Loop until we hit the closing bracket
		for ev.toks[ev.offset] != "}" {

			// Read the sub-expressions, recursively.
			key, err := ev.readExpression()
			if err != nil {
				return nil, err
			}

			// Check again we've not hit the end of the program
			if ev.offset >= len(ev.toks) {
				return nil, ErrEOF
			}

			// Read the sub-expressions, recursively.
			val, err2 := ev.readExpression()
			if err2 != nil {
				return nil, err2
			}

			// Check again we've not hit the end of the program
			if ev.offset >= len(ev.toks) {
				return nil, ErrEOF
			}

			hash.Set(key.ToString(), val)
		}

		// We bump the current read-position one more here,
		// which means we skip over the closing "}" character.
		ev.offset++

		return hash, nil

	case ")", "}":
		// We shouldn't ever hit these, because we skip over
		// the closing characters ")" and "}" when we handle
		// the corresponding opening character.
		return nil, errors.New("unexpected '" + token + "'")

	default:

		// Return a single atom/primitive.
		return ev.atom(token), nil
	}
}

// Evaluate executes the source that was passed in the constructor,
// using the given environment for storing/retrieving state.
//
// The return value of this function is that of the last expression which
// was executed.
func (ev *Eval) Evaluate(e *env.Environment) primitive.Primitive {

	// Reset our position so we can evaluate the same program
	// multiple times
	ev.offset = 0

	// Our output/return value
	var out primitive.Primitive

	// Default to "nil" not "<nil>"
	out = primitive.Nil{}

	// loop over all input
	for {
		// Get the next expression
		expr, err := ev.readExpression()

		if err != nil {
			// End of list?
			if err == ErrEOF {
				return out
			}
			break
		}

		// Evaluate, and save the result
		out = ev.eval(expr, e, true)

		// If this is an error then return that immediately
		switch out.(type) {
		case *primitive.Error, primitive.Error:
			return out
		}
	}
	return out
}

// Does the given list start with a call to the given function?
func (ev *Eval) startsWith(l primitive.List, val string) bool {
	// list must have one entry
	if len(l) < 1 {
		return false
	}

	// entry should match
	return l[0].ToString() == val
}

// quote/quote loop
func (ev *Eval) qqLoop(xs primitive.List) primitive.List {
	var acc primitive.List

	for i := len(xs) - 1; 0 <= i; i-- {
		elt := xs[i]
		switch e := elt.(type) {
		case primitive.List:
			if ev.startsWith(e, "splice-unquote") {
				tmp := primitive.List{}
				tmp = append(tmp, primitive.Symbol("concat"))
				tmp = append(tmp, e[1])
				tmp = append(tmp, acc)
				acc = tmp
				continue
			}
		default:
		}

		tmp := primitive.List{}
		tmp = append(tmp, primitive.Symbol("cons"))
		tmp = append(tmp, ev.quasiquote(elt))
		tmp = append(tmp, acc)
		acc = tmp
	}
	return acc
}

// quasiquote handler.
func (ev *Eval) quasiquote(exp primitive.Primitive) primitive.Primitive {
	switch a := exp.(type) {
	case primitive.Symbol:

		var c primitive.List
		c = append(c, primitive.Symbol("quote"))
		c = append(c, exp)
		return c
	case primitive.List:
		if ev.startsWith(a, "unquote") {
			return a[1]
		}
		return ev.qqLoop(a)

	default:
		return exp
	}
}

// isMacro tests if a given thing is a macro
func (ev *Eval) isMacro(exp primitive.Primitive, e *env.Environment) bool {

	// If we're not being called with a list then there's nothing to do
	l, ok := exp.(primitive.List)
	if !ok {
		return false
	}

	// If the list doesn't have a size it is not a macro.
	if len(l) < 1 {
		return false
	}

	// Find the thing we're gonna call.
	procExp := ev.eval(l[0], e, false)

	// Is it really a procedure we can call?
	proc, ok2 := procExp.(*primitive.Procedure)
	if !ok2 {
		return false
	}
	return proc.Macro
}

func (ev *Eval) macroExpand(exp primitive.Primitive, e *env.Environment) primitive.Primitive {

	// is this a macro?
	for ev.isMacro(exp, e) {
		exp = ev.eval(exp, e, false)
	}
	return exp
}

// eval evaluates a single expression appropriately.
func (ev *Eval) eval(exp primitive.Primitive, e *env.Environment, expandMacro bool) primitive.Primitive {

	// Bump our recursion count
	ev.recurse++

	// Ensure that when we exit we drop back down again
	defer func() {
		ev.recurse--
	}()

	// Arbitrary limit here.
	if ev.recurse > (1024 * 8) {
		if flag.Lookup("test.v") != nil {
			return primitive.Error("hit recursion limit")
		}
	}

	// Run in a loop - even though everything will be done
	// in one-shot, except for the case of evaluating a user-function.
	//
	// User-functions could be handled recursively but there was some
	// confusion about scoping earlier..

	for {

		//
		// We've been given a context, which we'll test at every
		// iteration of our main-loop.
		//
		// This is a little slow and inefficient, but we need
		// to allow our execution to be time-limited.
		//
		select {
		case <-ev.context.Done():
			return primitive.Error("context timeout - deadline exceeded")
		default:
			// nop
		}

		//
		// Expand any macro we might be dealing with
		//
		if expandMacro {
			exp = ev.macroExpand(exp, e)
		}

		//
		// Behaviour depends on the type of the primitive/expression
		// we've been given to execute.
		//
		switch obj := exp.(type) {

		// Booleans return themselves
		case primitive.Bool:
			return exp

		// Characters return themselves
		case primitive.Character:
			return exp

		// Errors return themselves
		case primitive.Error:
			return exp

		// Hashes return themselves, but the values should be
		// evaluated - see #8.
		case primitive.Hash:
			ret := primitive.NewHash()

			for x, y := range obj.Entries {

				val := ev.eval(y, e, expandMacro)

				ret.Set(x, val)
			}
			return ret

		// Numbers return themselves
		case primitive.Number:
			return exp

		// Nil returns itself
		case primitive.Nil:
			return exp

		// Strings return themselves
		case primitive.String:
			return exp

		// Symbols return the value they contain
		case primitive.Symbol:
			// A symbol with a ":" prefix is treated as a literal.
			if strings.HasPrefix(exp.ToString(), ":") {
				return exp
			}

			// Otherwise it's looked up in the environment.
			v, ok := e.Get(string(exp.(primitive.Symbol)))

			// If it wasn't found there, return a nil value
			if !ok {
				return primitive.Nil{}
			}

			// We need to cast it (our env. package stores "any")
			return v.(primitive.Primitive)

		// Lists return the result of applying the operation
		case primitive.List:

			listExp := exp.(primitive.List)

			// If the list has no entries then we abort
			if len(listExp) == 0 {
				return listExp
			}

			// special handling for some forms, based on the
			// first token/symbol
			switch listExp[0] {

			// (alias ..)
			case primitive.Symbol("alias"):
				// We need at least one pair.
				if len(listExp) < 3 {
					return primitive.ArityError()
				}

				// The arguments are gonna be a list of pairs
				args := listExp[1:]

				if len(args)%2 != 0 {
					return primitive.Error(fmt.Sprintf("(alias ..) must have an even length of arguments, got %v", args))
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
				return primitive.Nil{}

			// (do ..)
			case primitive.Symbol("do"):
				var ret primitive.Primitive
				for _, x := range listExp[1:] {
					ret = ev.eval(x, e, expandMacro)
				}
				return ret

			// (read
			case primitive.Symbol("read"):
				if len(listExp) != 2 {
					return primitive.ArityError()
				}

				arg := listExp[1].ToString()

				// Create a new evaluator with the list
				tmp := New(arg)

				// Read an expression with it.
				//
				// Note here we just _read_ the expression,
				// we don't evaluate it.
				//
				// So we don't need an environment, etc.
				//
				out, err := tmp.readExpression()
				if err != nil {
					return primitive.Error(fmt.Sprintf("failed to read %s:%s", arg, err.Error()))
				}

				// Return it.
				return out

			// (eval
			case primitive.Symbol("eval"):

				if len(listExp) != 2 {
					return primitive.ArityError()
				}

				switch val := listExp[1].(type) {

				// Evaluate
				case primitive.List:
					// Evaluate the list
					res := ev.eval(listExp[1], e, expandMacro)

					// Create a new evaluator with
					// the result as a string
					tmp := New(res.ToString())

					// Ensure that we have a suitable
					// child-environment.
					nEnv := env.NewEnvironment(e)

					// Now evaluate it.
					return tmp.Evaluate(nEnv)

				// symbol solely so we can do env. lookup
				case primitive.Symbol:
					str, ok := e.Get(val.ToString())
					if ok {
						tmp := New(str.(primitive.Primitive).ToString())
						nEnv := env.NewEnvironment(e)
						return tmp.Evaluate(nEnv)
					}
					return primitive.Nil{}

				// string eval
				case primitive.String:
					tmp := New(string(val))
					nEnv := env.NewEnvironment(e)
					return tmp.Evaluate(nEnv)

				default:
					return primitive.Error(fmt.Sprintf("unexpected type for eval %V.", listExp[1]))
				}

			// (define
			case primitive.Symbol("define"), primitive.Symbol("def!"):
				if len(listExp) < 3 {
					return primitive.ArityError()
				}
				symb, ok := listExp[1].(primitive.Symbol)
				if ok {

					val := ev.eval(listExp[2], e, expandMacro)
					e.Set(string(symb), val)
					return primitive.Nil{}
				}
				return primitive.Error(fmt.Sprintf("Expected a symbol, got %v", listExp[1]))

			// (defmacro!
			case primitive.Symbol("defmacro!"):
				if len(listExp) < 3 {
					return primitive.ArityError()
				}

				// name of macro
				symb, ok := listExp[1].(primitive.Symbol)
				if !ok {
					return primitive.Error(fmt.Sprintf("Expected a symbol, got %v", listExp[1]))
				}

				// macro body
				val := ev.eval(listExp[2], e, expandMacro)

				mac, ok2 := val.(*primitive.Procedure)
				if !ok2 {
					return primitive.Error(fmt.Sprintf("expected a function body for (defmacro..), got %v", val))
				}

				// this is now a macro
				mac.Macro = true
				e.Set(string(symb), mac)
				return primitive.Nil{}

			// (set!
			case primitive.Symbol("set!"):
				if len(listExp) < 3 {
					return primitive.ArityError()
				}

				// Get the symbol we're gonna set
				sym, ok := listExp[1].(primitive.Symbol)
				if !ok {
					return primitive.Error(fmt.Sprintf("tried to set a non-symbol %v", listExp[1]))
				}

				// Get the value.
				val := ev.eval(listExp[2], e, expandMacro)

				// Now set, either locally or in the parent scope.
				if len(listExp) == 4 {
					e.SetOuter(string(sym), val)
				} else {
					e.Set(string(sym), val)
				}
				return primitive.Nil{}

			// (quote ..)
			case primitive.Symbol("quote"):
				if len(listExp) < 2 {
					return primitive.ArityError()
				}
				return listExp[1]

			// (quasiquote ..)
			case primitive.Symbol("quasiquote"):
				if len(listExp) < 2 {
					return primitive.ArityError()
				}
				exp = ev.quasiquote(listExp[1])
				goto repeat_eval

			// (macroexpand ..)
			case primitive.Symbol("macroexpand"):
				if len(listExp) < 2 {
					return primitive.ArityError()
				}
				return ev.macroExpand(listExp[1], e)

			// (let
			case primitive.Symbol("let"):
				if len(listExp) < 2 {
					return primitive.ArityError()
				}

				newEnv := env.NewEnvironment(e)
				bindingsList, ok := listExp[1].(primitive.List)
				if !ok {
					return primitive.Error(fmt.Sprintf("argument is not a list, got %v", listExp[1]))
				}

				for _, binding := range bindingsList {

					// ensure we got a list
					bl, ok := binding.(primitive.List)
					if !ok {
						return primitive.Error(fmt.Sprintf("binding value is not a list, got %v", binding))
					}

					if len(bl) < 2 {
						return primitive.ArityError()
					}
					// get the value
					bindingVal := ev.eval(bl[1], newEnv, expandMacro)

					// The thing to set
					set, ok2 := bl[0].(primitive.Symbol)
					if !ok2 {
						return primitive.Error(fmt.Sprintf("binding name is not a symbol, got %v", bl[0]))
					}

					// Finally set the parameter
					newEnv.Set(string(set), bindingVal)
				}

				// Now we've populated the new
				// environment with the pairs we received
				// in the setup phase we can execute
				// the body.
				var ret primitive.Primitive
				for _, x := range listExp[2:] {
					ret = ev.eval(x, newEnv, expandMacro)
				}
				return ret

			// (let*
			case primitive.Symbol("let*"):
				// let should have two entries

				if len(listExp) < 2 {
					return primitive.ArityError()
				}

				newEnv := env.NewEnvironment(e)
				bindingsList, ok := listExp[1].(primitive.List)
				if !ok {
					return primitive.Error(fmt.Sprintf("argument is not a list, got %v", listExp[1]))
				}

				// Length of binding must be %2
				if len(bindingsList)%2 != 0 {
					return primitive.Error(fmt.Sprintf("list for (len*) must have even length, got %v", bindingsList))
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
						return primitive.Error(fmt.Sprintf("binding name is not a symbol, got %v", key))
					}

					// Finally set the parameter
					newEnv.Set(string(eKey), eVal)
				}

				// Now we've populated the new
				// environment with the pairs we received
				// in the setup phase we can execute
				// the body.
				var ret primitive.Primitive
				for _, x := range listExp[2:] {
					ret = ev.eval(x, newEnv, expandMacro)
				}
				return ret

			// (env
			case primitive.Symbol("env"):

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

				return c

			// (if
			case primitive.Symbol("if"):
				if len(listExp) < 3 {
					return primitive.ArityError()
				}

				test := ev.eval(listExp[1], e, expandMacro)

				// If we got an error inside the `if` then we return it
				e, eok := test.(primitive.Error)
				if eok {
					return e
				}

				// if the test was false then we return
				// the else-section
				if b, ok := test.(primitive.Bool); (ok && !bool(b)) || primitive.IsNil(test) {
					if len(listExp) < 4 {
						return primitive.Nil{}
					}
					exp = listExp[3]
					continue
				}

				// otherwise we handle the true-section.
				exp = listExp[2]
				continue

			// (try
			case primitive.Symbol("try"):
				if len(listExp) < 3 {
					return primitive.ArityError()
				}

				// first expression is what to execute: a list
				expr := listExp[1]

				// Cast the argument to a list
				expLst, ok1 := expr.(primitive.List)
				if !ok1 {
					return primitive.Error(fmt.Sprintf("expected a list for argument, got %v", listExp[1]))
				}

				// second expression is the catch: a list
				blk := listExp[2]
				blkLst, ok2 := blk.(primitive.List)
				if !ok2 {
					return primitive.Error(fmt.Sprintf("expected a list for argument, got %v", listExp[2]))
				}
				if len(blkLst) != 3 {
					return primitive.Error(fmt.Sprintf("list should have three elements, got %v", blkLst))
				}
				if !ev.startsWith(blkLst, "catch") {
					return primitive.Error(fmt.Sprintf("catch list should begin with 'catch', got %v", blkLst))
				}

				// Evaluate the expression
				out := ev.eval(expLst, e, expandMacro)

				// Evaluating the expression didn't return an error.
				//
				// Nothing to catch, all OK
				_, ok3 := out.(primitive.Error)
				if !ok3 {
					return out
				}

				// The catch statement is blkLst[0] - we tested for that already
				// The variable to bind is blkLst[1]
				// The form to execute with that is blkLst[2]
				tmpEnv := env.NewEnvironment(e)
				tmpEnv.Set(blkLst[1].ToString(), primitive.String(out.ToString()))

				return ev.eval(blkLst[2], tmpEnv, expandMacro)

			// (lambda
			case primitive.Symbol("lambda"), primitive.Symbol("fn*"):

				// ensure we have arguments
				if len(listExp) != 3 && len(listExp) != 4 {
					return primitive.ArityError()
				}

				// ensure that our arguments are a list
				argMarkers, ok := listExp[1].(primitive.List)
				if !ok {
					return primitive.Error(fmt.Sprintf("expected a list for arguments, got %v", listExp[1]))
				}

				// Collect arguments
				args := []primitive.Symbol{}
				for _, x := range argMarkers {

					xs, ok := x.(primitive.Symbol)
					if !ok {
						return primitive.Error(fmt.Sprintf("expected a symbol for an argument, got %v", x))
					}
					args = append(args, xs)
				}

				body := listExp[2]
				help := ""

				// If there's an optional help string ..
				if len(listExp) == 4 {
					help = listExp[2].ToString()
					body = listExp[3]

				}
				// This is a procedure, which will default
				// to not being a macro.
				//
				// To make it a macro it should be set with
				// "(defmacro!..)"
				return &primitive.Procedure{
					Args:  args,
					Body:  body,
					Env:   e,
					Help:  help,
					Macro: false,
				}

			// Anything else is either a built-in function,
			// or a user-function.
			default:

				// Find the thing we're gonna call.
				procExp := ev.eval(listExp[0], e, expandMacro)

				// Is it really a procedure we can call?
				proc, ok := procExp.(*primitive.Procedure)
				if !ok {
					return primitive.Error(fmt.Sprintf("argument '%s' not a function", listExp[0].ToString()))
				}

				// build up the arguments
				args := []primitive.Primitive{}

				// Is this a macro?
				if proc.Macro {

					// Then the arguments are NOT evaluated
					args = listExp[1:]

				} else {
					// We evaluate the arguments
					for _, argExp := range listExp[1:] {

						// Evaluate the arg
						evalArgExp := ev.eval(argExp, e, expandMacro)

						// Was it an error?  Then abort
						x, ok := evalArgExp.(primitive.Error)
						if ok {
							return primitive.Error(fmt.Sprintf("error expanding argument %v for call to (%s ..): %s", argExp, listExp[0], x.ToString()))
						}

						// Otherwise append it to the list we'll supply
						args = append(args, evalArgExp)
					}
				}

				// Is this function implemented in golang?
				if proc.F != nil {

					// Then call it.
					return proc.F(e, args)
				}

				//
				// Iterate over the arguments the
				// lambda has and count those that
				// are mandatory.
				//
				// i.e. If we see this we have two args:
				//
				// (define foo (lambda (a b) ...
				//
				// However for this we accept 1+
				//
				// (define bar (lambda (a &b) ..
				//
				// We didn't do this in the golang-implemented
				// primitives as they handle argument counting
				// themselves.
				//
				min := 0

				//
				// if this is non-empty then we add all parameters here as a llist
				//
				variadic := ""

				//
				// The list of arguments to add when working in a variadic fashion
				//
				var lst primitive.List

				//
				// Count the minimum number of arguments.
				//
				// A variadic argument may be nil of course.
				//
				for _, arg := range proc.Args {
					if !strings.HasPrefix(arg.ToString(), "&") {
						min++
					}
				}

				//
				// Check that the arguments supplied match those that are expected.
				//
				// Unless variadic arguments are expected, because in that case "anything" is fine.
				//
				if len(args) < min && (variadic == "") {
					return primitive.ArityError()
				}

				// Create a new environment/scope to set the
				// parameter values within.
				e = env.NewEnvironment(proc.Env)

				// For each of the arguments that have been supplied
				for i, x := range args {

					// If this is not more than the proc accepts
					if i < len(proc.Args) {

						// Get the parameter name
						tmp := proc.Args[i].ToString()

						// Is this variadic?
						//
						// Then save the name of the argument away, after removing
						// the prefix
						//
						if strings.HasPrefix(tmp, "&") {
							tmp = strings.TrimPrefix(tmp, "&")
							variadic = tmp
						}

						// Does the argument have a trailing type?
						if strings.Contains(tmp, ":") {

							before, after, found := strings.Cut(tmp, ":")

							// Did we find it?
							if found {

								// types that are allowed
								valid := make(map[string]bool)

								// Is anything possible?
								any := false

								// Record each one
								for _, typ := range strings.Split(after, ":") {

									// Any is special
									if typ == "any" {
										any = true
									}

									// Since we're calling `type` we need to
									// do some rewriting for the function-case,
									// which has distinct types.
									if typ == "function" {
										valid["procedure(lisp)"] = true
										valid["procedure(golang)"] = true
										valid["macro"] = true
										continue
									}

									valid[typ] = true
								}

								// See if the type matches
								_, ok := valid[x.Type()]

								if !ok && !any {
									return primitive.Error(fmt.Sprintf("type-validation failed: argument %s to %s was supposed to be %s, got %s", before, listExp[0].ToString(), after, x.Type()))
								}
							}

							// strip off the ":foo" part.
							tmp = string(before)
						}

						// And now set the value
						if variadic == "" {
							e.Set(tmp, x)
						}

					}

					// Variadic arguments?  Then save this arg away
					if len(variadic) > 0 {
						lst = append(lst, x)
					}
				}

				// For variadic arguments we can't set the value as we go,
				// because we have to wait until we've collected them all.
				//
				// So set them now.
				if len(variadic) > 0 {
					e.Set(variadic, lst)
				}

				// Here we go round the evaluation loop again.
				//
				// Which will execute the body of the function this time.
				//
				// TCO.
				exp = proc.Body
			}
		}
	repeat_eval:
	}
}
