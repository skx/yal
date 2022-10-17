// Package builtins contains the implementations of the lisp-callable
// functions which are implemented in golang.
//
// This package excludes the special forms, which have to be handled
// specially - for example "(let*)", "(if)", and "(eval..)".
package builtins

import (
	"bufio"
	"bytes"
	"fmt"
	"math"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/skx/yal/env"
	"github.com/skx/yal/primitive"

	_ "embed" // embedded-resource magic
)

//go:embed help.txt
var helpText string

// helpMap contains a map of help-text.
//
// it is populated at init-time, from helpText
var helpMap map[string]string

// regCache is a cache of compiled regular expression objects.
// These may persist between runs because a regular expression object
// is essentially constant.
var regCache map[string]*regexp.Regexp

// symCount is the count of symbols generated by the 'gensym' built-in
// function.
var symCount int

// init ensures that our regexp cache is populated, and that we build
// up a list of help-texts, keyed on function name.
func init() {

	// Create our maps.
	regCache = make(map[string]*regexp.Regexp)
	helpMap = make(map[string]string)

	// Convert the help-text to a string
	help := string(helpText)

	term := ""
	text := ""

	// process the help-text, embedded, line by line
	for _, line := range strings.Split(help, "\n") {

		// end of an entry?  Save it away
		if line == "%%" {
			helpMap[term] = text

			term = ""
			text = ""
			continue
		}
		if len(term) == 0 {
			// no term?  Then save one
			term = line
		} else {
			// Otherwise add the text to it.
			if len(text) > 0 {
				text += "\n"
			}
			text += line
		}
	}

	// All done, our help text should be available at run-time now.
}

