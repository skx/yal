// Package eval contains the core of our lisp interpreter.
//
// We require an environment to execute with, but basically all the
// core logic is here, or in the built-in functions which are added
// by the primitives package.
package eval

import (
	"context"
	"errors"
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
}

// New constructs a new evaluator.
func New(src string) *Eval {

	// Create with a default context.
	e := &Eval{
		context: context.Background(),
	}

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

// tokenize splits the input string into tokens, via a horrific regular
// expression which I don't understand!
func (ev *Eval) tokenize(str string) {
	re := regexp.MustCompile(`[\s,]*(~@|[\[\]{}()'` + "`" +
		`~^@]|"(?:\\.|[^\\"])*"|;.*|[^\s\[\]{}('"` + "`" +
		`,;)]*)`)

	for _, match := range re.FindAllStringSubmatch(str, -1) {

		// skip empty terms, or comments (which begin with ";").
		if (match[1] == "") || (match[1][0] == ';') {
			continue
		}
		ev.toks = append(ev.toks, match[1])
	}
}

// atom converts strings into symbols, booleans, etc, as appropriate.
func (ev *Eval) atom(token string) primitive.Primitive {
	switch token {
	case "#t", "true":
		return primitive.Bool(true)
	case "#f", "false":
		return primitive.Bool(false)
	case "nil":
		return primitive.Nil{}
	}
	if token[0] == '"' {
		return primitive.String(strings.ReplaceAll(strings.Trim(token, `"`), `\"`, `"`))
	}

	// if it isn't a number then it is a symbol
	f, err := strconv.ParseFloat(token, 64)
	if err == nil {
		return primitive.Number(f)
	}
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
	case "(":
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
	case ")":
		// We shouldn't ever hit this, because we skip over
		// the closing ")" when we handle "(".
		//
		// If a program is malformed we'll see it though
		return nil, errors.New("unexpected ')'")
	default:

		// Return just a single atom/primitive.
		//
		// (i.e. A non-list, non-quote, and non-string).
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
		out = ev.eval(expr, e)

		// If this is an error then return that immediately
		switch out.(type) {
		case *primitive.Error, primitive.Error:
			return out
		}
	}
	return out
}

