package eval

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/skx/yal/builtins"
	"github.com/skx/yal/env"
	"github.com/skx/yal/stdlib"
)

// This tests an infinite loop is handled
func TestTimeout(t *testing.T) {

	// Test code
	tst := `
(define r (lambda () (r)))

(r)
`
	// Load our standard library
	st := stdlib.Contents()
	std := string(st)

	// Timeout after a second
	ctx, cancel := context.WithTimeout(context.Background(), 1000*time.Millisecond)
	defer cancel()

	// Create a new interpreter
	l := New(std + "\n" + tst)

	// Ensure we get a timeout
	l.SetContext(ctx)

	// With a new environment
	env := env.New()

	// Populate the default primitives
	builtins.PopulateEnvironment(env)

	// Run it
	out := l.Evaluate(env)

	if !strings.Contains(out.ToString(), "deadline exceeded") {
		t.Fatalf("Didn't get the expected output:%s", out.ToString())
	}

}

// This function contains a bunch of table-driven tests which are
// designed to be simple.
func TestEvaluate(t *testing.T) {

	type TC struct {
		input  string
		output string
	}

	tests := []TC{
		// bools
		{"#t", "#t"},
		{"true", "#t"},
		{"#f", "#f"},
		{"false", "#f"},

		// if
		{"(if true true false)", "#t"},
		{"(if true true)", "#t"},
		{`(if false "false" "true")`, "true"},
		{"(if false false)", "nil"},

		// lambda
		{`(define sq (lambda (x) (* x x)))
                 ; comment
                 (sq 33)`, "1089"},
		{`(define sqrt (lambda (x) (# x 0.5))) (sqrt 9)`, "3"},
		{`(define sqrt (lambda (x) (# x 0.5))) (sqrt 100)`, "10"},

		// Let
		{"(let ((a 5)) (nil? a))", "#f"},
		{"(let ((a 5)) a)", "5"},
		{"(let ((a 5) (b 6)) a (* a b))", "30"},
		{"(let ((a 5)) (set! a 44) a)", "44"},
		{"(let ((a 5)) c)", "nil"},

		// lists
		{"'()", "()"},
		{"()", "()"},
		{"(car '(1 2 3))", "1"},
		{"(cdr '(1 2 3))", "(2 3)"},
		{"(begin 1 2)", "2"},

		// numbers
		{"3", "3"},

		// nil
		{"nil", "nil"},

		// eval
		{"(eval \"(+ 1 2)\")", "3"},
		{"(let ((a \"(+ 1 23)\")) (eval a))", "24"},
		{"(eval c)", "nil"},

		// cond
		{`(define a 44)
 (cond
  (quote
    (> a 20) "big"
    true     "fallback"))`,
			"big"},

		// maths
		{"(+ 3 1)", "4"},
		{"(- 3 1)", "2"},
		{"(/ 4 2)", "2"},
		{"(* 4 -1)", "-4"},
		{"(# 3 2)", "9"},

		// since we're variadic we start with the first
		// number, and apply the operation to any subsequent
		// ones.
		//
		// so "- 3" is "3 [- nothing]" as there were no more args
		//
		// Explicitly "(- 3)" is NOT "-3".
		//
		{"(+ 3)", "3"},
		{"(- 3)", "3"},
		{"(/ 3)", "3"},
		{"(* 3)", "3"},

		// strings
		{`"steve"`, "steve"},
		{`(split "steve" "")`, `(s t e v e)`},
		{`(join (list "s" "t" "e" "v" "e"))`, `steve`},

		// comparison
		{"(< 1 3)", "#t"},
		{"(< 10 3)", "#f"},
		{"(> 1 3)", "#f"},
		{"(> 10 3)", "#t"},
		{"(= 10 3)", "#f"},
		{"(= 10 10)", "#t"},
		{"(= -1 -1)", "#t"},

		// errors
		{"(invalid)", "ERROR{argument 'invalid' not a function}"},
		{"(eval 'foo 'bar)", "ERROR{Expected only a single argument}"},
		{"(eval 3)", "ERROR{unexpected type for eval %!V(primitive.Number=3).}"},
		{"(let 3)", "ERROR{argument is not a list, got 3}"},
		{"(let ((0 0)) )", "ERROR{binding name is not a symbol, got 0}"},
		{"(let ((0 )) )", "ERROR{arity-error: binding list had missing arguments}"},
		{"(let (3 3) )", "ERROR{binding value is not a list, got 3}"},

		{"(cond (quote 3))", "ERROR{expected pairs of two items}"},
		{"(quote )", "ERROR{arity-error: not enough arguments for (quote}"},
		{"(if )", "ERROR{arity-error: not enough arguments for (if ..)}"},
		{"(define )", "ERROR{arity-error: not enough arguments for (define ..)}"},
		{"(define \"steve\" 3)}", "ERROR{Expected a symbol, got steve}"},
		{"(lambda )}", "ERROR{wrong number of arguments}"},
		{"(lambda 3 4)}", "ERROR{expected a list for arguments, got 3}"},
		{"(define sq (lambda (x) (* x x))) (sq)", "ERROR{arity-error - function 'sq' requires 1 argument(s), 0 provided}"},
		{"(print (/ 3 0))", "ERROR{error expanding argument [/ 3 0]}"},
		{"(lambda (x 3) (nil))}", "ERROR{expected a symbol for an argument, got 3}"},
		{"(set! )", "ERROR{arity-error: not enough arguments for (set! ..)}"},
		{"(let )", "ERROR{arity-error: not enough arguments for (let ..)}"},

		{"))))", "nil"},
		{"'", "nil"},
		{"(3 3 ", "nil"},
		{"(((((", "nil"},
	}

	for _, test := range tests {

		// Load our standard library
		st := stdlib.Contents()
		std := string(st)

		// Create a new interpreter
		l := New(std + "\n" + test.input)

		// With a new environment
		env := env.New()

		// Populate the default primitives
		builtins.PopulateEnvironment(env)

		// Run it
		out := l.Evaluate(env)

		if out.ToString() != test.output {
			t.Fatalf("test '%s' should have produced '%s', but got '%s'", test.input, test.output, out.ToString())
		}

	}
}

