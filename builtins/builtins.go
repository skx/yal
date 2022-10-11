// Package builtins contains the implementations of the lisp-callable
// functions which are implemented in golang.
//
// Note that these builtings here don't have access to the run-time
// environment, which is why things like "gensym" have to be implemented
// in our core `eval.go` package.
//
// Updating the builtins to receive a references to the environment would
// allow some of the implementations to be moved.
package builtins

import (
	"fmt"
	"math"
	"os"
	"regexp"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/skx/yal/env"
	"github.com/skx/yal/primitive"
)

// PrimitiveFn is the type which represents a function signature for
// a lisp-usable function implemented in golang.
type PrimitiveFn func(args []primitive.Primitive) primitive.Primitive

// regCache is a cache of compiled regular expression objects.
// These may persist between runs because a regular expression object
// is essentially constant.
var regCache map[string]*regexp.Regexp

// init ensures that our regexp cache is populated
func init() {
	regCache = make(map[string]*regexp.Regexp)
}

// PopulateEnvironment registers our default primitives
func PopulateEnvironment(env *env.Environment) {

	// maths
	env.Set("#", &primitive.Procedure{F: expnFn})
	env.Set("%", &primitive.Procedure{F: modFn})
	env.Set("*", &primitive.Procedure{F: multiplyFn})
	env.Set("+", &primitive.Procedure{F: plusFn})
	env.Set("-", &primitive.Procedure{F: minusFn})
	env.Set("/", &primitive.Procedure{F: divideFn})

	// comparisons
	//
	// When it comes to comparisons there are several we could
	// use:
	//
	//  <
	//  <=
	//  >
	//  >=
	//
	// We only actually need to implement "<" in Golang, the rest
	// can be added in lisp.
	env.Set("<", &primitive.Procedure{F: ltFn})
	env.Set("=", &primitive.Procedure{F: equalsFn, Help: `= returns true if supplied with two numerical values, of equal value.
\tSee also: eq
\t Example: (print (= 3 a))`})
	env.Set("eq", &primitive.Procedure{F: eqFn, Help: `eq returns true if the two values supplied as parameters have the same type, and string representation.

\tSee also: =
\t Example: (print (eq "bob" 2))`})

	// Types
	env.Set("nil?", &primitive.Procedure{F: nilFn, Help: `nil? returns a true value if the specified argument is nil, or an empty list.`})
	env.Set("type", &primitive.Procedure{F: typeFn, Help: `type returns a string describing the type of the specified object.

Example:  (type "string") (type 3) (type type)`})

	// List
	env.Set("car", &primitive.Procedure{F: carFn, Help: `car returns the first item from the specified list.`})
	env.Set("cdr", &primitive.Procedure{F: cdrFn, Help: `cdr returns all items from the specified list, except the first.`})
	env.Set("cons", &primitive.Procedure{F: consFn, Help: `cons joins the two specified lists: FIXME`})
	env.Set("join", &primitive.Procedure{F: joinFn, Help: `join returns a string formed by converting every element of the supplied list into a string and concatenating them.`})
	env.Set("list", &primitive.Procedure{F: listFn, Help: `list creates and returns a list containing each of the specified arguments, in order.`})

	// Hash
	env.Set("contains?", &primitive.Procedure{F: containsFn, Help: `contains? returns true if the hash specified as the first argument contains the key specified as the second argument.`})
	env.Set("get", &primitive.Procedure{F: getFn, Help: `get returns the specified field from the specified hash.

\tSee also: set
\t Example: (get {:name "steve" :location "Europe" } ":name")`})
	env.Set("keys", &primitive.Procedure{F: keysFn, Help: `keys returns the keys which are present in the specified hash.

NOTE: Keys are returned in sorted order.

\tSee also: vals`})
	env.Set("set", &primitive.Procedure{F: setFn, Help: `set updates the specified hash, setting the value given by name.

\t: See also: get
\t: Example:  (set! person {:name "Steve"})  (set person :name "Bobby")`})
	env.Set("vals", &primitive.Procedure{F: valsFn, Help: `valus returns the values which are present in the specified hash.

NOTE: Values are returned in the order of their sorted keys.

\tSee also: keys`})

	// core
	env.Set("arch", &primitive.Procedure{F: archFn, Help: `arch returns a simple string describing the architecture the current host is running upon.

\tSee also: (os)
\t Example: (print (arch))`})
	env.Set("date", &primitive.Procedure{F: dateFn, Help: `date returns a list containing date-related fields; the day of the week, the day-number, the month-number, and the year.

\tSee also: (time)`})
	env.Set("error", &primitive.Procedure{F: errorFn, Help: `error raises an error with the specified argument as the explaination.

\t Example: (error "Expected foo to be bar!")`})
	env.Set("getenv", &primitive.Procedure{F: getenvFn, Help: `getenv returns the contents of the environmental-variable which was specified as the first argument.

\t Example: (print (getenv "HOME"))`})
	env.Set("ms", &primitive.Procedure{F: msFn, Help: `ms returns the current time as a number of milliseconds, it is useful for benchmarking.

\tSee also: now`})
	env.Set("now", &primitive.Procedure{F: nowFn, Help: `now returns the number of seconds since the Unix Epoch.

\tSee also: ms`})
	env.Set("os", &primitive.Procedure{F: osFn, Help: `os returns a simple string describing the operating system the current host is running.

\tSee also: (arch)
\t Example: (print (os))`,
	})
	env.Set("print", &primitive.Procedure{F: printFn, Help: `print is used to output text to the console.  It can be called with either an object/string to print, or a format-string and list of parameters.

\tSee also: sprintf
\t Example: (print "Hello, world")
\t Example: (print "Hello user %s you are %d" (getenv "USER") 32)`})
	env.Set("sort", &primitive.Procedure{F: sortFn, Help: `sort will sort the items in the list specified as the single argument, and return them as a new list.

\t Example: (print (sort 3 43  1 "Steve" "Adam"))
`})
	env.Set("sprintf", &primitive.Procedure{F: sprintfFn, Help: `sprintf allows formating values with a simple format-string.

\tSee also: print
\t Example: (sprintf "Today is %s" (weekday))`})
	env.Set("slurp", &primitive.Procedure{F: slurpFn, Help: `slurp returns the contents of the specified file.`})
	env.Set("time", &primitive.Procedure{F: timeFn, Help: `time returns a list containing time-related entries; the current hour, the current minute past the hour, and the current value of the seconds.

\tSee also: (date)`})

	// string
	env.Set("chr", &primitive.Procedure{F: chrFn, Help: `chr returns a string containing the single character who's ASCII code was provided.

\tSee also: ord
\t Example: (chr 42) ; => "*"`})

	env.Set("match", &primitive.Procedure{F: matchFn, Help: `match is used to perform regular expression matches.  The first parameter must be a suitable regular expression, supplied in string-form, and the second should be a value to test against.  If the second value is not a string it will be stringified prior to the test-attempt.

Any matches found will be returned as a list, with nil being returned on no match.

\t Example: (print (match "c.ke$" "cake"))`})
	env.Set("ord", &primitive.Procedure{F: ordFn, Help: `ord returns the ASCII code for the character provided as the first input.

\tSee also: chr
\t Example: (ord "a") ; => 97`})

	env.Set("split", &primitive.Procedure{F: splitFn, Help: `split accepts two string parameters, and splits the first string by the term specified as the second argument, returning a list of the results.

\tSee also: join
\t Example: (split "steve" "e") ; => ("st" "v")
\t Example: (split "steve" "")  ; => ("s" "t" "e" "v" "e")`})

	env.Set("str", &primitive.Procedure{F: strFn, Help: `str converts the parameter supplied to a string, and returns it.`})
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
			case 'e':
				out += string(rune(033))
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

// minusFn implements "-"
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
			if n == 0 {
				return primitive.Error("attempted division by zero")
			}

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

	a := int(args[0].(primitive.Number))
	b := int(args[1].(primitive.Number))
	if b == 0 {
		return primitive.Error("attempted division by zero")
	}
	return primitive.Number(a % b)
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

// equalsFn implements "="
func equalsFn(args []primitive.Primitive) primitive.Primitive {
	if len(args) != 2 {
		return primitive.Error("wrong number of arguments")
	}

	a := args[0]
	b := args[1]

	if a.Type() != "number" {
		return primitive.Error("argument was not a number")
	}
	if b.Type() != "number" {
		return primitive.Error("argument was not a number")
	}
	if a.(primitive.Number) == b.(primitive.Number) {
		return primitive.Bool(true)
	}
	return primitive.Bool(false)
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

// listFn implements "list"
func listFn(args []primitive.Primitive) primitive.Primitive {
	return primitive.List(args)
}

// carFn implements "car"
func carFn(args []primitive.Primitive) primitive.Primitive {

	if len(args) != 1 {
		return primitive.Error("wrong number of arguments")
	}

	// ensure we received a list
	if _, ok := args[0].(primitive.List); !ok {
		return primitive.Error("argument not a list")
	}

	// If we have at least one entry then return the first
	lst := args[0].(primitive.List)
	if len(lst) > 0 {
		return lst[0]
	}

	// Otherwise return nil
	return primitive.Nil{}
}

// cdrFn implements "cdr"
func cdrFn(args []primitive.Primitive) primitive.Primitive {

	if len(args) != 1 {
		return primitive.Error("wrong number of arguments")
	}

	// ensure we received a list
	if _, ok := args[0].(primitive.List); !ok {
		return primitive.Error("argument not a list")
	}

	lst := args[0].(primitive.List)
	if len(lst) > 0 {
		return lst[1:]
	}
	return primitive.Nil{}
}

// errorFn implements "error"
func errorFn(args []primitive.Primitive) primitive.Primitive {
	if len(args) != 1 {
		return primitive.Error("wrong number of arguments")
	}
	return primitive.Error(args[0].ToString())
}

// typeFn implements "type"
func typeFn(args []primitive.Primitive) primitive.Primitive {
	if len(args) != 1 {
		return primitive.Error("wrong number of arguments")
	}
	return primitive.String(args[0].Type())
}

// slurpFn returns the contents of the specified file
func slurpFn(args []primitive.Primitive) primitive.Primitive {
	if len(args) != 1 {
		return primitive.Error("wrong number of arguments")
	}

	fName := args[0].ToString()
	data, err := os.ReadFile(fName)
	if err != nil {
		return primitive.Error(fmt.Sprintf("error reading %s %s", fName, err))
	}
	return primitive.String(string(data))
}

// strFn implements "str"
func strFn(args []primitive.Primitive) primitive.Primitive {
	if len(args) != 1 {
		return primitive.Error("wrong number of arguments")
	}
	return primitive.String(args[0].ToString())
}

// nilFn implements nil?
func nilFn(args []primitive.Primitive) primitive.Primitive {
	if len(args) != 1 {
		return primitive.Error("wrong number of arguments")
	}

	// nil is nil (yeah, really)
	if primitive.IsNil(args[0]) {
		return primitive.Bool(true)
	}

	// an empty list is nil.
	if list, ok := args[0].(primitive.List); ok {
		return primitive.Bool(len(list) == 0)
	}
	return primitive.Bool(false)

}

// consFn implements (cons).
func consFn(args []primitive.Primitive) primitive.Primitive {
	if len(args) < 1 {
		return primitive.Error("wrong number of arguments")
	}

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
}

// osFn implements (os)
func osFn(args []primitive.Primitive) primitive.Primitive {
	return primitive.String(runtime.GOOS)
}

// archFn implements (os)
func archFn(args []primitive.Primitive) primitive.Primitive {
	return primitive.String(runtime.GOARCH)
}

// dateFn returns the current (Weekday, DD, MM, YYYY) as a list.
func dateFn(args []primitive.Primitive) primitive.Primitive {
	var ret primitive.List

	t := time.Now()

	name := t.Weekday().String()
	day := t.Day()
	mon := int(t.Month())
	year := t.Year()

	ret = append(ret, primitive.String(name))
	ret = append(ret, primitive.Number(day))
	ret = append(ret, primitive.Number(mon))
	ret = append(ret, primitive.Number(year))

	return ret
}

// printFn implements (print).
func printFn(args []primitive.Primitive) primitive.Primitive {
	// no args
	if len(args) < 1 {
		return primitive.Error("wrong number of arguments")
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
}

// (sprintf "fmt" "arg1" ... "argN")
func sprintfFn(args []primitive.Primitive) primitive.Primitive {

	// we need 2+ arguments
	if len(args) < 2 {
		return primitive.Error("wrong number of arguments")
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
}

// (split "str" "by")
func splitFn(args []primitive.Primitive) primitive.Primitive {

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
}

// timeFn returns the current (HH, MM, SS) as a list.
func timeFn(args []primitive.Primitive) primitive.Primitive {
	var ret primitive.List

	t := time.Now()

	hr := t.Hour()
	mn := t.Minute()
	sc := t.Second()

	ret = append(ret, primitive.Number(hr))
	ret = append(ret, primitive.Number(mn))
	ret = append(ret, primitive.Number(sc))

	return ret
}

// (join (1 2 3)
func joinFn(args []primitive.Primitive) primitive.Primitive {

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
}

// sortFn implements (sort)
func sortFn(args []primitive.Primitive) primitive.Primitive {
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

}

// getFn is the implementation of `(get hash key)`
func getFn(args []primitive.Primitive) primitive.Primitive {

	// We need two arguments
	if len(args) != 2 {
		return primitive.Error("invalid argument count")
	}

	// First is a Hash
	if _, ok := args[0].(primitive.Hash); !ok {
		return primitive.Error("argument not a hash")
	}

	tmp := args[0].(primitive.Hash)
	return tmp.Get(args[1].ToString())
}

// keysFn is the implementation of `(keys hash)`
func keysFn(args []primitive.Primitive) primitive.Primitive {

	// We need a single argument
	if len(args) != 1 {
		return primitive.Error("invalid argument count")
	}

	// First is a Hash
	if _, ok := args[0].(primitive.Hash); !ok {
		return primitive.Error("argument not a hash")
	}

	// Create the list to hold the result
	var c primitive.List

	// Cast the argument
	tmp := args[0].(primitive.Hash)

	// Get the keys as a list
	keys := []string{}

	// Add the keys
	for x := range tmp.Entries {
		keys = append(keys, x)
	}

	// Sort the list
	sort.Strings(keys)

	// Now append
	for _, x := range keys {
		c = append(c, primitive.String(x))
	}

	return c
}

// valsFn is the implementation of `(vals hash)`
func valsFn(args []primitive.Primitive) primitive.Primitive {

	// We need a single argument
	if len(args) != 1 {
		return primitive.Error("invalid argument count")
	}

	// First is a Hash
	if _, ok := args[0].(primitive.Hash); !ok {
		return primitive.Error("argument not a hash")
	}

	// Create the list to hold the result
	var c primitive.List

	// Cast the argument
	tmp := args[0].(primitive.Hash)

	// Get the keys as a list
	keys := []string{}

	// Add the keys
	for x := range tmp.Entries {
		keys = append(keys, x)
	}

	// Sort the list
	sort.Strings(keys)

	// Now append the value
	for _, x := range keys {
		c = append(c, tmp.Entries[x])
	}

	return c
}

// containsFn implements (contains?)
func containsFn(args []primitive.Primitive) primitive.Primitive {

	// We need a pair of arguments
	if len(args) != 2 {
		return primitive.Error("invalid argument count")
	}

	// First is a Hash
	hsh, ok := args[0].(primitive.Hash)
	if !ok {
		return primitive.Error("argument not a hash")
	}

	// The second should be a string, but other things can be converted
	str, ok := args[1].(primitive.String)
	if !ok {
		str = primitive.String(args[1].ToString())
	}

	_, found := hsh.Entries[str.ToString()]
	if found {
		return primitive.Bool(true)
	}

	return primitive.Bool(false)

}

// setFn is the implementation of `(set hash key val)`
func setFn(args []primitive.Primitive) primitive.Primitive {

	// We need three arguments
	if len(args) != 3 {
		return primitive.Error("invalid argument count")
	}

	// First is a Hash
	if _, ok := args[0].(primitive.Hash); !ok {
		return primitive.Error("argument not a hash")
	}

	tmp := args[0].(primitive.Hash)
	tmp.Set(args[1].ToString(), args[2])
	return args[2]
}

// getenvFn is the implementation of `(getenv "PATH")`
func getenvFn(args []primitive.Primitive) primitive.Primitive {

	// If we have only a single argument
	if len(args) != 1 {
		return primitive.Error("invalid argument count")
	}

	// Which is a string
	if _, ok := args[0].(primitive.String); !ok {
		return primitive.Error("argument not a string")
	}

	// Return the value
	str := args[0].(primitive.String)
	return primitive.String(os.Getenv(string(str)))
}

// nowFn is the implementation of `(now)`
func nowFn(args []primitive.Primitive) primitive.Primitive {
	return primitive.Number(time.Now().Unix())
}

// msFn is the implementation of `(ms)`
func msFn(args []primitive.Primitive) primitive.Primitive {
	return primitive.Number(time.Now().UnixNano() / int64(time.Millisecond))
}

// matchFn is the implementation of (match ..)
func matchFn(args []primitive.Primitive) primitive.Primitive {

	// We need two arguments
	if len(args) != 2 {
		return primitive.Error("invalid argument count")
	}

	// First argument is a string (which is a regexp)
	if _, ok := args[0].(primitive.String); !ok {
		return primitive.Error("argument not a string")
	}

	// Second is what we'll match
	pat := args[0].ToString()
	txt := args[1].ToString()

	// Look for a cached regexp
	r, ok := regCache[pat]
	if !ok {
		// OK it wasn't found, so compile it.
		var err error
		r, err = regexp.Compile(pat)

		// Ensure it compiled
		if err != nil {
			return primitive.Error(fmt.Sprintf("failed to compile regexp '%s':%s", pat, err.Error()))
		}

		// store in the cache for next time
		regCache[pat] = r
	}

	res := r.FindStringSubmatch(txt)

	if len(res) > 0 {

		// Return the items in a list
		var tmp primitive.List

		if len(res) > 0 {
			for i := 0; i < len(res); i++ {

				tmp = append(tmp, primitive.String(res[i]))
			}
		}

		return tmp
	}

	// No match
	return primitive.Nil{}

}

// chrFn is the implementation of (chr ..)
func chrFn(args []primitive.Primitive) primitive.Primitive {

	if len(args) != 1 {
		return primitive.Error("wrong number of arguments")
	}

	if _, ok := args[0].(primitive.Number); !ok {
		return primitive.Error("argument not a number")
	}

	i := args[0].(primitive.Number)
	rune := rune(i)

	return primitive.String(rune)
}

// ordFn is the implementation of (ord ..)
func ordFn(args []primitive.Primitive) primitive.Primitive {

	if len(args) != 1 {
		return primitive.Error("wrong number of arguments")
	}

	if _, ok := args[0].(primitive.String); !ok {
		return primitive.Error("argument not a string")
	}

	// We convert this to an array of runes because we
	// want to handle unicode strings.
	i := []rune(args[0].ToString())

	if len(i) > 0 {
		s := rune(i[0])
		return primitive.Number(float64(rune(s)))
	}
	return primitive.Number(0)
}
