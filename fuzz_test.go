//go:build go1.18
// +build go1.18

package main

import (
	"strings"
	"testing"

	"github.com/skx/yal/builtins"
	"github.com/skx/yal/env"
	"github.com/skx/yal/eval"
	"github.com/skx/yal/primitive"
	"github.com/skx/yal/stdlib"
)

func FuzzYAL(f *testing.F) {

	// empty string
	f.Add([]byte(""))

	// simple entries
	f.Add([]byte("(/ 1 30)"))
	f.Add([]byte("(print (+ 3 2))"))
	f.Add([]byte("()"))
	f.Add([]byte("(list 3 4 5)"))

	// Known errors are listed here.
	//
	// The purpose of fuzzing is to find panics, or unexpected errors.
	//
	// Some programs are obviously invalid though, so we don't want to
	// report those known-bad things.
	known := []string{
		"not a function",
	}

	f.Fuzz(func(t *testing.T, input []byte) {

		// Create a new environment
		environment := env.New()

		// Populate the default primitives
		builtins.PopulateEnvironment(environment)

		// Read the standard library
		pre := stdlib.Contents()

		// Prepend that to the users' script
		src := string(pre) + "\n" + string(input)

		// Create a new interpreter with that source
		interpreter := eval.New(src)

		// Now evaluate the input using the specified environment
		out := interpreter.Evaluate(environment)

		found := false

		switch out.(type) {
		case *primitive.Error, primitive.Error:
			// does it look familiar?
			for _, v := range known {
				if strings.Contains(out.ToString(), v) {
					found = true
				}
			}

			// raise an error
			if !found {
				t.Fatalf("error parsing %s:%v", input, out)
			}
		}
	})
}