// PopulateEnvironment registers our default primitives
func PopulateEnvironment(env *env.Environment) {

	//
	// bind the functions - sorted order
	//
	env.Set("#", &primitive.Procedure{F: expnFn, Help: helpMap["#"], Args: []primitive.Symbol{primitive.Symbol("N"), primitive.Symbol("b")}})
	env.Set("%", &primitive.Procedure{F: modFn, Help: helpMap["%"], Args: []primitive.Symbol{primitive.Symbol("a"), primitive.Symbol("b")}})
	env.Set("*", &primitive.Procedure{F: multiplyFn, Help: helpMap["*"], Args: []primitive.Symbol{primitive.Symbol("N"), primitive.Symbol("arg1..argN")}})
	env.Set("+", &primitive.Procedure{F: plusFn, Help: helpMap["+"], Args: []primitive.Symbol{primitive.Symbol("N"), primitive.Symbol("arg1..argN")}})
	env.Set("-", &primitive.Procedure{F: minusFn, Help: helpMap["-"], Args: []primitive.Symbol{primitive.Symbol("N"), primitive.Symbol("arg1..argN")}})
	env.Set("/", &primitive.Procedure{F: divideFn, Help: helpMap["/"], Args: []primitive.Symbol{primitive.Symbol("N"), primitive.Symbol("arg1..argN")}})
	env.Set("/=", &primitive.Procedure{F: inequalityFn, Help: helpMap["/="], Args: []primitive.Symbol{primitive.Symbol("N"), primitive.Symbol("arg1..argN")}})
	env.Set("<", &primitive.Procedure{F: ltFn, Help: helpMap["<"], Args: []primitive.Symbol{primitive.Symbol("a"), primitive.Symbol("b")}})
	env.Set("=", &primitive.Procedure{F: equalsFn, Help: helpMap["="], Args: []primitive.Symbol{primitive.Symbol("arg1"), primitive.Symbol("arg2 .. argN")}})
	env.Set("arch", &primitive.Procedure{F: archFn, Help: helpMap["arch"]})
	env.Set("car", &primitive.Procedure{F: carFn, Help: helpMap["car"], Args: []primitive.Symbol{primitive.Symbol("list")}})
	env.Set("cdr", &primitive.Procedure{F: cdrFn, Help: helpMap["cdr"], Args: []primitive.Symbol{primitive.Symbol("list")}})
	env.Set("chr", &primitive.Procedure{F: chrFn, Help: helpMap["chr"], Args: []primitive.Symbol{primitive.Symbol("num")}})
	env.Set("cons", &primitive.Procedure{F: consFn, Help: helpMap["cons"], Args: []primitive.Symbol{primitive.Symbol("a"), primitive.Symbol("b")}})
	env.Set("contains?", &primitive.Procedure{F: containsFn, Help: helpMap["contains?"], Args: []primitive.Symbol{primitive.Symbol("hash"), primitive.Symbol("key")}})
	env.Set("date", &primitive.Procedure{F: dateFn, Help: helpMap["date"]})
	env.Set("directory:entries", &primitive.Procedure{F: directoryEntriesFn, Help: helpMap["directory:entries"]})
	env.Set("directory?", &primitive.Procedure{F: directoryFn, Help: helpMap["directory?"], Args: []primitive.Symbol{primitive.Symbol("path")}})
	env.Set("eq", &primitive.Procedure{F: eqFn, Help: helpMap["eq"], Args: []primitive.Symbol{primitive.Symbol("a"), primitive.Symbol("b")}})
	env.Set("error", &primitive.Procedure{F: errorFn, Help: helpMap["error"], Args: []primitive.Symbol{primitive.Symbol("message")}})
	env.Set("exists?", &primitive.Procedure{F: existsFn, Help: helpMap["exists?"], Args: []primitive.Symbol{primitive.Symbol("path")}})
	env.Set("file:lines", &primitive.Procedure{F: fileLinesFn, Help: helpMap["file:lines"], Args: []primitive.Symbol{primitive.Symbol("path")}})
	env.Set("file:read", &primitive.Procedure{F: fileReadFn, Help: helpMap["file:read"], Args: []primitive.Symbol{primitive.Symbol("path")}})
	env.Set("file:stat", &primitive.Procedure{F: fileStatFn, Help: helpMap["file:stat"], Args: []primitive.Symbol{primitive.Symbol("path")}})
	env.Set("file:write", &primitive.Procedure{F: fileWriteFn, Help: helpMap["file:write"], Args: []primitive.Symbol{primitive.Symbol("path"), primitive.Symbol("content")}})
	env.Set("file?", &primitive.Procedure{F: fileFn, Help: helpMap["file?"], Args: []primitive.Symbol{primitive.Symbol("path")}})
	env.Set("gensym", &primitive.Procedure{F: gensymFn, Help: helpMap["gensym"]})
	env.Set("get", &primitive.Procedure{F: getFn, Help: helpMap["get"], Args: []primitive.Symbol{primitive.Symbol("hash"), primitive.Symbol("key")}})
	env.Set("getenv", &primitive.Procedure{F: getenvFn, Help: helpMap["getenv"], Args: []primitive.Symbol{primitive.Symbol("key")}})
	env.Set("glob", &primitive.Procedure{F: globFn, Help: helpMap["glob"], Args: []primitive.Symbol{primitive.Symbol("pattern")}})
	env.Set("help", &primitive.Procedure{F: helpFn, Help: helpMap["help"], Args: []primitive.Symbol{primitive.Symbol("function")}})
	env.Set("join", &primitive.Procedure{F: joinFn, Help: helpMap["join"], Args: []primitive.Symbol{primitive.Symbol("list")}})
	env.Set("keys", &primitive.Procedure{F: keysFn, Help: helpMap["keys"], Args: []primitive.Symbol{primitive.Symbol("hash")}})
	env.Set("list", &primitive.Procedure{F: listFn, Help: helpMap["list"], Args: []primitive.Symbol{primitive.Symbol("arg1"), primitive.Symbol("arg...")}})
	env.Set("match", &primitive.Procedure{F: matchFn, Help: helpMap["match"], Args: []primitive.Symbol{primitive.Symbol("regexp"), primitive.Symbol("str")}})
	env.Set("ms", &primitive.Procedure{F: msFn, Help: helpMap["ms"]})
	env.Set("nil?", &primitive.Procedure{F: nilFn, Help: helpMap["nil?"], Args: []primitive.Symbol{primitive.Symbol("object")}})
	env.Set("now", &primitive.Procedure{F: nowFn, Help: helpMap["now"]})
	env.Set("ord", &primitive.Procedure{F: ordFn, Help: helpMap["ord"], Args: []primitive.Symbol{primitive.Symbol("char")}})
	env.Set("os", &primitive.Procedure{F: osFn, Help: helpMap["os"]})
	env.Set("print", &primitive.Procedure{F: printFn, Help: helpMap["print"], Args: []primitive.Symbol{primitive.Symbol("arg1..argN")}})
	env.Set("random", &primitive.Procedure{F: randomFn, Help: helpMap["random"], Args: []primitive.Symbol{primitive.Symbol("max")}})
	env.Set("set", &primitive.Procedure{F: setFn, Help: helpMap["set"], Args: []primitive.Symbol{primitive.Symbol("hash"), primitive.Symbol("key"), primitive.Symbol("val")}})
	env.Set("shell", &primitive.Procedure{F: shellFn, Help: helpMap["shell"], Args: []primitive.Symbol{primitive.Symbol("list")}})
	env.Set("sort", &primitive.Procedure{F: sortFn, Help: helpMap["sort"], Args: []primitive.Symbol{primitive.Symbol("list")}})
	env.Set("split", &primitive.Procedure{F: splitFn, Help: helpMap["split"], Args: []primitive.Symbol{primitive.Symbol("str"), primitive.Symbol("by")}})
	env.Set("sprintf", &primitive.Procedure{F: sprintfFn, Help: helpMap["sprintf"], Args: []primitive.Symbol{primitive.Symbol("arg1..argN")}})
	env.Set("str", &primitive.Procedure{F: strFn, Help: helpMap["str"], Args: []primitive.Symbol{primitive.Symbol("object")}})
	env.Set("time", &primitive.Procedure{F: timeFn, Help: helpMap["time"]})
	env.Set("type", &primitive.Procedure{F: typeFn, Help: helpMap["type"], Args: []primitive.Symbol{primitive.Symbol("object")}})
	env.Set("vals", &primitive.Procedure{F: valsFn, Help: helpMap["vals"], Args: []primitive.Symbol{primitive.Symbol("hash")}})
}

