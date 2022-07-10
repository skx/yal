// Package eval contains the core of our lisp interpreter.
//
// We require an environment to execute with, but basically all the
// core logic is here, or in the built-in functions which are added
// by the primitives package.
package eval

import (
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
}

// New constructs a new evaluator.
func New(src string) *Eval {
	e := &Eval{}

	// tokenize our input program into a series of terms
	e.tokenize(src)

	return e
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
	case "#t":
		return primitive.Bool(true)
	case "#f":
		return primitive.Bool(false)
	case "nil":
		return primitive.Nil{}
	}
	if token[0] == '"' {
		return primitive.String(strings.ReplaceAll(strings.Trim(token, `"`), `\"`, `"`))
	}
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

		// Now rewrite our input a little if it is a (define ..) list
		if len(list) > 0 && list[0] == primitive.Symbol("define") {
			// (define (f ...) (...)) => (define f (lambda (...) (...)))
			if argsList, ok := list[1].(primitive.List); ok {
				return primitive.List{ev.atom("define"), argsList[0], primitive.List{ev.atom("lambda"), argsList[1:], list[2]}}, nil
			}
		}

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
	// Out value
	var out primitive.Primitive

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
			v, _ := e.Get(string(exp.(primitive.Symbol)))

			// If it wasn't found then return a nil value
			if v == nil {
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
				return listExp[1]

			// (define
			case primitive.Symbol("define"):
				val := ev.eval(listExp[2], e)
				e.Set(string(listExp[1].(primitive.Symbol)), val)
				return primitive.Nil{}

			// (set!
			case primitive.Symbol("set!"):
				val := ev.eval(listExp[2], e)
				e.Set(string(listExp[1].(primitive.Symbol)), val)
				return primitive.Nil{}

			// (let
			case primitive.Symbol("let"):
				newEnv := env.NewEnvironment(e)
				bindingsList := listExp[1].(primitive.List)
				for _, binding := range bindingsList {
					bindingVal := ev.eval(binding.(primitive.List)[1], e)
					newEnv.Set(string(binding.(primitive.List)[0].(primitive.Symbol)), bindingVal)
				}
				var ret primitive.Primitive
				for _, x := range listExp[2:] {
					ret = ev.eval(x, newEnv)
				}
				return ret

			// (if
			case primitive.Symbol("if"):
				test := ev.eval(listExp[1], e)
				if b, ok := test.(primitive.Bool); (ok && !bool(b)) || primitive.IsNil(test) {
					if len(listExp) < 4 {
						return primitive.Nil{}
					}
					exp = listExp[3]
					continue
				}
				exp = listExp[2]
				continue

			// (lambda
			case primitive.Symbol("lambda"):
				args := []primitive.Symbol{}
				for _, x := range listExp[1].(primitive.List) {
					args = append(args, x.(primitive.Symbol))
				}

				return &primitive.Procedure{
					Args: args,
					Body: listExp[2],
					Env:  e,
				}

			// Anything else is either a built-in function,
			// or a user-function.
			default:
				procExp := ev.eval(listExp[0], e)
				proc, ok := procExp.(*primitive.Procedure)
				if !ok {
					return primitive.Error(fmt.Sprintf("argument '%s' not a function", procExp))
				}

				// build up the arguments to pass to the function
				args := []primitive.Primitive{}
				for _, argExp := range listExp[1:] {
					evalArgExp := ev.eval(argExp, e)
					args = append(args, evalArgExp)
				}

				// Is this implemented in golang?
				if proc.F != nil {
					return proc.F(args)
				}

				// If not then it's a user-function,
				// create a new environment/scope to set the
				// parameter values, and evaluate the body.
				e = env.NewEnvironment(proc.Env)
				for i, x := range proc.Args {
					e.Set(string(x), args[i])
				}

				// Here we go round the loop again.
				exp = proc.Body
			}
		}
	}
}
