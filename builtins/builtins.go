// Package builtins contains the implementations of the lisp-callable
// functions which are implemented in golang.
package builtins

import (
	"fmt"
	"math"
	"sort"
	"strconv"
	"strings"

	"github.com/skx/yal/env"
	"github.com/skx/yal/primitive"
)

// PrimitiveFn is the type which represents a function signature for
// a lisp-usable function implemented in golang
type PrimitiveFn func(args []primitive.Primitive) primitive.Primitive

// PopulateEnvironment registers our default primitives
func PopulateEnvironment(env *env.Environment) {

	// maths
	env.Set("+", &primitive.Procedure{F: plusFn})
	env.Set("-", &primitive.Procedure{F: minusFn})
	env.Set("*", &primitive.Procedure{F: multiplyFn})
	env.Set("/", &primitive.Procedure{F: divideFn})
	env.Set("%", &primitive.Procedure{F: modFn})
	env.Set("#", &primitive.Procedure{F: expnFn})

	// comparisions
	env.Set("<", &primitive.Procedure{F: ltFn})
	env.Set("=", &primitive.Procedure{F: eqFn})

	// When it comes to comparisons there are several we could
	// use:
	//
	//  <
	//  <=
	//  >
	//  >=
	//  =
	//
	// We only actually need to implement "<" in Golang, the rest
	// can be added in lisp.

	// List
	env.Set("list", &primitive.Procedure{F: func(args []primitive.Primitive) primitive.Primitive {
		return primitive.List(args)
	}})
	env.Set("car", &primitive.Procedure{F: func(args []primitive.Primitive) primitive.Primitive {
		if len(args) != 1 {
			return primitive.Error("wrong number of arguments")
		}

		// ensure we received a list
		if _, ok := args[0].(primitive.List); !ok {
			return primitive.Error("argument not a list")
		}

		return args[0].(primitive.List)[0]
	}})
	env.Set("cdr", &primitive.Procedure{F: func(args []primitive.Primitive) primitive.Primitive {
		if len(args) != 1 {
			return primitive.Error("wrong number of arguments")
		}

		// ensure we received a list
		if _, ok := args[0].(primitive.List); !ok {
			return primitive.Error("argument not a list")
		}

		return args[0].(primitive.List)[1:]
	}})
	env.Set("sort", &primitive.Procedure{F: func(args []primitive.Primitive) primitive.Primitive {

		// If we have only a single argument
		if len(args) != 1 {
			return primitive.Error("invalid argument count")
		}

		// Which is a list
		if _, ok := args[0].(primitive.List); !ok {
			return primitive.Error("argument not a list")
		}

		// Cast
		l := args[0].(primitive.List)

		// Copy
		var c primitive.List
		c = append(c, l...)

		// Sort the copy of the list
		sort.Slice(c, func(i, j int) bool {

			// If we have numbers we can sort
			if _, ok := c[i].(primitive.Number); ok {
				if _, ok := c[j].(primitive.Number); ok {

					a, _ := strconv.ParseFloat(c[i].ToString(), 64)
					b, _ := strconv.ParseFloat(c[j].ToString(), 64)

					return a < b
				}
			}

			// Otherwise we sort as strings
			a := c[i].ToString()
			b := c[j].ToString()
			return a < b
		})

		return c
	}})

	// nil
	env.Set("nil?", &primitive.Procedure{F: func(args []primitive.Primitive) primitive.Primitive {
		// nil is nil (yeah, really)
		if primitive.IsNil(args[0]) {
			return primitive.Bool(true)
		}

		// an empty list is nil.
		if list, ok := args[0].(primitive.List); ok {
			return primitive.Bool(len(list) == 0)
		}
		return primitive.Bool(false)
	}})

	env.Set("cons", &primitive.Procedure{F: func(args []primitive.Primitive) primitive.Primitive {
		if len(args) == 1 {
			return primitive.List{args[0]}
		}
		if args[1] == nil || primitive.IsNil(args[1]) {
			return primitive.List{args[0]}
		}
		if _, ok := args[1].(primitive.List); ok {
			return append(primitive.List{args[0]}, args[1].(primitive.List)...)
		}
		return primitive.List{args[0], args[1]}
	}})

	// string functions
	env.Set("split", &primitive.Procedure{F: func(args []primitive.Primitive) primitive.Primitive {

		// We require two arguments
		if len(args) != 2 {
			return primitive.Error("invalid argument count")
		}

		// Both arguments must be strings
		if _, ok := args[0].(primitive.String); !ok {
			return primitive.Error("argument not a string")
		}
		if _, ok := args[1].(primitive.String); !ok {
			return primitive.Error("argument not a string")
		}

		// split
		out := strings.Split(args[0].ToString(), args[1].ToString())

		var c primitive.List

		for _, x := range out {
			c = append(c, primitive.String(x))
		}

		return c
	}})

	// string functions
	env.Set("join", &primitive.Procedure{F: func(args []primitive.Primitive) primitive.Primitive {

		// We require one argument
		if len(args) != 1 {
			return primitive.Error("invalid argument count")
		}

		// The argument must be a list
		if _, ok := args[0].(primitive.List); !ok {
			return primitive.Error("argument not a list")
		}

		tmp := ""

		for _, t := range args[0].(primitive.List) {
			tmp += t.ToString()
		}

		return primitive.String(tmp)
	}})

	// type
	env.Set("type", &primitive.Procedure{F: func(args []primitive.Primitive) primitive.Primitive {
		if len(args) != 1 {
			return primitive.Error("wrong number of arguments")
		}
		return primitive.String(args[0].Type())
	}})

	// logic
	env.Set("or", &primitive.Procedure{F: func(args []primitive.Primitive) primitive.Primitive {

		// For each argument
		for _, arg := range args {

			switch v := arg.(type) {

			// Bool?
			case primitive.Bool:
				if v {
					return primitive.Bool(true)
				}

			// list
			case primitive.List:

				for _, a := range v {

					// See if we can cast to a bool
					b, ok := a.(primitive.Bool)
					if ok {
						// OK it was - is it true?
						if b {
							return primitive.Bool(true)
						}
					} else {
						if !primitive.IsNil(arg) {
							return primitive.Bool(true)
						}
					}
				}
			}
		}
		return primitive.Bool(false)
	}})

	env.Set("and", &primitive.Procedure{F: func(args []primitive.Primitive) primitive.Primitive {

		// For each argument
		for _, arg := range args {

			switch v := arg.(type) {

			// Bool?
			case primitive.Bool:
				if !v {
					return primitive.Bool(false)
				}

			// Nil
			case primitive.Nil:
				return primitive.Bool(false)

			// list
			case primitive.List:

				for _, a := range v {

					// See if we can cast to a bool
					b, ok := a.(primitive.Bool)
					if ok {
						// OK it was - is it true?
						if !b {
							return primitive.Bool(false)
						}
					} else {
						if primitive.IsNil(arg) {
							return primitive.Bool(false)
						}
					}
				}
			}
		}
		return primitive.Bool(true)
	}})

	// Print
	env.Set("sprintf", &primitive.Procedure{F: func(args []primitive.Primitive) primitive.Primitive {
		// no args
		if len(args) < 1 {
			return primitive.Nil{}
		}

		// one arg
		if len(args) == 1 {
			// expand
			str := expandStr(args[0].ToString())

			// show & return
			fmt.Println(str)
			return primitive.String(str)
		}

		// OK format-string
		frmt := expandStr(args[0].ToString())
		parm := []any{}

		for i, a := range args {
			if i == 0 {
				continue
			}
			parm = append(parm, a.ToString())
		}

		out := fmt.Sprintf(frmt, parm...)
		return primitive.String(out)
	}})

	env.Set("print", &primitive.Procedure{F: func(args []primitive.Primitive) primitive.Primitive {
		// no args
		if len(args) < 1 {
			return primitive.Nil{}
		}

		// one arg
		if len(args) == 1 {
			// expand
			str := expandStr(args[0].ToString())

			// show & return
			fmt.Println(str)
			return primitive.String(str)
		}

		// OK format-string
		frmt := expandStr(args[0].ToString())
		parm := []any{}

		for i, a := range args {
			if i == 0 {
				continue
			}
			parm = append(parm, a.ToString())
		}

		out := fmt.Sprintf(frmt, parm...)
		fmt.Println(out)
		return primitive.String(out)
	}})
	// Convert an object to a string
	env.Set("str", &primitive.Procedure{F: func(args []primitive.Primitive) primitive.Primitive {
		return primitive.String(args[0].ToString())
	}})

}

