//go:build go1.18
// +build go1.18

package main

import (
	"context"
	"os"
	"path"
	"strings"
	"testing"
	"time"

	"github.com/skx/yal/builtins"
	"github.com/skx/yal/env"
	"github.com/skx/yal/eval"
	"github.com/skx/yal/primitive"
	"github.com/skx/yal/stdlib"
)

func FuzzYAL(f *testing.F) {

	// We're running fuzzing, and that means we need
	// to disable "shell".  That is done via the use
	// of an environmental variable
	os.Setenv("FUZZ", "FUZZ")

	// empty string
	f.Add([]byte(""))

	// simple entries
	f.Add([]byte("(/ 1 30)"))
	f.Add([]byte("(print (+ 3 2))"))
	f.Add([]byte("()"))
	f.Add([]byte("; This is a comment"))
	f.Add([]byte("(list 3 4 5)"))

	// bigger entries
	f.Add([]byte(`
(print "Our mathematical functions allow 2+ arguments, e.g: %s = %s"
  (quote (+ 1 2 3 4 5 6)) (+ 1 2 3 4 5 6))
`))
	f.Add([]byte(`
;; Define a function, 'fact', to calculate factorials.
(define fact (lambda (n)
  (if (<= n 1)
    1
      (* n (fact (- n 1))))))

;; Invoke the factorial function, using apply
(apply (list 1 2 3 4 5 6 7 8 9 10)
  (lambda (x)
    (print "%s! => %s" x (fact x))))
`))

	f.Add([]byte(`
; Split a string into a list, reverse it, and join it
(let* (input "Steve Kemp")
  (do
   (print "Starting string: %s" input)
   (print "Reversed string: %s" (join (reverse (split "Steve Kemp" ""))))))
`))

	f.Add([]byte(`
;; Now create a utility function to square a number
(define sq (lambda (x) (* x x)))

;; For each item in the range 1-10, print it, and the associated square.
;; Awesome!  Much Wow!
(apply (nat 11)
      (lambda (x)
        (print "%s\tsquared is %s" x (sq x))))
`))

	f.Add([]byte(`
;;
;; Setup a list of integers, and do a few things with it.
;;
(let* (vals '(32 92 109 903 31 3 -93 -31 -17 -3))
     (print "Working with the list: %s " vals)
     (print "\tBiggest item is %s"       (max vals))
     (print "\tSmallest item is %s"      (min vals))
     (print "\tReversed list is %s "     (reverse vals))
     (print "\tSorted list is %s "       (sort vals))
     (print "\tFirst item is %s "        (first vals))
     (print "\tRemaining items %s "      (rest vals))
   )
`))

	f.Add([]byte(`
;; We have a built-in eval function, which operates upon symbols, or strings.
(define e "(+ 3 4)")
(print "Eval of '%s' resulted in %s" e (eval e))
`))

	// Recurse forever
	f.Add([]byte(`
(define r (lambda () (r))) (r)`))

	// Macros
	f.Add([]byte(`
(defmacro! unless (fn* (pred a &b) ` + "`" + `(if (! ~pred) ~a ~b)))
(unless false (print "OK")
`))

	// Type checking
	f.Add([]byte(`define blah (lambda (a:list) (print "I received the list %s" a)))`))
	f.Add([]byte(`define blah (lambda (a:string) (print "I received the string %s" a)))`))
	f.Add([]byte(`define blah (lambda (a:number) (print "I received the number %s" a)))`))
	f.Add([]byte(`define blah (lambda (a:any) (print "I received the arg %s" a)))`))
	f.Add([]byte(`define blah (lambda (a) (print "I received the arg %s" a)))`))

	// Find each of our examples, as these are valid code samples
	files, err := os.ReadDir("examples")
	if err != nil {
		f.Fatalf("failed to read examples/ directory %s", err)
	}

	// Load each example as a fuzz-source
	for _, file := range files {
		path := path.Join("examples", file.Name())

		data, err := os.ReadFile(path)
		if err != nil {
			f.Fatalf("Failed to load %s %s", path, err)
		}
		f.Add(data)
	}

	// Known errors are listed here.
	//
	// The purpose of fuzzing is to find panics, or unexpected errors.
	// Some programs are obviously invalid though, and we don't want to
	// report those known-bad things.
	//
	known := []string{
		"arityerror",
		"catch list should begin with 'catch'", // try/catch
		"deadline exceeded",                    // context timeout
		"division by zero",
		"error expanding argument",
		"expected a function body",
		"expected a list",
		"expected a symbol",
		"failed to compile regexp",
		"failed to open", // file:lines
		"invalid character literal",
		"is not a symbol",
		"list should have three elements", // try
		"must be greater than zero",       // random
		"must have even length",
		"not a character",
		"not a function",
		"not a hash",
		"not a list",
		"not a number",
		"not a procedure",
		"not a string",
		"out of bounds", // nth
		"recursion limit",
		"syntax error in pattern", // glob
		"tried to set a non-symbol",
		"typeerror - ",
		"unexpected type",
	}

	// Read the standard library only once.
	std := string(stdlib.Contents()) + "\n"

	f.Fuzz(func(t *testing.T, input []byte) {

		// Timeout after a second
		ctx, cancel := context.WithTimeout(context.Background(), 1000*time.Millisecond)
		defer cancel()

		// Create a new environment
		environment := env.New()

		// Populate the default primitives
		builtins.PopulateEnvironment(environment)

		// Prepend the standard-library to the users' script
		src := std + string(input)

		// Create a new interpreter with the combined source
		interpreter := eval.New(src)

		// Ensure we timeout after 1 second
		interpreter.SetContext(ctx)

		// Now evaluate the input using the specified environment
		out := interpreter.Evaluate(environment)

		switch out.(type) {
		case *primitive.Error, primitive.Error:
			str := strings.ToLower(out.ToString())

			// does it look familiar?
			for _, v := range known {
				if strings.Contains(str, v) {
					return
				}
			}
			t.Fatalf("error processing input %s:%v", input, out)
		}
	})
}