// This function tests our standard library.
func TestStandardLibrary(t *testing.T) {

	type TC struct {
		input  string
		output string
	}

	tests := []TC{

		// (boolean?
		{input: "(boolean? true)", output: "#t"},
		{input: "(boolean? #t)", output: "#t"},
		{input: "(boolean? false)", output: "#t"},
		{input: "(boolean? #f)", output: "#t"},
		{input: "(boolean? \"steve\")", output: "#f"},
		{input: "(boolean? (list 1 2 3))", output: "#f"},
		{input: "(boolean? 3)", output: "#f"},

		// first/last
		{input: "(first (list 10 11 12))", output: "10"},
		{input: "(rest  (list 10 11 12))", output: "(11 12)"},

		// inc/dec
		{input: "(inc 10)", output: "11"},
		{input: "(dec 11)", output: "10"},

		// zero?/one?
		{input: "(zero? 0)", output: "#t"},
		{input: "(zero? 10)", output: "#f"},
		// TODO- FIXMIE		{input: "(zero? \"steve\")", output: "#f"},
		{input: "(one? 0)", output: "#f"},
		{input: "(one? 1)", output: "#t"},
		// TODO - FIXME {input: "(one? \"steve\")", output: "#f"},

		// map
		{input: `
(define sq (lambda (x) (* x x)))
(map (list 1 2 3 4 5) (lambda (x) (sq x)))
`,
			output: "(1 4 9 16 25)"},

		// range
		{input: "(range -5 5 1)", output: "(-5 -4 -3 -2 -1 0 1 2 3 4)"},
		{input: "(range 1 11 2)", output: "(1 3 5 7 9)"},

		// seq/nat
		{input: "(seq 10)", output: "(0 1 2 3 4 5 6 7 8 9)"},
		{input: "(nat 10)", output: "(1 2 3 4 5 6 7 8 9)"},

		{input: "(join (reverse (split \"Steve\" \"\")))", output: "evetS"},
	}

	for _, test := range tests {

		// Load our standard library
		st := stdlib.Contents()
		std := string(st)

		// Create a new interpreter
		l := New(std + "\n" + test.input)

		// With a new environment
		env := env.New()

		// Populate the default primitives
		builtins.PopulateEnvironment(env)

		// Run it
		out := l.Evaluate(env)

		if out.ToString() != test.output {
			t.Fatalf("test '%s' should have produced '%s', but got '%s'", test.input, test.output, out.ToString())
		}

	}
}