// Convert a string such as "steve\tkemp" into "steve<TAB>kemp"
func expandStr(input string) string {
	out := ""

	// Walk the string character by character
	i := 0
	l := len(input)

	for i < l {

		// current character
		c := input[i]

		// look for "\n", "\t", etc.
		if c == '\\' && (i+1) < l {

			next := input[i+1]
			switch next {
			case 't':
				out += "\t"
			case 'n':
				out += "\n"
			case 'r':
				out += "\r"
			case '\\':
				out += "\\"
			default:
				// unknown escapes will be left alone
				out += "\\" + string(next)
			}

			// Bump the count once, to skip the "\".
			//
			// At the end of the loop we bump again, which will
			// skip the character after that
			i++
		} else {
			out += string(c)
		}
		i++
	}

	return out
}

// plusFn implements "+"
func plusFn(args []primitive.Primitive) primitive.Primitive {

	// ensure we have at least one argument
	if len(args) < 1 {
		return primitive.Error("invalid argument count")
	}

	// the first argument must be a number.
	v, ok := args[0].(primitive.Number)
	if !ok {
		return primitive.Error(fmt.Sprintf("argument '%s' was not a number", args[0].ToString()))
	}

	// now process all the rest of the arguments
	for _, i := range args[1:] {

		// check we have a number
		n, ok := i.(primitive.Number)
		if ok {
			v += n
		} else {
			return primitive.Error(fmt.Sprintf("argument %s was not a number", i.ToString()))
		}
	}
	return primitive.Number(v)
}

