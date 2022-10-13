// Package main contains a simple CLI driver for our lisp interpreter.
//
// All the logic is contained within the `main` function, and it merely
// reads the contents of the user-supplied filename, prepends the standard
// library to that content, and executes it.
//
// Notably we don't contain a REPL-mode at the moment.
package main

import (
	"flag"
	"fmt"
	"math/rand"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/skx/yal/builtins"
	"github.com/skx/yal/env"
	"github.com/skx/yal/eval"
	"github.com/skx/yal/primitive"
	"github.com/skx/yal/stdlib"
)

func main() {

	// (gensym) needs a decent random seed
	rand.Seed(time.Now().UnixNano())

	// Look to see if we're gonna execute a statement
	exp := flag.String("e", "", "A string to evaluate.")
	hlp := flag.Bool("h", false, "Should we show help information, and exit?")
	flag.Parse()

	// Ensure we have an argument
	if len(flag.Args()) < 1 && (*exp == "") && !*hlp {
		fmt.Printf("Usage: yal [-e expr] file.lisp\n")
		return
	}

	// Source we'll execute, either from the CLI, or a file
	src := ""

	if *exp != "" {
		src = *exp
	}

	// If we have a file, then read the content
	if len(flag.Args()) > 0 {
		content, err := os.ReadFile(flag.Args()[0])
		if err != nil {
			fmt.Printf("Error reading %s:%s\n", os.Args[1], err)
			return
		}

		src = string(content)
	}

	// Create a new environment
	environment := env.New()

	// If we got any command-line arguments then save them
	if len(flag.Args()) > 0 {

		x := primitive.List{}
		for _, arg := range flag.Args() {
			x = append(x, primitive.String(arg))
		}

		environment.Set("os.args", x)
	} else {
		// Otherwise we'll set an empty list.
		environment.Set("os.args", primitive.List{})
	}

	// Populate the default primitives
	builtins.PopulateEnvironment(environment)

	// Show the help?
	if *hlp {

		// Read the standard library
		pre := stdlib.Contents()

		// Create a new interpreter with that source
		interpreter := eval.New(string(pre))

		// Now evaluate the library, so we get the help for the
		// built-in functions
		interpreter.Evaluate(environment)

		// Build up a list of all things known in the environment
		keys := []string{}

		items := environment.Items()
		for k := range items {
			keys = append(keys, k)
		}

		// sort the items
		sort.Strings(keys)

		for _, key := range keys {

			val, _ := environment.Get(key)

			prc, ok := val.(*primitive.Procedure)
			if ok && len(prc.Help) > 0 {
				txt := prc.Help

				fmt.Printf("%s\n", key)
				fmt.Printf("%s\n", strings.Repeat("=", len(key)))
				fmt.Printf("%s\n\n\n\n", txt)
			}

		}

		return
	}

	// Read the standard library
	pre := stdlib.Contents()

	// Prepend that to the users' script
	src = string(pre) + "\n" + src

	// Create a new interpreter with that source
	interpreter := eval.New(src)

	// Now evaluate the input using the specified environment
	out := interpreter.Evaluate(environment)

	// Did we get an error?  Then show it.
	if _, ok := out.(primitive.Error); ok {
		fmt.Printf("Error running: %v\n", out)
	}
}
