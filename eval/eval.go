// Package eval contains the core of our lisp interpreter.
//
// We require an environment to execute with, but basically all the
// core logic is here, or in the built-in functions which are added
// by the primitives package.
package eval

import (
	"bufio"
	"context"
	"errors"
	"flag"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/skx/yal/env"
	"github.com/skx/yal/primitive"
)

// ErrEOF is used to indicate when we've finished parsing
var ErrEOF = errors.New("unexpected EOF")

// ErrTimeout is used to say that we've timed out
var ErrTimeout = errors.New("context timeout - deadline exceeded")

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

	// structs contains a list of known structures.
	//
	// The key is the name of the structure, and the value is an
	// array of the fields that structure possesses.
	structs map[string][]string

	// accessors contains struct field lookups
	//
	// The key is the name of the fake method, the value the name of
	// the field to get/set
	accessors map[string]string

	// STDIN is an input-reader used for the (read) function, when
	// called with no arguments.
	STDIN *bufio.Reader

	// STDOUT is the writer which we should use for "(print)", but we
	// currently don't.
	STDOUT *bufio.Writer
}

// New constructs a new lisp interpreter.
func New(src string) *Eval {

	// Create with a default context.
	e := &Eval{

		// aliases holds function aliases
		aliases: make(map[string]string),

		// context used for timeout-testing
		context: context.Background(),

		// symbols is an interning cache
		symbols: make(map[string]primitive.Primitive),

		// structs contains the names and expected field-names
		// of user-defined structures.
		structs: make(map[string][]string),

		// accessors contains the names of generated get/set
		// functions for field access within structs
		accessors: make(map[string]string),
	}

	// Setup default input/output streams
	e.STDIN = bufio.NewReader(os.Stdin)
	e.STDOUT = bufio.NewWriter(os.Stdout)

	// Setup the default symbol-table (interned) entries.

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

	// tokenize our input program into a series of terms
	e.tokenize(src)

	return e
}