// minusFn implements "+"
func minusFn(args []primitive.Primitive) primitive.Primitive {

	// ensure we have at least one argument
	if len(args) < 1 {
		return primitive.Error("invalid argument count")
	}

	// the first argument must be a number.
	v, ok := args[0].(primitive.Number)
	if !ok {
		return primitive.Error(fmt.Sprintf("argument '%s' was not a number", args[0].ToString()))
	}

	// now process all the rest of the arguments
	for _, i := range args[1:] {

		// check we have a number
		n, ok := i.(primitive.Number)
		if ok {
			v -= n
		} else {
			return primitive.Error(fmt.Sprintf("argument %s was not a number", i.ToString()))
		}
	}
	return primitive.Number(v)
}

// multiplyFn implements "*"
func multiplyFn(args []primitive.Primitive) primitive.Primitive {
	// ensure we have at least one argument
	if len(args) < 1 {
		return primitive.Error("invalid argument count")
	}

	// the first argument must be a number.
	v, ok := args[0].(primitive.Number)
	if !ok {
		return primitive.Error(fmt.Sprintf("argument '%s' was not a number", args[0].ToString()))
	}

	// now process all the rest of the arguments
	for _, i := range args[1:] {

		// check we have a number
		n, ok := i.(primitive.Number)
		if ok {
			v *= n
		} else {
			return primitive.Error(fmt.Sprintf("argument %s was not a number", i.ToString()))
		}
	}
	return primitive.Number(v)
}

// divideFn implements "/"
func divideFn(args []primitive.Primitive) primitive.Primitive {
	// ensure we have at least one argument
	if len(args) < 1 {
		return primitive.Error("invalid argument count")
	}

	// the first argument must be a number.
	v, ok := args[0].(primitive.Number)
	if !ok {
		return primitive.Error(fmt.Sprintf("argument '%s' was not a number", args[0].ToString()))
	}

	// now process all the rest of the arguments
	for _, i := range args[1:] {

		// check we have a number
		n, ok := i.(primitive.Number)
		if ok {
			v /= n
		} else {
			return primitive.Error(fmt.Sprintf("argument %s was not a number", i.ToString()))
		}
	}
	return primitive.Number(v)
}

// modFn implements "%"
func modFn(args []primitive.Primitive) primitive.Primitive {
	if len(args) != 2 {
		return primitive.Error("wrong number of arguments")
	}
	if _, ok := args[0].(primitive.Number); !ok {
		return primitive.Error("argument not a number")
	}
	if _, ok := args[1].(primitive.Number); !ok {
		return primitive.Error("argument not a number")
	}
	return primitive.Number(int(args[0].(primitive.Number)) % int(args[1].(primitive.Number)))
}

// expnFn implements "#"
func expnFn(args []primitive.Primitive) primitive.Primitive {
	if len(args) != 2 {
		return primitive.Error("wrong number of arguments")
	}
	if _, ok := args[0].(primitive.Number); !ok {
		return primitive.Error("argument not a number")
	}
	if _, ok := args[1].(primitive.Number); !ok {
		return primitive.Error("argument not a number")
	}
	return primitive.Number(math.Pow(float64(args[0].(primitive.Number)), float64(args[1].(primitive.Number))))
}

// ltFn implements "<"
func ltFn(args []primitive.Primitive) primitive.Primitive {
	if len(args) != 2 {
		return primitive.Error("wrong number of arguments")
	}

	if _, ok := args[0].(primitive.Number); !ok {
		return primitive.Error("argument not a number")
	}
	if _, ok := args[1].(primitive.Number); !ok {
		return primitive.Error("argument not a number")
	}
	return primitive.Bool(args[0].(primitive.Number) < args[1].(primitive.Number))
}

// eqFn implements "eq"
func eqFn(args []primitive.Primitive) primitive.Primitive {
	if len(args) != 2 {
		return primitive.Error("wrong number of arguments")
	}

	a := args[0]
	b := args[1]

	if a.Type() != b.Type() {
		return primitive.Bool(false)
	}
	if a.ToString() != b.ToString() {
		return primitive.Bool(false)
	}
	return primitive.Bool(true)
}
