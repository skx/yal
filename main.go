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
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/skx/yal/builtins"
	"github.com/skx/yal/env"
	"github.com/skx/yal/eval"
	"github.com/skx/yal/primitive"
	"github.com/skx/yal/stdlib"
)

var (
	version = "unreleased"
	sha1sum = "unknown"
)

func main() {

	// (gensym) needs a decent random seed, as does (random).
	rand.Seed(time.Now().UnixNano())

	// Look to see if we're gonna execute a statement, or dump our help.
	exp := flag.String("e", "", "A string to evaluate.")
	hlp := flag.Bool("h", false, "Should we show help information, and exit?")
	ver := flag.Bool("v", false, "Should we show our version, and exit?")
	flag.Parse()

	// Showing the version?
	if *ver {
		fmt.Printf("%s [%s]\n", version, sha1sum)
		return
	}

	// Ensure we have an argument, if we don't have flags.
	if len(flag.Args()) < 1 && (*exp == "") && !*hlp {
		fmt.Printf("Usage: yal [-e expr] [-h] [-v] file.lisp\n")
		return
	}

	// Source we'll execute, either from the CLI, or a file
	src := ""

	if *exp != "" {
		src = *exp
	}

	// If we have a file, then read the content
	if len(flag.Args()) > 0 && !*hlp {
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

		// Are we just showing a specific thing?
		show := flag.Args()

		// Read the standard library
		pre := stdlib.Contents()

		// Create a new interpreter with that source
		interpreter := eval.New(string(pre))

		// Now evaluate the library, so that all our core
		// functions are defined.
		//
		// We can only get help for functions which are present.
		interpreter.Evaluate(environment)

		// Show aliased functions, separately
		aliased := interpreter.Aliased()

		// Build up a list of all things known in the environment
		keys := []string{}

		// Save the known "things", because we want show them
		// in sorted-order.
		items := environment.Items()
		for k := range items {
			keys = append(keys, k)
		}

		// sort the items
		sort.Strings(keys)

		// Now we have a list of sorted things.
		for _, key := range keys {

			// get the item.
			val, _ := environment.Get(key)

			// Is it a procedure?
			prc, ok := val.(*primitive.Procedure)

			// If it isn't a procedure skip it
			if !ok {
				continue
			}

			// No help? skip it
			if len(prc.Help) == 0 {
				continue
			}

			// Get the text
			txt := prc.Help

			// Is this an aliased function?
			target, ok := aliased[key]
			if ok {

				txt = fmt.Sprintf("%s is an alias for %s.", key, target)
			}

			// Build up the arguments
			args := ""

			if len(prc.Args) > 0 {

				for _, arg := range prc.Args {
					args += " " + arg.ToString()
				}
				args = strings.TrimSpace(args)
				args = " (" + args + ")"
			}

			entry := key + args + "\n"
			entry += strings.Repeat("=", len(key+args)) + "\n"
			entry += txt + "\n\n\n"

			// Are we going to show this?
			match := false
			for _, x := range show {

				r, er := regexp.Compile(x)
				if er != nil {
					fmt.Printf("Error compiling regexp %s:%s", show, er)
					return
				}

				res := r.FindStringSubmatch(entry)
				if len(res) > 0 {
					match = true
				}
			}

			if (len(show) == 0) || match {
				fmt.Printf("%s", entry)
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
