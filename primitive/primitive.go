// Package primitive contains the definitions of our primitive types,
// which are "nil", "bool", "number", "string", and "list".
package primitive

import (
	"fmt"

	"github.com/skx/yal/env"
)

// Primitive is the interface of all our types
type Primitive interface {

	// Convert this primitive to a string
	ToString() string

	// Return the type of this object
	Type() string
}

// IsNil tests whether an expression is nil.
func IsNil(e Primitive) bool {
	var n Nil
	return e == n
}

// PopulateEnvironment registers our default primitives
func PopulateEnvironment(env *env.Environment) {

	env.Set("+", &Procedure{F: func(args []Primitive) Primitive {
		v := args[0].(Number)
		for _, i := range args[1:] {
			v += i.(Number)
		}
		return Number(v)
	}})

	env.Set("-", &Procedure{F: func(args []Primitive) Primitive {
		v := args[0].(Number)
		for _, i := range args[1:] {
			v -= i.(Number)
		}
		return Number(v)
	}})

	env.Set("*", &Procedure{F: func(args []Primitive) Primitive {
		v := args[0].(Number)
		for _, i := range args[1:] {
			v *= i.(Number)
		}
		return Number(v)
	}})

	env.Set("/", &Procedure{F: func(args []Primitive) Primitive {
		v := args[0].(Number)
		for _, i := range args[1:] {
			v /= i.(Number)
		}
		return Number(v)
	}})

	// Only added "<" + ">"
	//
	// "<=" and ">=" can be implemented in lisp :)
	//
	env.Set("<", &Procedure{F: func(args []Primitive) Primitive {
		if _, ok := args[0].(Number); !ok {
			return Error("argument not a number")
		}
		if _, ok := args[1].(Number); !ok {
			return Error("argument not a number")
		}
		return Bool(args[0].(Number) < args[1].(Number))
	}})

	env.Set(">", &Procedure{F: func(args []Primitive) Primitive {
		if _, ok := args[0].(Number); !ok {
			return Error("argument not a number")
		}
		if _, ok := args[1].(Number); !ok {
			return Error("argument not a number")
		}
		return Bool(args[0].(Number) > args[1].(Number))
	}})

	env.Set("=", &Procedure{F: func(args []Primitive) Primitive {
		if _, ok := args[0].(Number); !ok {
			return Error("argument not a number")
		}
		if _, ok := args[1].(Number); !ok {
			return Error("argument not a number")
		}
		return Bool(args[0].(Number) == args[1].(Number))
	}})

	env.Set("%", &Procedure{F: func(args []Primitive) Primitive {
		if _, ok := args[0].(Number); !ok {
			return Error("argument not a number")
		}
		if _, ok := args[1].(Number); !ok {
			return Error("argument not a number")
		}
		return Number(int(args[0].(Number)) % int(args[1].(Number)))
	}})

	// List
	env.Set("list", &Procedure{F: func(args []Primitive) Primitive {
		return List(args)
	}})
	env.Set("car", &Procedure{F: func(args []Primitive) Primitive {
		return args[0].(List)[0]
	}})
	env.Set("cdr", &Procedure{F: func(args []Primitive) Primitive {
		return args[0].(List)[1:]
	}})

	// nil
	env.Set("nil?", &Procedure{F: func(args []Primitive) Primitive {
		// nil is nil (yeah, really)
		if IsNil(args[0]) {
			return Bool(true)
		}

		// an empty list is nil.
		if list, ok := args[0].(List); ok {
			return Bool(len(list) == 0)
		}
		return Bool(false)
	}})

	env.Set("cons", &Procedure{F: func(args []Primitive) Primitive {
		if len(args) == 1 {
			return List{args[0]}
		}
		if args[1] == nil || IsNil(args[1]) {
			return List{args[0]}
		}
		if _, ok := args[1].(List); ok {
			return append(List{args[0]}, args[1].(List)...)
		}
		return List{args[0], args[1]}
	}})

	// type
	env.Set("type", &Procedure{F: func(args []Primitive) Primitive {
		return String(args[0].Type())
	}})

	// equality
	env.Set("eq", &Procedure{F: func(args []Primitive) Primitive {
		a := args[0]
		b := args[1]

		if a.Type() != b.Type() {
			return Bool(false)
		}
		if a.ToString() != b.ToString() {
			return Bool(false)
		}
		return Bool(true)
	}})

	// Print
	env.Set("print", &Procedure{F: func(args []Primitive) Primitive {
		// no args
		if len(args) < 1 {
			return Nil{}
		}

		// one arg
		if len(args) == 1 {
			str := expandStr(args[0].ToString())

			fmt.Println(str)
			return Nil{}
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
		fmt.Println(fmt.Sprintf(frmt, parm...))
		return Nil{}
	}})

	// Convert an object to a string
	env.Set("str", &Procedure{F: func(args []Primitive) Primitive {
		return String(args[0].ToString())
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
		if c == '\\' && i < l {

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
