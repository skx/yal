package eval

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/skx/yal/builtins"
	"github.com/skx/yal/env"
	"github.com/skx/yal/primitive"
	"github.com/skx/yal/stdlib"
)

// TestAliased ensures we have some aliases
func TestAliased(t *testing.T) {

	// Load our standard library
	st := stdlib.Contents()
	std := string(st)

	// Create a new interpreter to evaluate 3
	l := New(std + "\n3")

	// With a new environment
	env := env.New()

	// Populate the default primitives
	builtins.PopulateEnvironment(env)

	// Run it
	out := l.Evaluate(env)
	if out.ToString() != "3" {
		t.Fatalf("got wrong result")
	}

	// Get the aliases
	a := l.Aliased()
	size := len(a)

	// Count of functions with a similar name
	//
	// i.e "hms" is like "time:hms"
	//
	similar := 0

	//
	// Look for similarities
	//
	for k, v := range a {

		if strings.Contains(v, k) {
			t.Logf("%s contains %s\n", v, k)
			similar++
		}
	}

	if similar == 0 {
		t.Fatalf("found no similar aliases")
	}

	// Now execute another piece of code, in the same
	// interpreter.
	//
	// We'll add a new alias and see if that gets found
	out = l.Execute(env, "(do (alias testing +) (testing 30 12))")

	if out.ToString() != "42" {
		t.Fatalf("execution had unexpected result")
	}

	if l.aliases["testing"] != "+" {
		t.Fatalf("testing alias wasn't set")
	}

	// Get the updated alias.
	b := l.Aliased()

	// Bigger, by one entry
	if (size + 1) != len(b) {
		t.Fatalf("failed to find newly added alias - old %v new %v", a, b)
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
		// comment
		{"; Foo\n;Bar\n", "nil"},
		{"; Foo\n;Bar\n#f", "#f"},
		{"#!/usr/bin/yal\n;Bar\n", "nil"},
		{"#!/usr/bin/yal\n", "nil"},
		{"#!/usr/bin/yal\n#t", "#t"},

		// integers: in hex
		{"0xff", "255"},
		{"0xFF", "255"},
		{"0xFf", "255"},
		{"0xfF", "255"},

		// integers: in binary
		{"0b00000000", "0"},
		{"0b00000001", "1"},
		{"0b00000010", "2"},
		{"0b00000100", "4"},
		{"0b00001000", "8"},
		{"0b10000000", "128"},

		// bools
		{"#t", "#t"},
		{"true", "#t"},
		{"#f", "#f"},
		{"false", "#f"},

		// character literals
		{"#\\\\n", "\n"},
		{"#\\a", "a"},
		{"#\\AB", "ERROR{invalid character literal: AB}"},

		// literals
		{":foo", ":foo"},

		// hashes
		{"{:age 34}", "{\n\t:age => 34\n}"},
		{"(get {:age 34, :alive true} :alive)", "#t"},
		{"{:age (+ 3 1)}", "{\n\t:age => 4\n}"},

		// if
		{"(if true true false)", "#t"},
		{"(if true true)", "#t"},
		{`(if false "false" "true")`, "true"},
		{"(if false false)", "nil"},

		// symbol
		{`(set! foo (symbol bar)) (symbol? foo)`, `#t`},

		// macroexpand - args are not evaluated
		{`(defmacro! foo (fn* (x) x)) (macroexpand (foo (+ 1 2)))`, "(+ 1 2)"},
		// quote
		{`(define lst (quote (b c)))
                  '(a lst d)`, "(a lst d)"},

		// quasiquote
		{`(define lst (quote (b c)))
                  (quasiquote (a (unquote lst) d))`,
			"(a (b c) d)"},

		// splice-unquote
		{`(define lst (quote (b c)))
                  (quasiquote (a (splice-unquote lst) d))`,
			"(a b c d)"},

		// expand a macro
		{`(defmacro! steve (fn* () "steve"))
                  (macroexpand (steve))`,
			"steve"},

		// lambda
		{`(define sq (lambda (x) (* x x)))
                 ; comment
                 (sq 33)`, "1089"},
		{`(define sqrt (lambda (x) (# x 0.5))) (sqrt 9)`, "3"},
		{`(define sqrt (lambda (x) (# x 0.5))) (sqrt 100)`, "10"},

		// gensym - just test that there's an 11 character return
		{"(length (split (str (gensym)) \"\"))", "11"},

		// Let
		{"(let ((a 5)) (nil? a))", "#f"},
		{"(let ((a 5)) a)", "5"},
		{"(let ((a 5) (b 6)) a (* a b))", "30"},
		{"(let ((a 5)) (set! a 44) a)", "44"},
		{"(let ((a 5)) c)", "nil"},
		{"(let ((a 3) (b (+ a 3))) b))", "6"},

		// let*
		{"(let* (z 9) z)", "9"},
		{"(let* (x 9) x)", "9"},
		{"(let* (z (+ 2 3)) (+ 1 z))", "6"},
		{"(let* (p (+ 2 3) q (+ 2 p)) (+ p q))", "12"},
		{"(def! y (let* (z 7) z)) y", "7"},

		// (set!) inside (let) will only modify the local scope
		{`
(set! a 3)
(let ((b 33))
  (set! a 4321))
a
`,
			"3"},

		// (set! a b TRUE) inside (let) will modify the parent scope
		{`
(define a 3)
(let ((b 33))
  (set! a 4321 true))
a
`,
			"4321"},

		// lists
		{"'()", "()"},
		{"()", "()"},
		{"(car '(1 2 3))", "1"},
		{"(cdr '(1 2 3))", "(2 3)"},
		{"(do 1 2)", "2"},

		// numbers
		{"3", "3"},

		// nil
		{"nil", "nil"},

		// eval
		{"(eval \"(+ 1 2)\")", "3"},
		{"(let ((a \"(+ 1 23)\")) (eval a))", "24"},
		{"(eval c)", "nil"},
		{"(read \"(+ 3 4)\")", "(+ 3 4)"},
		{"(eval (read \"(+ 3 4)\"))", "7"},

		// try caught an error
		{"(try (/ 1 0) (catch e 3))", "3"},
		// try no error to catch
		{"(try (/ 1 1) (catch e 3))", "1"},

		// quoting options
		// quasiquote
		{"`1", "1"},
		{"`(", "nil"},

		// unquote
		{"~`1", "ERROR{argument 'unquote' not a function}"},
		{"~`\"", "nil"},

		// splice-unquote
		{"~@1", "ERROR{argument 'splice-unquote' not a function}"},
		{"~@\"", "nil"},

		// cond
		{`(define a 44)
 (cond
    (> a 20) "big"
    true     "small")`,
			"big"},
		{`(define a 4)
 (cond
    (> a 20) "big"
    true     "small")`,
			"small"},

		{"(cond true 7 true 8)", "7"},
		{"(cond false 7 true 8)", "8"},
		{"(cond false 7 false 8 \"else\" 9)", "9"},
		{"(cond false 7 (= 2 2) 8 \"else\" 9)", "8"},
		{"(cond false 7 false 8 false 9)", "nil"},

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
		{"(* 3)", "3"},

		// strings
		{`"steve"`, "steve"},
		{`(split "steve" "")`, `(s t e v e)`},
		{`(join (list "s" "t" "e" "v" "e"))`, `steve`},

		// alias
		{`(alias explode split) (explode "steve" "")`, `(s t e v e)`},
		{`(alias ** #) (** 3 3)`, `27`},

		// comparison
		{"(< 1 3)", "#t"},
		{"(< 10 3)", "#f"},
		{"(> 1 3)", "#f"},
		{"(> 10 3)", "#t"},
		{"(= 10 3)", "#f"},
		{"(= 10 10)", "#t"},
		{"(= -1 -1)", "#t"},

		// we have a LOT of built ins, but not 200
		{"(> (length (env))  10)", "#t"},
		{"(> (length (env))  50)", "#t"},
		{"(< (length (env)) 200)", "#t"},

		// errors
		{"(symbol)", primitive.ArityError().ToString()},
		{"(symbol 1 2)", primitive.ArityError().ToString()},
		{"(invalid)", "ERROR{argument 'invalid' not a function}"},
		{"(set! 3 4)", "ERROR{tried to set a non-symbol 3}"},
		{"(eval 'foo 'bar)", primitive.ArityError().ToString()},
		{"(eval 3)", "ERROR{unexpected type for eval %!V(primitive.Number=3).}"},
		{"(let 3)", "ERROR{argument is not a list, got 3}"},
		{"(let ((0 0)) )", "ERROR{binding name is not a symbol, got 0}"},
		{"(let ((0 )) )", primitive.ArityError().ToString()},
		{"(let (3 3) )", "ERROR{binding value is not a list, got 3}"},

		{"(let*)", primitive.ArityError().ToString()},
		{"(let* 32)", "ERROR{argument is not a list, got 32}"},
		{"(let* (a 3 b))", "ERROR{list for (len*) must have even length, got [a 3 b]}"},
		{"(let* (a 3 3 b))", "ERROR{binding name is not a symbol, got 3}"},

		{"(struct foo bar)  (type (foo 3))", "struct-foo"},
		{"(struct foo bar)  (foo 3 3)", primitive.ArityError().ToString()},
		{"(struct foo bar)  (foo? nil)", "#f"},
		{"(struct foo bar)  (foo? (foo 3))", "#t"},
		{"(struct a name) (struct b name)  (a? (b 3))", "#f"},
		{"(struct a name) (struct b name)  (b? (b 3))", "#t"},

		{"(struct)", primitive.ArityError().ToString()},
		{"(do (struct foo bar ) (foo?))", primitive.ArityError().ToString()},
		{"(error )", primitive.ArityError().ToString()},
		{"(quote )", primitive.ArityError().ToString()},
		{"(quasiquote )", primitive.ArityError().ToString()},
		{"(macroexpand )", primitive.ArityError().ToString()},
		{"(if )", primitive.ArityError().ToString()},
		{"(if (/ 1 0) #t #f)", "ERROR{attempted division by zero}"},
		{"(define )", primitive.ArityError().ToString()},
		{"(define \"steve\" 3)}", "ERROR{Expected a symbol, got steve}"},
		{"(lambda )}", primitive.ArityError().ToString()},
		{"(lambda 3 4)}", "ERROR{expected a list for arguments, got 3}"},
		{"(define sq (lambda (x) (* x x))) (sq)", primitive.ArityError().ToString()},
		{"(print (/ 3 0))", "ERROR{error expanding argument [/ 3 0] for call to (print ..): ERROR{attempted division by zero}}"},
		{"(lambda (x 3) (nil))}", "ERROR{expected a symbol for an argument, got 3}"},
		{"(set! )", primitive.ArityError().ToString()},
		{"(let )", primitive.ArityError().ToString()},
		{`
(define fizz (lambda (n:number)
  (cond
      (/ n 0)  (print "fizzbuzz")
      #t       (print n))))

(fizz 3)
`, "ERROR{attempted division by zero}"},
		{"(error \"CAKE-FAIL\")", "ERROR{CAKE-FAIL}"},

		{"(defmacro!)", primitive.ArityError().ToString()},
		{"(defmacro! 1 2)", "ERROR{Expected a symbol, got 1}"},
		{"(defmacro! foo 2)", "ERROR{expected a function body for (defmacro..), got 2}"},

		{"(read foo bar)", primitive.ArityError().ToString()},
		{"(read \")\")", "ERROR{failed to read ):unexpected ')'}"},
		{"(read \"}\")", "ERROR{failed to read }:unexpected '}'}"},
		{"'", "nil"},
		{"(3 3 ", "nil"},
		{"(((((", "nil"},
		{"))))", "nil"},
		{"{{{{{{", "nil"},
		{"{ ", "nil"},
		{"{ :name ", "nil"},
		{"{ :name { ", "nil"},
		{"{ :age 333  ", "nil"},
		{"}}}}}}", "nil"},

		{"(alias foo)", primitive.ArityError().ToString()},
		{"(alias foo print)", "nil"},
		{"(alias foo bar print)", "ERROR{(alias ..) must have an even length of arguments, got [foo bar print]}"},

		// try / catch
		{"(try 3)", primitive.ArityError().ToString()},
		{"(try 3 3)", "ERROR{expected a list for argument, got 3}"},
		{"(try (/ 1 0) 3)", "ERROR{expected a list for argument, got 3}"},
		{"(try (/ 1 0) (/ 1 0) (/ 3 9))", "ERROR{catch list should begin with 'catch', got [/ 1 0]}"},
		{"(try (/ 1 0) (catch x))", "ERROR{list should have three elements, got [catch x]}"},

		// type failures
		{input: "(define blah (lambda (a:list) (print a))) (blah 3)", output: "ERROR{type-validation failed: argument a to blah was supposed to be list, got number}"},
		{input: "(define blah (lambda (a:string) (print a))) (blah 3)", output: "ERROR{type-validation failed: argument a to blah was supposed to be string, got number}"},
		{input: "(define blah (lambda (a:number) (print a))) (blah '(3))", output: "ERROR{type-validation failed: argument a to blah was supposed to be number, got list}"},
		{input: "(define blah (lambda (a:function) (print a))) (blah '(3))", output: "ERROR{type-validation failed: argument a to blah was supposed to be function, got list}"},
		{input: "(define blah (lambda (a:any) (print a))) (blah '(3))", output: "(3)"},
	}

	for _, test := range tests {

		t.Run(test.input, func(t *testing.T) {

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
		})
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
		{input: "(zero? \"steve\")", output: "ERROR{argument was not a number}"},
		{input: "(one? 0)", output: "#f"},
		{input: "(one? 1)", output: "#t"},
		{input: "(one? \"steve\")", output: "ERROR{argument was not a number}"},

		// map
		{input: `
(define sq (lambda (x) (* x x)))
(map (list 1 2 3 4 5) (lambda (x) (sq x)))
`,
			output: "(1 4 9 16 25)"},

		// range
		{input: "(range -5 5 1)", output: "(-5 -4 -3 -2 -1 0 1 2 3 4 5)"},
		{input: "(range 1 11 2)", output: "(1 3 5 7 9 11)"},

		// seq/nat
		{input: "(seq 10)", output: "(0 1 2 3 4 5 6 7 8 9 10)"},
		{input: "(nat 10)", output: "(1 2 3 4 5 6 7 8 9 10)"},

		{input: "(join (reverse (split \"Steve\" \"\")))", output: "evetS"},
	}

	for _, test := range tests {

		t.Run(test.input, func(t *testing.T) {

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
		})

	}
}

