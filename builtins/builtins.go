// Package builtins contains the implementations of the lisp-callable
// functions which are implemented in golang.
package builtins

import (
	"fmt"

	"github.com/skx/yal/env"
	"github.com/skx/yal/primitive"
)

// PopulateEnvironment registers our default primitives
func PopulateEnvironment(env *env.Environment) {

	env.Set("+", &primitive.Procedure{F: func(args []primitive.Primitive) primitive.Primitive {
		v := args[0].(primitive.Number)
		for _, i := range args[1:] {
			v += i.(primitive.Number)
		}
		return primitive.Number(v)
	}})

	env.Set("-", &primitive.Procedure{F: func(args []primitive.Primitive) primitive.Primitive {
		v := args[0].(primitive.Number)
		for _, i := range args[1:] {
			v -= i.(primitive.Number)
		}
		return primitive.Number(v)
	}})

	env.Set("*", &primitive.Procedure{F: func(args []primitive.Primitive) primitive.Primitive {
		v := args[0].(primitive.Number)
		for _, i := range args[1:] {
			v *= i.(primitive.Number)
		}
		return primitive.Number(v)
	}})

	env.Set("/", &primitive.Procedure{F: func(args []primitive.Primitive) primitive.Primitive {
		v := args[0].(primitive.Number)
		for _, i := range args[1:] {
			v /= i.(primitive.Number)
		}
		return primitive.Number(v)
	}})

	// Only added "<" + ">"
	//
	// "<=" and ">=" can be implemented in lisp :)
	//
	env.Set("<", &primitive.Procedure{F: func(args []primitive.Primitive) primitive.Primitive {
		if _, ok := args[0].(primitive.Number); !ok {
			return primitive.Error("argument not a number")
		}
		if _, ok := args[1].(primitive.Number); !ok {
			return primitive.Error("argument not a number")
		}
		return primitive.Bool(args[0].(primitive.Number) < args[1].(primitive.Number))
	}})

	env.Set(">", &primitive.Procedure{F: func(args []primitive.Primitive) primitive.Primitive {
		if _, ok := args[0].(primitive.Number); !ok {
			return primitive.Error("argument not a number")
		}
		if _, ok := args[1].(primitive.Number); !ok {
			return primitive.Error("argument not a number")
		}
		return primitive.Bool(args[0].(primitive.Number) > args[1].(primitive.Number))
	}})

	env.Set("=", &primitive.Procedure{F: func(args []primitive.Primitive) primitive.Primitive {
		if _, ok := args[0].(primitive.Number); !ok {
			return primitive.Error("argument not a number")
		}
		if _, ok := args[1].(primitive.Number); !ok {
			return primitive.Error("argument not a number")
		}
		return primitive.Bool(args[0].(primitive.Number) == args[1].(primitive.Number))
	}})

	env.Set("%", &primitive.Procedure{F: func(args []primitive.Primitive) primitive.Primitive {
		if _, ok := args[0].(primitive.Number); !ok {
			return primitive.Error("argument not a number")
		}
		if _, ok := args[1].(primitive.Number); !ok {
			return primitive.Error("argument not a number")
		}
		return primitive.Number(int(args[0].(primitive.Number)) % int(args[1].(primitive.Number)))
	}})

	// primitive.List
	env.Set("list", &primitive.Procedure{F: func(args []primitive.Primitive) primitive.Primitive {
		return primitive.List(args)
	}})
	env.Set("car", &primitive.Procedure{F: func(args []primitive.Primitive) primitive.Primitive {
		return args[0].(primitive.List)[0]
	}})
	env.Set("cdr", &primitive.Procedure{F: func(args []primitive.Primitive) primitive.Primitive {
		return args[0].(primitive.List)[1:]
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

	// type
	env.Set("type", &primitive.Procedure{F: func(args []primitive.Primitive) primitive.Primitive {
		return primitive.String(args[0].Type())
	}})

	// equality
	env.Set("eq", &primitive.Procedure{F: func(args []primitive.Primitive) primitive.Primitive {
		a := args[0]
		b := args[1]

		if a.Type() != b.Type() {
			return primitive.Bool(false)
		}
		if a.ToString() != b.ToString() {
			return primitive.Bool(false)
		}
		return primitive.Bool(true)
	}})

	// Print
	env.Set("print", &primitive.Procedure{F: func(args []primitive.Primitive) primitive.Primitive {
		// no args
		if len(args) < 1 {
			return primitive.Nil{}
		}

		// one arg
		if len(args) == 1 {
			str := expandStr(args[0].ToString())

			fmt.Println(str)
			return primitive.Nil{}
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
		return primitive.Nil{}
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
