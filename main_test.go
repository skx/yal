package main

import (
	"fmt"
	"testing"

	"github.com/skx/yal/builtins"
	"github.com/skx/yal/env"
	"github.com/skx/yal/eval"
	"github.com/skx/yal/primitive"
	"github.com/skx/yal/stdlib"
)

// The interpreter we're going to execute with.
var interpreter *eval.Eval

// The environment contains the primitives the interpreter uses.
var environment *env.Environment

// Create the interpreter, and parse the source of our benchmark script.
//
// Only do this once, at startup.
func init() {
	// Create a new environment
	environment = env.New()

	// Populate with the default primitives
	builtins.PopulateEnvironment(environment)

	// The script we're going to run
	content := `
(define fact (lambda (n)
  (if (<= n 1)
    1
      (* n (fact (- n 1))))))

(fact 100)
`

	// Read the standard library
	pre := stdlib.Contents()

	// Prepend that to the users' script
	src := string(pre) + "\n" + string(content)

	// Create a new interpreter with that source
	interpreter = eval.New(src)
}

// fact is a benchmark implementation in pure-go for comparison purposes.
func fact(n int64) int64 {
	if n == 0 {
		return 1
	}
	return n * fact(n-1)
}

// BenchmarkGoFactorial allows running the golang benchmark.
func BenchmarkGoFactorial(b *testing.B) {
	for i := 0; i < b.N; i++ {
		fact(100)
	}
}

// BenchmarkYALFactorial allows running the lisp benchmark.
func BenchmarkYALFactorial(b *testing.B) {
	var out primitive.Primitive

	for i := 0; i < b.N; i++ {

		// Run 100!
		out = interpreter.Evaluate(environment)
	}

	// Did we get an error?  Then show it.
	if _, ok := out.(primitive.Error); ok {
		fmt.Printf("Error running: %v\n", out)
	}

}
