// Package main contains a simple CLI driver for our lisp interpreter.
//
// All the logic is contained within the `main` function, and it merely
// reads the contents of the user-supplied filename, prepends the standard
// library to that content, and executes it.
package main

import (
	"flag"
	"fmt"
	"os"
	"path"
	"regexp"
	"sort"
	"strings"

	"github.com/chzyer/readline"

	"github.com/skx/yal/builtins"
	"github.com/skx/yal/config"
	"github.com/skx/yal/env"
	"github.com/skx/yal/eval"
	"github.com/skx/yal/primitive"
	"github.com/skx/yal/stdlib"
)

var (
	version = "unreleased"
	sha1sum = "unknown"

	// ENV is the environment the interpreter uses.
	ENV *env.Environment

	// LISP is the actual interpreter.
	LISP *eval.Eval
)

// versionFn is the implementation of the (version) primitive.
func versionFn(env *env.Environment, args []primitive.Primitive) primitive.Primitive {
	return primitive.String(version)
}

// create handles the setup of our global interpreter and environment.
//
// The standard-library will be loaded, and os.args will be populated.
func create() {

	// Create a new environment
	ENV = env.New()

	// Setup the I/O
	ENV.SetIOConfig(config.DefaultIO())

	// Populate the default primitives
	builtins.PopulateEnvironment(ENV)

	// Add the (version) function
	ENV.Set("version",
		&primitive.Procedure{
			F:    versionFn,
			Help: "Return the version of the interpreter.\n\nSee-also: arch, os",
			Args: []primitive.Symbol{}})

	// Build up a list of the command-line arguments
	args := primitive.List{}

	// Adding them to the list
	for _, arg := range flag.Args() {
		args = append(args, primitive.String(arg))
	}

	// Before setting them in the environment
	ENV.Set("os.args", args)

	// Read the standard library
	txt := stdlib.Contents()

	// Create a new interpreter with that source
	LISP = eval.New(string(txt))

	// Now evaluate the input using the specified environment
	out := LISP.Evaluate(ENV)

	// Did we get an error?  Then show it.
	if _, ok := out.(primitive.Error); ok {
		fmt.Printf("Error executing standard-library: %v\n", out)
		os.Exit(1)
	}
}

// help - show help information.
//
// Either all functions, or only those that match the regular expressions
// supplied.
func help(show []string) {

	// Patterns is a cache of regexps, to ensure we only compile
	// them once.
	var patterns []*regexp.Regexp

	// Compile each supplied pattern, and save it away.
	for _, pat := range show {

		r, er := regexp.Compile(pat)
		if er != nil {
			fmt.Printf("Error compiling regexp %s:%s", show, er)
			return
		}

		patterns = append(patterns, r)
	}

	// We want to show aliased functions separately, so we have to
	// find them - via the interpreter which executed the stdlib
	// at create() time.
	aliased := LISP.Aliased()

	// Build up a list of all things known in the environment
	keys := []string{}

	// Save the known "things", because we want show them in sorted-order.
	items := ENV.Items()
	for k := range items {
		keys = append(keys, k)
	}

	// sort the known-things (i.e. environment keys)
	sort.Strings(keys)

	// Now we have a list of sorted things.
	for _, key := range keys {

		// get the item from the environment.
		val, _ := ENV.Get(key)

		// Is it a procedure?
		prc, ok := val.(*primitive.Procedure)

		// If it isn't a procedure skip it.
		if !ok {
			continue
		}

		// If there is no help then skip it.
		if len(prc.Help) == 0 {
			continue
		}

		// Get the text
		txt := prc.Help

		// Is this an aliased function?
		target, ok := aliased[key]
		if ok {
			// If so change the text.
			txt = fmt.Sprintf("%s is an alias for %s.", key, target)
		}

		// Build up the arguments to the procedure.
		args := ""

		if len(prc.Args) > 0 {

			for _, arg := range prc.Args {
				args += " " + arg.ToString()

				// Default value for this argument?
				def, ok := prc.Defaults[arg]
				if ok {
					args += "[default:"
					args += def.ToString()
					args += "]"
				}
			}
			args = strings.TrimSpace(args)
			args = " (" + args + ")"
		}

		// Build up a complete list of the entry we'll output.
		entry := key + args + "\n"
		entry += strings.Repeat("=", len(key+args)) + "\n"
		entry += txt + "\n\n\n"

		// Are we going to show this?
		//
		// No filtering?  Then yes
		if len(show) == 0 {
			fmt.Printf("%s", entry)
			continue
		}

		// Otherwise test each supplied pattern against the text,
		// and if one matches show it and continue.
		for _, test := range patterns {

			res := test.FindStringSubmatch(entry)
			if len(res) > 0 {
				fmt.Printf("%s", entry)
				continue
			}
		}
	}
}