// Built in functions

// archFn implements (os)
func archFn(env *env.Environment, args []primitive.Primitive) primitive.Primitive {
	return primitive.String(runtime.GOARCH)
}

// carFn implements "car"
func carFn(env *env.Environment, args []primitive.Primitive) primitive.Primitive {

	if len(args) != 1 {
		return primitive.ArityError()
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
func cdrFn(env *env.Environment, args []primitive.Primitive) primitive.Primitive {

	if len(args) != 1 {
		return primitive.ArityError()
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

// chrFn is the implementation of (chr ..)
func chrFn(env *env.Environment, args []primitive.Primitive) primitive.Primitive {

	if len(args) != 1 {
		return primitive.ArityError()
	}

	if _, ok := args[0].(primitive.Number); !ok {
		return primitive.Error("argument not a number")
	}

	i := args[0].(primitive.Number)
	rune := rune(i)

	return primitive.String(rune)
}

// consFn implements (cons).
func consFn(env *env.Environment, args []primitive.Primitive) primitive.Primitive {
	if len(args) < 1 {
		return primitive.ArityError()
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

// containsFn implements (contains?)
func containsFn(env *env.Environment, args []primitive.Primitive) primitive.Primitive {

	// We need a pair of arguments
	if len(args) != 2 {
		return primitive.ArityError()
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

// dateFn returns the current (Weekday, DD, MM, YYYY) as a list.
func dateFn(env *env.Environment, args []primitive.Primitive) primitive.Primitive {
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

// directoryEntriesFn returns the files beneath given path, recursively.
func directoryEntriesFn(env *env.Environment, args []primitive.Primitive) primitive.Primitive {

	// We only need a single argument
	if len(args) != 1 {
		return primitive.ArityError()
	}

	// Which is a string
	pth, ok := args[0].(primitive.String)
	if !ok {
		return primitive.Error("argument not a string")
	}

	var res primitive.List

	_ = filepath.Walk(pth.ToString(), func(path string, info os.FileInfo, err error) error {

		if err != nil {
			return nil
		}

		res = append(res, primitive.String(path))
		return nil
	})

	return res
}

// directoryFn returns whether the given path exists, and is a directory
func directoryFn(env *env.Environment, args []primitive.Primitive) primitive.Primitive {

	// We only need a single argument
	if len(args) != 1 {
		return primitive.ArityError()
	}

	// Which is a string
	str, ok := args[0].(primitive.String)
	if !ok {
		return primitive.Error("argument not a string")
	}

	// Stat the entry
	info, err := os.Stat(str.ToString())

	// No error and isDir then true?  Otherwise false
	//
	// i.e. swallow errors
	if err == nil {
		if info.IsDir() {
			return primitive.Bool(true)
		}
	}
	return primitive.Bool(false)
}

// divideFn implements "/"
func divideFn(env *env.Environment, args []primitive.Primitive) primitive.Primitive {
	// ensure we have at least one argument
	if len(args) < 1 {
		return primitive.ArityError()
	}

	// the first argument must be a number.
	v, ok := args[0].(primitive.Number)
	if !ok {
		return primitive.Error(fmt.Sprintf("argument '%s' was not a number", args[0].ToString()))
	}

	// If there is only one argument then we return the
	// reciprocal.
	//
	// (i.e. "(/ 3)" == ".3333"
	if len(args) == 1 {
		return primitive.Number(1 / v)
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

// eqFn implements "eq"
func eqFn(env *env.Environment, args []primitive.Primitive) primitive.Primitive {
	if len(args) != 2 {
		return primitive.ArityError()
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

// equalsFn implements "="
func equalsFn(env *env.Environment, args []primitive.Primitive) primitive.Primitive {

	// We need at least two arguments
	if len(args) < 2 {
		return primitive.ArityError()
	}

	// First argument must be a number.
	nA, ok := args[0].(primitive.Number)
	if !ok {
		return primitive.Error("argument was not a number")
	}

	// Now we'll loop over all other numbers
	//
	// If we got something that was NOT the same as our
	// initial number we can terminate early but we don't
	// because it is important to also report on failures to
	// validate types - which we can't do if we bail.
	//
	ret := primitive.Bool(true)

	for _, i := range args[1:] {

		// check we have a number
		nB, ok2 := i.(primitive.Number)

		if !ok2 {
			return primitive.Error("argument was not a number")
		}

		// Record our failure, but keep testing in case
		// we have a type violation to report in a later
		// argument.
		if nB != nA {
			ret = primitive.Bool(false)
		}
	}

	return ret
}

// errorFn implements "error"
func errorFn(env *env.Environment, args []primitive.Primitive) primitive.Primitive {
	if len(args) != 1 {
		return primitive.ArityError()
	}
	return primitive.Error(args[0].ToString())
}

// existsFn returns whether the given path exists.
func existsFn(env *env.Environment, args []primitive.Primitive) primitive.Primitive {

	// We only need a single argument
	if len(args) != 1 {
		return primitive.ArityError()
	}

	// Which is a string
	str, ok := args[0].(primitive.String)
	if !ok {
		return primitive.Error("argument not a string")
	}

	if _, err := os.Stat(str.ToString()); err == nil {
		// path/to/whatever exists
		return primitive.Bool(true)
	}

	return primitive.Bool(false)
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

// expnFn implements "#"
func expnFn(env *env.Environment, args []primitive.Primitive) primitive.Primitive {
	if len(args) != 2 {
		return primitive.ArityError()
	}

	if _, ok := args[0].(primitive.Number); !ok {
		return primitive.Error("argument not a number")
	}
	if _, ok := args[1].(primitive.Number); !ok {
		return primitive.Error("argument not a number")
	}
	return primitive.Number(math.Pow(float64(args[0].(primitive.Number)), float64(args[1].(primitive.Number))))
}

// fileFn returns whether the given path exists, and is a file (or rather is not a directory).
func fileFn(env *env.Environment, args []primitive.Primitive) primitive.Primitive {
	// We only need a single argument
	if len(args) != 1 {
		return primitive.ArityError()
	}

	// Which is a string
	str, ok := args[0].(primitive.String)
	if !ok {
		return primitive.Error("argument not a string")
	}

	// stat the path
	info, err := os.Stat(str.ToString())

	// no error then return true, unless we've got a directory
	//
	// i.e. swallow errors
	if err == nil {
		if !info.IsDir() {
			return primitive.Bool(true)
		}
	}
	return primitive.Bool(false)
}

// fileLinesFn implements (file:lines)
func fileLinesFn(env *env.Environment, args []primitive.Primitive) primitive.Primitive {
	// We only need a single argument
	if len(args) != 1 {
		return primitive.ArityError()
	}

	// Which is a string
	fName, ok := args[0].(primitive.String)
	if !ok {
		return primitive.Error("argument not a string")
	}

	// Return value,
	var res primitive.List

	// Open the file
	file, err := os.Open(fName.ToString())
	if err != nil {
		return primitive.Error(fmt.Sprintf("failed to open %s:%s", fName.ToString(), err))
	}
	defer file.Close()

	// Read each line, and append to our list.
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		res = append(res, primitive.String(scanner.Text()))
	}

	// All done.
	return res
}

// fileReadFn implements (file:read)
func fileReadFn(env *env.Environment, args []primitive.Primitive) primitive.Primitive {
	// We only need a single argument
	if len(args) != 1 {
		return primitive.ArityError()
	}

	// Which is a string
	fName, ok := args[0].(primitive.String)
	if !ok {
		return primitive.Error("argument not a string")
	}

	data, err := os.ReadFile(fName.ToString())
	if err != nil {
		return primitive.Error(fmt.Sprintf("error reading %s %s", fName.ToString(), err))
	}
	return primitive.String(string(data))
}

// fileStatFn implements (file:stat)
//
// Return value is (NAME SIZE UID GID MODE)
func fileStatFn(env *env.Environment, args []primitive.Primitive) primitive.Primitive {
	// We only need a single argument
	if len(args) != 1 {
		return primitive.ArityError()
	}

	// Which is a string
	fName, ok := args[0].(primitive.String)
	if !ok {
		return primitive.Error("argument not a string")
	}

	// Stat the entry
	info, err := os.Stat(fName.ToString())

	if err != nil {
		return primitive.Nil{}
	}

	// If we're not on Linux the Stat_t type won't be available,
	// so we'd default to the current user.
	UID := os.Getuid()
	GID := os.Getgid()

	// But if we can get the "real" values, then use them.
	if stat, ok := info.Sys().(*syscall.Stat_t); ok {
		UID = int(stat.Uid)
		GID = int(stat.Gid)
	}

	var res primitive.List

	res = append(res, primitive.String(info.Name()))
	res = append(res, primitive.Number(info.Size()))
	res = append(res, primitive.Number(UID))
	res = append(res, primitive.Number(GID))
	res = append(res, primitive.String(info.Mode().String()))

	return res
}

// fileWriteFn implements file:write
func fileWriteFn(env *env.Environment, args []primitive.Primitive) primitive.Primitive {
	// We need two arguments
	if len(args) != 2 {
		return primitive.ArityError()
	}

	// Path is a string
	path, ok := args[0].(primitive.String)
	if !ok {
		return primitive.Error("argument not a string")
	}

	// Content is a string
	content, ok := args[1].(primitive.String)
	if !ok {
		return primitive.Error("argument not a string")
	}

	err := os.WriteFile(path.ToString(), []byte(content.ToString()), 0777)
	if err != nil {
		return primitive.Error(fmt.Sprintf("failed to write to %s:%s", path.ToString(), err))
	}
	return primitive.Nil{}
}

// gensymFn is the implementation of (gensym ..)
func gensymFn(env *env.Environment, args []primitive.Primitive) primitive.Primitive {
	// symbol characters
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

	// generate prefix
	b := make([]rune, 5)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}

	// generate with count
	symCount++
	str := fmt.Sprintf("%s%06d", string(b), symCount)
	sym := primitive.Symbol(str)
	return sym
}

// getFn is the implementation of `(get hash key)`
func getFn(env *env.Environment, args []primitive.Primitive) primitive.Primitive {

	// We need two arguments
	if len(args) != 2 {
		return primitive.ArityError()
	}

	// First is a Hash
	if _, ok := args[0].(primitive.Hash); !ok {
		return primitive.Error("argument not a hash")
	}

	tmp := args[0].(primitive.Hash)
	return tmp.Get(args[1].ToString())
}

// getenvFn is the implementation of `(getenv "PATH")`
func getenvFn(env *env.Environment, args []primitive.Primitive) primitive.Primitive {

	// If we have only a single argument
	if len(args) != 1 {
		return primitive.ArityError()
	}

	// Which is a string
	if _, ok := args[0].(primitive.String); !ok {
		return primitive.Error("argument not a string")
	}

	// Return the value
	str := args[0].(primitive.String)
	return primitive.String(os.Getenv(string(str)))
}

// globFn is the implementation of `(glob "pattern")`
func globFn(env *env.Environment, args []primitive.Primitive) primitive.Primitive {

	// If we have only a single argument
	if len(args) != 1 {
		return primitive.ArityError()
	}

	// Which is a string
	str, ok := args[0].(primitive.String)
	if !ok {
		return primitive.Error("argument not a string")
	}

	// Run the glob
	out, err := filepath.Glob(str.ToString())

	if err != nil {
		return primitive.Error(fmt.Sprintf("error running glob(%s): %s", str.ToString(), err))
	}

	var ret primitive.List

	for _, ent := range out {
		ret = append(ret, primitive.String(ent))
	}

	return ret
}

// helpFn is the implementation of `(help fn)`
func helpFn(env *env.Environment, args []primitive.Primitive) primitive.Primitive {
	// We need a single argument
	if len(args) != 1 {
		return primitive.ArityError()
	}

	// Which is a function
	proc, ok := args[0].(*primitive.Procedure)

	if !ok {
		return primitive.Error("argument not a procedure")
	}

	// Return value
	str := ""

	for _, arg := range proc.Args {
		if len(str) == 0 {
			str = "Arguments"
		}
		str += " " + arg.ToString()
	}
	if len(str) > 0 {
		str += "\n"
	}
	str += proc.Help
	return primitive.String(str)
}

// inequalityFn implements /=
func inequalityFn(env *env.Environment, args []primitive.Primitive) primitive.Primitive {

	// We need at least two arguments
	if len(args) < 2 {
		return primitive.ArityError()
	}

	// First argument must be a number.
	nA, ok := args[0].(primitive.Number)
	if !ok {
		return primitive.Error("argument was not a number")
	}

	// Now we'll loop over all other numbers
	//
	// If we got something that was already seen we can
	// terminate early but we don't because it is important
	// to also report on failures to validate types - which
	// we can't do if we bail.
	//
	ret := primitive.Bool(true)

	// Keep track of things we've seen here
	seen := make(map[float64]bool)
	seen[float64(nA)] = true

	for _, i := range args[1:] {

		// check we have a number
		nB, ok2 := i.(primitive.Number)

		if !ok2 {
			return primitive.Error("argument was not a number")
		}

		// Have we seen this?
		_, found := seen[float64(nB)]
		if found {
			ret = primitive.Bool(false)
		}
		seen[float64(nB)] = true
	}

	return ret
}

// (join (1 2 3)
func joinFn(env *env.Environment, args []primitive.Primitive) primitive.Primitive {

	// We require one argument
	if len(args) != 1 {
		return primitive.ArityError()
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

// keysFn is the implementation of `(keys hash)`
func keysFn(env *env.Environment, args []primitive.Primitive) primitive.Primitive {

	// We need a single argument
	if len(args) != 1 {
		return primitive.ArityError()
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

// listFn implements "list"
func listFn(env *env.Environment, args []primitive.Primitive) primitive.Primitive {
	return primitive.List(args)
}

// ltFn implements "<"
func ltFn(env *env.Environment, args []primitive.Primitive) primitive.Primitive {
	if len(args) != 2 {
		return primitive.ArityError()
	}

	if _, ok := args[0].(primitive.Number); !ok {
		return primitive.Error("argument not a number")
	}
	if _, ok := args[1].(primitive.Number); !ok {
		return primitive.Error("argument not a number")
	}
	return primitive.Bool(args[0].(primitive.Number) < args[1].(primitive.Number))
}

// matchFn is the implementation of (match ..)
func matchFn(env *env.Environment, args []primitive.Primitive) primitive.Primitive {

	// We need two arguments
	if len(args) != 2 {
		return primitive.ArityError()
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

// minusFn implements "-"
func minusFn(env *env.Environment, args []primitive.Primitive) primitive.Primitive {

	// ensure we have at least one argument
	if len(args) < 1 {
		return primitive.ArityError()
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

// modFn implements "%"
func modFn(env *env.Environment, args []primitive.Primitive) primitive.Primitive {
	if len(args) != 2 {
		return primitive.ArityError()
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

// msFn is the implementation of `(ms)`
func msFn(env *env.Environment, args []primitive.Primitive) primitive.Primitive {
	return primitive.Number(time.Now().UnixNano() / int64(time.Millisecond))
}

// multiplyFn implements "*"
func multiplyFn(env *env.Environment, args []primitive.Primitive) primitive.Primitive {
	// ensure we have at least one argument
	if len(args) < 1 {
		return primitive.ArityError()
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

// nilFn implements nil?
func nilFn(env *env.Environment, args []primitive.Primitive) primitive.Primitive {
	if len(args) != 1 {
		return primitive.ArityError()
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

// nowFn is the implementation of `(now)`
func nowFn(env *env.Environment, args []primitive.Primitive) primitive.Primitive {
	return primitive.Number(time.Now().Unix())
}

// ordFn is the implementation of (ord ..)
func ordFn(env *env.Environment, args []primitive.Primitive) primitive.Primitive {

	if len(args) != 1 {
		return primitive.ArityError()
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

// osFn implements (os)
func osFn(env *env.Environment, args []primitive.Primitive) primitive.Primitive {
	return primitive.String(runtime.GOOS)
}

// plusFn implements "+"
func plusFn(env *env.Environment, args []primitive.Primitive) primitive.Primitive {

	// ensure we have at least one argument
	if len(args) < 1 {
		return primitive.ArityError()
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

// printFn implements (print).
func printFn(env *env.Environment, args []primitive.Primitive) primitive.Primitive {
	// no args
	if len(args) < 1 {
		return primitive.ArityError()
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

// randomFn implements (random).
func randomFn(env *env.Environment, args []primitive.Primitive) primitive.Primitive {
	if len(args) != 1 {
		return primitive.ArityError()
	}

	// ensure we received a number
	num, ok := args[0].(primitive.Number)

	if !ok {
		return primitive.Error("argument not a number")
	}

	return primitive.Number(rand.Intn(int(num)))

}

// setFn is the implementation of `(set hash key val)`
func setFn(env *env.Environment, args []primitive.Primitive) primitive.Primitive {

	// We need three arguments
	if len(args) != 3 {
		return primitive.ArityError()
	}

	// First is a Hash
	if _, ok := args[0].(primitive.Hash); !ok {
		return primitive.Error("argument not a hash")
	}

	tmp := args[0].(primitive.Hash)
	tmp.Set(args[1].ToString(), args[2])
	return args[2]
}

// shellFn runs a command via the shell
func shellFn(env *env.Environment, args []primitive.Primitive) primitive.Primitive {

	// We need one argument
	if len(args) != 1 {
		return primitive.ArityError()
	}

	// The argument must be a list
	lst, ok := args[0].(primitive.List)
	if !ok {
		return primitive.Error("argument not a list")
	}

	// Command to  run, and arguments
	cArgs := []string{}

	for _, arg := range lst {
		cArgs = append(cArgs, arg.ToString())
	}

	// If we're running a test-case we'll stop here, because
	// fuzzing might run commands.
	if os.Getenv("FUZZ") != "" {
		return primitive.List{}
	}

	cmd := exec.Command(cArgs[0], cArgs[1:]...)
	var outb, errb bytes.Buffer
	cmd.Stdout = &outb
	cmd.Stderr = &errb
	err := cmd.Run()
	if err != nil {
		return primitive.Error(fmt.Sprintf("error running command %s:%s", lst, err))
	}

	var ret primitive.List
	ret = append(ret, primitive.String(outb.String()))
	ret = append(ret, primitive.String(errb.String()))

	return ret
}

// sortFn implements (sort)
func sortFn(env *env.Environment, args []primitive.Primitive) primitive.Primitive {
	// If we have only a single argument
	if len(args) != 1 {
		return primitive.ArityError()
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

// (split "str" "by")
func splitFn(env *env.Environment, args []primitive.Primitive) primitive.Primitive {

	// We require two arguments
	if len(args) != 2 {
		return primitive.ArityError()
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

// (sprintf "fmt" "arg1" ... "argN")
func sprintfFn(env *env.Environment, args []primitive.Primitive) primitive.Primitive {

	// we need 2+ arguments
	if len(args) < 2 {
		return primitive.ArityError()
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

// strFn implements "str"
func strFn(env *env.Environment, args []primitive.Primitive) primitive.Primitive {
	if len(args) != 1 {
		return primitive.ArityError()
	}
	return primitive.String(args[0].ToString())
}

// timeFn returns the current (HH, MM, SS) as a list.
func timeFn(env *env.Environment, args []primitive.Primitive) primitive.Primitive {
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

// typeFn implements "type"
func typeFn(env *env.Environment, args []primitive.Primitive) primitive.Primitive {
	if len(args) != 1 {
		return primitive.ArityError()
	}
	return primitive.String(args[0].Type())
}

// valsFn is the implementation of `(vals hash)`
func valsFn(env *env.Environment, args []primitive.Primitive) primitive.Primitive {

	// We need a single argument
	if len(args) != 1 {
		return primitive.ArityError()
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