// eval evaluates a single expression appropriately.
func (ev *Eval) eval(exp primitive.Primitive, e *env.Environment) primitive.Primitive {

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
		// Behaviour depends on the type of the primitive/expression
		// we've been given to execute.
		//
		switch exp.(type) {

		// Numbers return themselves
		case primitive.Number:
			return exp

		// Booleans return themselves
		case primitive.Bool:
			return exp

		// Strings return themselves
		case primitive.String:
			return exp

		// Nil returns itself
		case primitive.Nil:
			return exp

		// Symbols return the value they contain
		case primitive.Symbol:
			v, ok := e.Get(string(exp.(primitive.Symbol)))

			// If it wasn't found then return a nil value
			if !ok {
				return primitive.Nil{}
			}
			// Otherwise cast it (our env. package stores "any")
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

			// (begin ..)
			case primitive.Symbol("begin"):
				var ret primitive.Primitive
				for _, x := range listExp[1:] {
					ret = ev.eval(x, e)
				}
				return ret

			// (quote
			case primitive.Symbol("quote"):
				if len(listExp) < 2 {
					return primitive.Error("arity-error: not enough arguments for (quote")
				}
				return listExp[1]

			// (eval
			case primitive.Symbol("eval"):

				if len(listExp) != 2 {
					return primitive.Error("Expected only a single argument")
				}

				switch val := listExp[1].(type) {

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
			case primitive.Symbol("define"):
				if len(listExp) < 3 {
					return primitive.Error("arity-error: not enough arguments for (define ..)")
				}
				symb, ok := listExp[1].(primitive.Symbol)
				if ok {

					val := ev.eval(listExp[2], e)
					e.Set(string(symb), val)
					return primitive.Nil{}
				}
				return primitive.Error(fmt.Sprintf("Expected a symbol, got %v", listExp[1]))

			// (set!
			case primitive.Symbol("set!"):
				if len(listExp) < 3 {
					return primitive.Error("arity-error: not enough arguments for (set! ..)")
				}
				val := ev.eval(listExp[2], e)
				e.Set(string(listExp[1].(primitive.Symbol)), val)
				return primitive.Nil{}

			// (let
			case primitive.Symbol("let"):
				if len(listExp) < 2 {
					return primitive.Error("arity-error: not enough arguments for (let ..)")
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
						return primitive.Error("arity-error: binding list had missing arguments")
					}
					// get the value
					bindingVal := ev.eval(bl[1], e)

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
					ret = ev.eval(x, newEnv)
				}
				return ret

			// (cond
			case primitive.Symbol("cond"):

				// Cast the argument to a list
				l := listExp[1].(primitive.List)

				// skip the quote
				l = l[1:]

				// Iterate over the list in pairs
				for i := 0; i < len(l); i += 2 {

					var section []primitive.Primitive
					if i > len(l)-2 {
						section = l[i:]
					} else {
						section = l[i : i+2]
					}

					// Test that worked
					if len(section) != 2 {
						return primitive.Error("expected pairs of two items")
					}

					// The two parts of this section
					test := section[0]
					eval := section[1]

					// need to eval test now.
					res := ev.eval(test, e)

					// If we got an error then we return
					// it to our caller.
					e, eok := res.(primitive.Error)
					if eok {
						return e
					}

					// Was it a success?  Then
					// goto our exit.
					//
					// This is horrid, but it short-circuits
					// the evaluation of the rest of the
					// list-pairs.
					if b, ok := res.(primitive.Bool); ok && bool(b) {
						// we'll execute this statement
						// when we return
						exp = eval
						goto repeat_eval
					}

				}

			// (if
			case primitive.Symbol("if"):
				if len(listExp) < 3 {
					return primitive.Error("arity-error: not enough arguments for (if ..)")
				}

				test := ev.eval(listExp[1], e)

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

			// (lambda
			case primitive.Symbol("lambda"):

				// ensure we have arguments
				if len(listExp) < 3 {
					return primitive.Error("wrong number of arguments")
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

				return &primitive.Procedure{
					Args: args,
					Body: listExp[2],
					Env:  e,
				}

			// Anything else is either a built-in function,
			// or a user-function.
			default:

				// Find the thing we're gonna call.
				procExp := ev.eval(listExp[0], e)

				// Is it really a procedure we can call?
				proc, ok := procExp.(*primitive.Procedure)
				if !ok {
					return primitive.Error(fmt.Sprintf("argument '%s' not a function", listExp[0].ToString()))
				}

				// build up the arguments
				args := []primitive.Primitive{}
				for _, argExp := range listExp[1:] {

					// Evaluate the arg
					evalArgExp := ev.eval(argExp, e)

					// Was it an error?  Then abort
					_, ok := evalArgExp.(primitive.Error)
					if ok {
						return primitive.Error(fmt.Sprintf("error expanding argument %v", argExp))
					}

					// Otherwise collect to invoke
					args = append(args, evalArgExp)
				}

				// Is this implemented in golang?
				if proc.F != nil {

					// Then call it.
					return proc.F(args)
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
				for _, x := range proc.Args {
					if !strings.HasPrefix(x.ToString(), "&") {
						min++
					}
				}

				//
				// Check that the arguments supplied
				// match those that are expected.
				//
				if len(args) < min {
					return primitive.Error(fmt.Sprintf("arity-error - function '%s' requires %d argument(s), %d provided", listExp[0].ToString(), min, len(args)))
				}

				// Create a new environment/scope to set the
				// parameter values within.
				e = env.NewEnvironment(proc.Env)

				// For each of the arguments that have been
				// supplied
				for i, x := range args {

					// If this is not more than the
					// proc accepts
					if i < len(proc.Args) {

						// Get the parameter name
						tmp := proc.Args[i].ToString()

						// Strip off any "&" prefix
						tmp = strings.TrimPrefix(tmp, "&")

						// Does the argument have
						// a trailing type?
						if strings.Contains(tmp, ":") {

							before, after, found := strings.Cut(tmp, ":")

							if found {
								switch after {
								case "string":
									_, ok := x.(primitive.String)
									if !ok {
										return primitive.Error(fmt.Sprintf("type-validation failed: argument %s to %s was supposed to be %s, but got %v", before, listExp[0].ToString(), after, x))
									}
								case "number":
									_, ok := x.(primitive.Number)
									if !ok {
										return primitive.Error(fmt.Sprintf("type-validation failed: argument %s to %s was supposed to be %s, but got %v", before, listExp[0].ToString(), after, x))
									}
								case "function":
									_, ok := x.(*primitive.Procedure)
									if !ok {
										return primitive.Error(fmt.Sprintf("type-validation failed: argument %s to %s was supposed to be %s, but got %v", before, listExp[0].ToString(), after, x))
									}
								case "list":
									_, ok := x.(primitive.List)
									if !ok {
										return primitive.Error(fmt.Sprintf("type-validation failed:argument %s to %s was supposed to be %s, but got %v", before, listExp[0].ToString(), after, x))
									}
								case "any":
									// nop
								default:
									return primitive.Error(fmt.Sprintf("unknown type-specification %s", after))
								}

							}

							// strip off the ":foo"
							// part.
							tmp = string(before)
						}

						// And now set the value
						e.Set(tmp, x)
					}
				}

				// Here we go round the loop again.
				//
				// Which will execute the body of the function
				// this time.
				exp = proc.Body
			}
		}
	repeat_eval:
	}
}