func main() {

	// define our command-line flags
	exp := flag.String("e", "", "A string to evaluate.")
	hlp := flag.Bool("h", false, "Show help information and exit.")
	lsp := flag.Bool("lsp", false, "Launch the LSP mode")
	ver := flag.Bool("v", false, "Show our version and exit.")
	deb := flag.Bool("debug", false, "Show debug output during execution (to STDERR).")

	// Parse our command-line flags
	flag.Parse()

	// Showing the version?
	if *ver {
		fmt.Printf("%s [%s]\n", version, sha1sum)
		return
	}

	// create the interpreter.
	//
	// This populates the environment, by executing the standard-library.
	//
	// This saves time because:
	//
	//     -h will require the stdlib to be loaded, to dump help info.
	//
	// OR
	//
	//     executing the users' code, via "-e" or a file, will need
	//     that present too.
	//
	create()

	//
	// By default we have no STDERR handler wired up, but if we set the
	// debug flag we'll send that to the actual console's STDERR stream
	if *deb {

		// Get config
		iohelper := ENV.GetIOConfig()

		// Setup a destination for STDERR
		iohelper.STDERR = os.Stderr

		// Update
		ENV.SetIOConfig(iohelper)
	}

	// LSP?
	if *lsp {
		lspStart()
		return
	}

	// showing the help?
	if *hlp {
		help(flag.Args())
		return
	}

	// Executing an expression?
	if *exp != "" {

		// Now evaluate the input using the specified environment
		out := LISP.Execute(ENV, string(*exp))

		// Did we get an error?  Then show it.
		if _, ok := out.(primitive.Error); ok {
			fmt.Printf("Error executing the supplied expression: %v\n", out)
			os.Exit(1)
		}
		os.Exit(0)
	}

	// If we have a file, then read the content.
	if len(flag.Args()) > 0 {
		content, err := os.ReadFile(flag.Args()[0])
		if err != nil {
			fmt.Printf("Error reading %s:%s\n", os.Args[1], err)
			return
		}

		// Now evaluate the input using the specified environment
		out := LISP.Execute(ENV, string(content))

		// Did we get an error?  Then show it.
		if _, ok := out.(primitive.Error); ok {
			fmt.Printf("Error executing %s: %v\n", os.Args[1], out)
			os.Exit(1)
		}
		os.Exit(0)
	}

	//
	// We'll read a config file from HOME, if it exists, and also persist history there.
	//
	home := os.Getenv("HOME")

	//
	// History file will be set if HOME wasn't empty.
	//
	hist := ""
	if home != "" {
		hist = path.Join(home, ".yal.history")
	}

	//
	// No arguments mean this is our REPL
	//
	// Create a readline-helper for reading the input from the user.
	//
	rl, err := readline.NewEx(&readline.Config{
		Prompt:                 "> ",
		HistoryFile:            hist,
		HistorySearchFold:      true,
		DisableAutoSaveHistory: false,
	})

	//
	// There should be no error creating our readline-helper,
	// but if there is then it is fatal.
	//
	if err != nil {
		fmt.Printf("Failed to initialize readlin: %s\n", err)
		os.Exit(1)
	}
	defer rl.Close()

	//
	// Get the home directory, and load ~/.yalrc if present
	//
	if home != "" {

		// Build the path
		file := path.Join(home, ".yalrc")

		// Read the content
		content, err := os.ReadFile(file)
		if err == nil {

			// Execute the contents
			out := LISP.Execute(ENV, string(content))
			if _, ok := out.(primitive.Error); ok {
				fmt.Printf("Error executing ~/.yalrc %v\n", out)
			}
		}
	}

	//
	// We allow multi-line input, and build up the thing to execute
	// into this temporary string if that is the case.
	//
	src := ""

	for {

		//
		// Different prompt for first-line and additional lines
		//
		if src == "" {
			rl.SetPrompt("> ")
		} else {
			rl.SetPrompt(">>> ")
		}

		//
		// Read input
		//
		line, err := rl.Readline()
		if err != nil {
			break
		}

		//
		// Save it to anything we might have previously read
		//
		src += line
		src = strings.TrimSpace(src)

		// Allow the user to exit
		if src == "exit" || src == "quit" {
			break
		}

		//
		// Attempt to guess if the line is complete.
		//
		// A line is complete if there are matching numbers
		// of opening and closing brackets.
		///
		open := strings.Count(src, "(")
		close := strings.Count(src, ")")

		if open < close {
			fmt.Printf("Malformed expression: %v", src)
			src = ""
			continue
		}
		if open == close {

			out := LISP.Execute(ENV, src)

			// If the result wasn't nil then show it
			if _, ok := out.(primitive.Nil); !ok {
				fmt.Printf("%v\n", out.ToString())
			}

			src = ""
		}
	}

	// All done.
}