// Aliased returns records of anything that has been aliased with "(alias ..)"
func (ev *Eval) Aliased() map[string]string {
	return ev.aliases
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
		expr, err := ev.readExpression(e)

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

// Execute will load the new code in the given src, and execute it
// using the specified environment.
//
// This allows a single interpreter to be reused to execute multiple
// expressions, with persistent state.
func (ev *Eval) Execute(e *env.Environment, src string) primitive.Primitive {

	// Reset our source
	ev.tokenize(src)

	// Now execute that source
	return (ev.Evaluate(e))
}

// SetContext allows a context to be passed to the evaluator.
//
// The context allows you to setup a timeout/deadline for the
// execution of user-supplied scripts.
func (ev *Eval) SetContext(ctx context.Context) {
	ev.context = ctx
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

	// See if this is a number with a hex/binary prefix
	based := strings.ToLower(token)
	if strings.HasPrefix(based, "0x") || strings.HasPrefix(based, "0b") {

		// If so then parse as an integer
		n, err := strconv.ParseInt(based, 0, 64)
		if err == nil {

			// Assuming it worked save it in our interned
			// table and return it.
			num := primitive.Number(n)
			ev.symbols[token] = num
			return num
		}
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

// eval evaluates a single expression appropriately.
//
// We have special cases for the simple values, for example numbers, strings,
// and similar primitive types just return themselves.
//
// Otherwise we have two special cases to handle:
//
// Symbols return the appropriate value from the environment, and
// lists involve invoking functions (or our special built-in forms).
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
			return primitive.Error(ErrTimeout.Error())
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
		// Simple types return themselves literally.
		//
		// For example a string returns itself, as
		// does a number, a boolean, or a hash.
		//
		// Lists are the things that are complex in
		// lisp, as they represent function calls.
		//
		if exp.IsSimpleType() {
			return exp
		}

		//
		// After simple types we have to deal with symbols, and lists.
		//
		sym, isSymbol := exp.(primitive.Symbol)
		if isSymbol {

			// A symbol with a ":" prefix is treated as a literal.
			if strings.HasPrefix(sym.ToString(), ":") {
				return exp
			}

			// Otherwise it's looked up in the environment.
			v, ok := e.Get(string(sym))

			// If it wasn't found there, return a nil value
			if !ok {
				return primitive.Nil{}
			}

			// We need to cast it (our env. package stores "any")
			return v.(primitive.Primitive)
		}

		//
		// Now we're only dealing with lists
		//
		listExp, listOk := exp.(primitive.List)

		//
		// But just in case we're not ..
		//
		if !listOk {
			return primitive.Error(fmt.Sprintf("argument not a list for a function call: %v", exp))
		}

		//
		// Is this an empty list?  Then just return it
		//
		if len(listExp) == 0 {
			return listExp
		}

		//
		// Is the first term a Symbol?
		//
		// If so then we're evaluating a special form,
		// these are implemented in "specials.go"
		//
		sym, isSymbol = listExp[0].(primitive.Symbol)
		if isSymbol {

			//
			// If that was handled then return the result.
			//
			// Otherwise we'll keep going and we'll treat this
			// as a usual function-call.
			//
			res, ok := ev.evalSpecialForm(sym.ToString(), listExp[1:], e, expandMacro)
			if ok {
				return res
			}
		}

		//
		// If we reached here then we've got a list as input,
		// and that is not a list that represents a call to a
		// built-in "special".
		//
		// So we need to work out what it is:
		//
		//  1. A syntethic method, relate to struct.
		//
		//  2. A golang-implemented primitive.
		//
		//  3. A user-defined function.
		//

		// The thing we'll call
		thing := listExp[0]

		// args supplied to this call
		listArgs := listExp[1:]

		// Is this a structure field access?
		access, okA := ev.accessors[thing.ToString()]
		if okA {

			// We have a single argument for the get-method
			// and two for the set-method
			if len(listArgs) != 1 && len(listArgs) != 2 {
				return primitive.ArityError()
			}

			// Get the first argument and ensure it is a hash
			obj := ev.eval(listArgs[0], e, expandMacro)
			hsh, okH := obj.(primitive.Hash)
			if !okH {
				return primitive.Error(fmt.Sprintf("expected a hash, got %v", obj))
			}

			// One argument?  Read the value
			if len(listArgs) == 1 {
				return hsh.Get(access)
			}

			// Two arguments?  Set the value, and return it
			val := ev.eval(listArgs[1], e, expandMacro)
			hsh.Set(access, val)
			return val

		}

		// Is this a structure creation?
		fields, ok := ev.structs[thing.ToString()]
		if ok {

			// ensure that we have some fields that
			// match those we expect.
			if len(listArgs) > len(fields) {
				return primitive.ArityError()
			}

			// Create a hash to store the state
			hash := primitive.NewHash()

			// However mark this as a "struct",
			// rather than a hash.
			hash.SetStruct(thing.ToString())

			// Set the fields, ensuring we evaluate them
			//
			// If some fields are unspecified they become nil.
			for i, name := range fields {
				if i < len(listArgs) {
					hash.Set(name, ev.eval(listArgs[i], e, expandMacro))
				} else {
					hash.Set(name, primitive.Nil{})
				}
			}
			return hash
		}

		// Is this a type-check on a struct?
		if strings.HasSuffix(thing.ToString(), "?") {

			// Get the thing that is being tested.
			typeName := strings.TrimSuffix(thing.ToString(), "?")

			// Does that represent a known-type?
			_, ok2 := ev.structs[typeName]
			if ok2 {

				// OK now we're sure we're not colliding
				// with another function test the argument
				// count.
				if len(listExp) != 2 {
					return primitive.ArityError()
				}

				// OK a type-check on a known struct
				//
				// Note we evaluate the object, because it
				// was probably a symbol, or return object
				// of some kind.
				obj := ev.eval(listExp[1], e, expandMacro)

				// is it a hash?
				hsh, ok2 := obj.(primitive.Hash)
				if !ok2 {
					// nope - if it isn't a hash
					// then it can't be a struct.
					return primitive.Bool(false)
				}

				// is the struct-type the same as the type name?
				if hsh.GetStruct() == typeName {
					return primitive.Bool(true)
				}
				return primitive.Bool(false)
			}

			// just a method call with a trailing "?".
			//
			// could be "string?", "contains?", etc,
			// so we fall-through and keep processing as
			// per usual.
		}

		// Find the thing we're gonna call.
		procExp := ev.eval(thing, e, expandMacro)

		// Is it really a procedure we can call?
		proc, ok := procExp.(*primitive.Procedure)
		if !ok {
			return primitive.Error(fmt.Sprintf("argument '%s' not a function", thing.ToString()))
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
		// if this is non-empty then we add all parameters here as a list
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
		// Variadic arguments would add _extra_ arguments, so this check
		// is still safe for those.
		//
		if len(args) < min {
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

					//
					// Type-check the supplied argument
					//
					argName, argTypes, found := strings.Cut(tmp, ":")

					//
					// Type-check
					//
					if found {
						err := ev.typeCheck(argTypes, x.Type())

						if err != nil {
							return primitive.TypeError(fmt.Sprintf("argument %s to %s was supposed to be %s, got %s", argName, thing.ToString(), argTypes, x.Type()))
						}
					}

					// strip off the ":foo" part.
					tmp = string(argName)
				}

				// And now set the value
				if variadic == "" {
					e.Set(tmp, x)
				}

			}

			// Variadic arguments?  Then save this arg away to
			// our temporary list, and set it.
			if len(variadic) > 0 {
				lst = append(lst, x)
				e.Set(variadic, lst)
			}
		}

		// Here we go round the evaluation loop again.
		//
		// Which will execute the body of the function this time.
		//
		// TCO.
		exp = proc.Body
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

// macroExpand expands the given macro, recursively.
func (ev *Eval) macroExpand(exp primitive.Primitive, e *env.Environment) primitive.Primitive {

	// is this a macro?
	for ev.isMacro(exp, e) {
		exp = ev.eval(exp, e, false)
	}
	return exp
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

// readExpression uses recursion to read a complete expression from
// our internal array of tokens - as produced by `tokenize`.
func (ev *Eval) readExpression(e *env.Environment) (primitive.Primitive, error) {

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
		quoted, err := ev.readExpression(e)
		if err != nil {
			return nil, err
		}
		return primitive.List{ev.atom("quote"), quoted}, nil

	case "`":
		// `... => (quasiquote ...)
		quoted, err := ev.readExpression(e)
		if err != nil {
			return nil, err
		}
		return primitive.List{ev.atom("quasiquote"), quoted}, nil

	case "~", ",":
		// ~... => (unquote ...)
		quoted, err := ev.readExpression(e)
		if err != nil {
			return nil, err
		}
		return primitive.List{ev.atom("unquote"), quoted}, nil

	case "~@", "`,", ",@":
		// ~@... => (splice-unquote ...)
		quoted, err := ev.readExpression(e)
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
			expr, err := ev.readExpression(e)
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
			key, err := ev.readExpression(e)
			if err != nil {
				return nil, err
			}

			// Check again we've not hit the end of the program
			if ev.offset >= len(ev.toks) {
				return nil, ErrEOF
			}

			// Read the sub-expressions, recursively.
			val, err2 := ev.readExpression(e)
			if err2 != nil {
				return nil, err2
			}

			// Check again we've not hit the end of the program
			if ev.offset >= len(ev.toks) {
				return nil, ErrEOF
			}

			// Ensure the value is evaluated
			v := ev.eval(val, e, true)

			hash.Set(key.ToString(), v)
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

// Does the given list start with a call to the given function?
func (ev *Eval) startsWith(l primitive.List, val string) bool {
	// list must have one entry
	if len(l) < 1 {
		return false
	}

	// entry should match
	return l[0].ToString() == val
}

// tokenize splits the input string into tokens, via a horrific regular
// expression which I don't understand!
func (ev *Eval) tokenize(str string) {

	// Reset our position
	ev.offset = 0
	ev.toks = []string{}

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

// typeCheck is called to type-check arguments.
//
// types contains a ":"-separated list of types that are acceptable, and supplied contains
// the name of the type which was actually supplied.
func (ev *Eval) typeCheck(types string, supplied string) error {

	// types that are allowed
	valid := make(map[string]bool)

	// Is anything possible?
	any := false

	// Record each one
	for _, typ := range strings.Split(types, ":") {

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
		} else {
			valid[typ] = true
		}
	}

	// See if the type matches
	_, ok := valid[supplied]

	if !ok && !any {
		return errors.New("invalid type")
	}

	return nil

}
