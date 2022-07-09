package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/skx/yal/env"
	"github.com/skx/yal/eval"
	"github.com/skx/yal/primitive"
	"github.com/skx/yal/stdlib"
)

func main() {

	if len(os.Args) < 1 {
		fmt.Printf("Usage: ml file.lisp\n")
		return
	}

	content, err := ioutil.ReadFile(os.Args[1])
	if err != nil {
		fmt.Printf("error reading %s:%s\n", os.Args[1], err)
		return
	}

	// Create a new environment
	environment := env.New()

	// Populate the default primitives
	primitive.PopulateEnvironment(environment)

	// Read the standard library
	pre := stdlib.Contents()

	// Prepend that to the users' script
	src := string(pre) + "\n" + string(content)

	// Create a new interpreter with that source
	interpreter := eval.New(src)

	// Now evaluate the input using the specified environment
	out := interpreter.Evaluate(environment)

	// Show the result
	fmt.Printf("%v\n", out)
}
