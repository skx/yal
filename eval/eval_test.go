package eval

import (
	"testing"

	"github.com/skx/yal/builtins"
	"github.com/skx/yal/env"
	"github.com/skx/yal/stdlib"
)

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

		// lists
		{"'()", "()"},
		{"(car '(1 2 3))", "1"},
		{"(cdr '(1 2 3))", "(2 3)"},

		// numbers
		{"3", "3"},

		// nil
		{"nil", "nil"},

		// mathes
		{"(+ 3 1)", "4"},
		{"(- 3 1)", "2"},
		{"(/ 4 2)", "2"},
		{"(* 4 -1)", "-4"},
		{"(# 3 2)", "9"},

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