func TestStartsWith(t *testing.T) {
	l := primitive.List{}

	e := New("")

	if e.startsWith(l, "steve") {
		t.Fatalf("unexpected match")
	}
}

// TestStdlibHelp is designed to ensure that our standard library
// functions have help-documentation.
func TestStdlibHelp(t *testing.T) {

	// Load our standard library
	st := stdlib.Contents()
	std := string(st)

	// Create a new interpreter
	l := New(std)

	// With a new environment
	env := env.New()

	// Populate the default primitives
	builtins.PopulateEnvironment(env)

	// Run it
	_ = l.Evaluate(env)

	// Now we should have an environment which is
	// populated with functions
	for name, val := range env.Items() {

		proc, ok := val.(*primitive.Procedure)

		if !ok {
			t.Skip("ignoring non-procedure entry in environment " + name)
		}

		t.Run(name, func(t *testing.T) {

			if len(proc.Help) == 0 {
				t.Fatalf("empty help for %s", name)
			}
		})
	}

}

// This tests an infinite loop is handled
func TestTimeout(t *testing.T) {

	// Test code - run an infinite loop, incrementing a variable.
	tst := `
(set! a 1)
(while true
 (do
   (sleep)
   (set! a (+ a 1) true)))
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
	ev := env.New()

	// Add a new function
	ev.Set("sleep",
		&primitive.Procedure{
			Help: "sleep delays for two seconds",
			F: func(e *env.Environment, args []primitive.Primitive) primitive.Primitive {
				fmt.Printf("Sleeping for two seconds")
				time.Sleep(2 * time.Second)
				return primitive.Nil{}
			}})

	// Populate the default primitives
	builtins.PopulateEnvironment(ev)

	// Run it
	out := l.Evaluate(ev)

	// Test for both possible errors here.
	//
	// We should get the context error, but sometimes we don't
	// the important thing is we DON'T hang forever
	if !strings.Contains(out.ToString(), "deadline exceeded") {
		t.Fatalf("Didn't get the expected output.  Got: %s", out.ToString())
	}
}
