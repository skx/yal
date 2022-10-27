package builtins

import (
	"math"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/skx/yal/env"
	"github.com/skx/yal/eval"
	"github.com/skx/yal/primitive"
	"github.com/skx/yal/stdlib"
)

// ENV contains a local environment for the test functions
var ENV *env.Environment

// init ensures our environment pointer is up to date.
func init() {
	ENV = env.New()
}

func TestArch(t *testing.T) {

	// No arguments
	out := archFn(ENV, []primitive.Primitive{})

	// Will lead to a number
	e, ok := out.(primitive.String)
	if !ok {
		t.Fatalf("expected string, got %v", out)
	}

	if e.ToString() != runtime.GOARCH {
		t.Fatalf("got wrong value for runtime architecture")
	}
}

// Test (car
func TestCar(t *testing.T) {

	// No arguments
	out := carFn(ENV, []primitive.Primitive{})

	// Will lead to an error
	e, ok := out.(primitive.Error)
	if !ok {
		t.Fatalf("expected error, got %v", out)
	}
	if e != primitive.ArityError() {
		t.Fatalf("got error, but wrong one %v", out)
	}

	// One argument
	out = carFn(ENV, []primitive.Primitive{
		primitive.Number(3),
	})

	// Will lead to an error
	e, ok = out.(primitive.Error)
	if !ok {
		t.Fatalf("expected error, got %v", out)
	}
	if !strings.Contains(string(e), "not a list") {
		t.Fatalf("got error, but wrong one %v", out)
	}

	// Now a list
	out = carFn(ENV, []primitive.Primitive{
		primitive.List{
			primitive.Number(3),
			primitive.Number(4),
		},
	})

	// No error
	r, ok2 := out.(primitive.Number)
	if !ok2 {
		t.Fatalf("expected number, got %v", out)
	}
	if r.ToString() != "3" {
		t.Fatalf("got wrong result : %v", r)
	}

	// Now a list which is empty
	out = carFn(ENV, []primitive.Primitive{
		primitive.List{},
	})

	// No error
	_, ok3 := out.(primitive.Nil)
	if !ok3 {
		t.Fatalf("expected nil, got %v", out)
	}

}

// Test (cdr
func TestCdr(t *testing.T) {

	// No arguments
	out := cdrFn(ENV, []primitive.Primitive{})

	// Will lead to an error
	e, ok := out.(primitive.Error)
	if !ok {
		t.Fatalf("expected error, got %v", out)
	}
	if e != primitive.ArityError() {
		t.Fatalf("got error, but wrong one %v", out)
	}

	// One argument
	out = cdrFn(ENV, []primitive.Primitive{
		primitive.Number(3),
	})

	// Will lead to an error
	e, ok = out.(primitive.Error)
	if !ok {
		t.Fatalf("expected error, got %v", out)
	}
	if !strings.Contains(string(e), "not a list") {
		t.Fatalf("got error, but wrong one %v", out)
	}

	// Now a list
	out = cdrFn(ENV, []primitive.Primitive{
		primitive.List{
			primitive.Number(3),
			primitive.Number(4),
			primitive.Number(5),
		},
	})

	// No error
	r, ok2 := out.(primitive.List)
	if !ok2 {
		t.Fatalf("expected list, got %v", out)
	}
	if r.ToString() != "(4 5)" {
		t.Fatalf("got wrong result : %v", r)
	}

	// Now a list which is empty
	out = cdrFn(ENV, []primitive.Primitive{
		primitive.List{},
	})

	// No error
	_, ok3 := out.(primitive.Nil)
	if !ok3 {
		t.Fatalf("expected nil, got %v", out)
	}
}

// TestCharEquals tests "chare=" (character equality)
func TestCharEquals(t *testing.T) {

	// No arguments
	out := charEqualsFn(ENV, []primitive.Primitive{})

	// Will lead to an error
	e, ok := out.(primitive.Error)
	if !ok {
		t.Fatalf("expected error, got %v", out)
	}
	if e != primitive.ArityError() {
		t.Fatalf("got error, but wrong one %v", out)
	}

	// One bogus argument
	out = charEqualsFn(ENV, []primitive.Primitive{
		primitive.Number(33),
		primitive.Character("a"),
	})

	// Will lead to an error
	e, ok = out.(primitive.Error)
	if !ok {
		t.Fatalf("expected error, got %v", out)
	}
	if !strings.Contains(string(e), "was not a character") {
		t.Fatalf("got error, but wrong one %v", out)
	}

	//
	// Now a real one: equal
	//
	out = charEqualsFn(ENV, []primitive.Primitive{
		primitive.Character("a"),
		primitive.Character("a"),
	})

	// Will work
	n, ok2 := out.(primitive.Bool)
	if !ok2 {
		t.Fatalf("expected bool, got %v", out)
	}
	if n != true {
		t.Fatalf("got wrong result")
	}

	//
	// Now a real one: equal - but multiple values
	//
	out = charEqualsFn(ENV, []primitive.Primitive{
		primitive.Character("b"),
		primitive.Character("b"),
		primitive.Character("b"),
		primitive.Character("b"),
		primitive.Character("b"),
		primitive.Character("b"),
	})

	// Will work
	n, ok2 = out.(primitive.Bool)
	if !ok2 {
		t.Fatalf("expected bool, got %v", out)
	}
	if n != true {
		t.Fatalf("got wrong result")
	}

	//
	// Now a real one: unequal values
	//
	out = charEqualsFn(ENV, []primitive.Primitive{
		primitive.Character("a"),
		primitive.Character("b"),
		primitive.Character("c"),
		primitive.Character("d"),
	})

	// Will work
	n, ok2 = out.(primitive.Bool)
	if !ok2 {
		t.Fatalf("expected bool, got %v", out)
	}
	if n != false {
		t.Fatalf("got wrong result")
	}

	//
	// Now with wrong types - last one is wrong
	//
	out = charEqualsFn(ENV, []primitive.Primitive{
		primitive.Character("a"),
		primitive.Character("b"),
		primitive.Character("c"),
		primitive.Character("d"),
		primitive.String("e"),
	})

	e, ok = out.(primitive.Error)
	if !ok {
		t.Fatalf("expected error, got %v", out)
	}
	if !strings.Contains(string(e), "was not a character") {
		t.Fatalf("got error, but wrong one '%v'", e)
	}

	//
	// Now with wrong types
	//
	out = charEqualsFn(ENV, []primitive.Primitive{
		primitive.Character("b"),
		primitive.Number(9),
	})

	// Will report an error
	e, ok = out.(primitive.Error)
	if !ok {
		t.Fatalf("expected error, got %v", out)
	}
	if !strings.Contains(string(e), "was not a character") {
		t.Fatalf("got error, but wrong one %v", out)
	}
}

// test char<
func TestCharLt(t *testing.T) {

	// No arguments
	out := charLtFn(ENV, []primitive.Primitive{})

	// Will lead to an error
	e, ok := out.(primitive.Error)
	if !ok {
		t.Fatalf("expected error, got %v", out)
	}
	if e != primitive.ArityError() {
		t.Fatalf("got error, but wrong one %v", out)
	}

	// Argument which isn't a character
	out = charLtFn(ENV, []primitive.Primitive{
		primitive.String("foo"),
		primitive.String("foo"),
	})

	// Will lead to an error
	e, ok = out.(primitive.Error)
	if !ok {
		t.Fatalf("expected error, got %v", out)
	}
	if !strings.Contains(string(e), "not a character") {
		t.Fatalf("got error, but wrong one %v", out)
	}

	// Argument which isn't a character
	out = charLtFn(ENV, []primitive.Primitive{
		primitive.Character("a"),
		primitive.String("foo"),
	})

	// Will lead to an error
	e, ok = out.(primitive.Error)
	if !ok {
		t.Fatalf("expected error, got %v", out)
	}
	if !strings.Contains(string(e), "not a character") {
		t.Fatalf("got error, but wrong one %v", out)
	}

	//
	// Now a real one
	//
	out = charLtFn(ENV, []primitive.Primitive{
		primitive.Character("a"),
		primitive.Character("b"),
	})

	// Will work
	n, ok2 := out.(primitive.Bool)
	if !ok2 {
		t.Fatalf("expected bool, got %v", out)
	}
	if n != true {
		t.Fatalf("got wrong result")
	}

}
func TestChr(t *testing.T) {

	// no arguments
	out := chrFn(ENV, []primitive.Primitive{})

	// Will lead to an error
	e, ok := out.(primitive.Error)
	if !ok {
		t.Fatalf("expected error, got %v", out)
	}
	if e != primitive.ArityError() {
		t.Fatalf("got error, but wrong one:%s", out)
	}

	// First argument must be a number
	out = chrFn(ENV, []primitive.Primitive{
		primitive.String("4"),
	})

	// Will lead to an error
	e, ok = out.(primitive.Error)
	if !ok {
		t.Fatalf("expected error, got %v", out)
	}
	if !strings.Contains(string(e), "not a number") {
		t.Fatalf("got error, but wrong one %v", out)
	}

	// Now a valid call 42 => "*"
	val := chrFn(ENV, []primitive.Primitive{
		primitive.Number(42),
	})

	r, ok2 := val.(primitive.Character)
	if !ok2 {
		t.Fatalf("expected character, got %v", val)
	}
	if r.ToString() != "*" {
		t.Fatalf("got wrong result %v", r)
	}
}

func TestCons(t *testing.T) {

	// No arguments
	out := consFn(ENV, []primitive.Primitive{})

	// Will lead to an error
	e, ok := out.(primitive.Error)
	if !ok {
		t.Fatalf("expected error, got %v", out)
	}
	if e != primitive.ArityError() {
		t.Fatalf("got error, but wrong one %v", out)
	}

	// one argument, string -> list
	out = consFn(ENV, []primitive.Primitive{
		primitive.String("steve"),
	})

	out, ok2 := out.(primitive.List)
	if !ok2 {
		t.Errorf("expected list")
	}
	if out.ToString() != "(steve)" {
		t.Fatalf("wrong result")
	}

	// A list with a nil second element is gonna be truncated
	out = consFn(ENV, []primitive.Primitive{
		primitive.String("steve"),
		primitive.Nil{},
	})

	out, ok2 = out.(primitive.List)
	if !ok2 {
		t.Errorf("expected list")
	}
	if out.ToString() != "(steve)" {
		t.Fatalf("wrong result")
	}

	// A list and a number
	a := []primitive.Primitive{
		primitive.List{
			primitive.Number(3),
			primitive.Number(4),
		},
		primitive.Number(5),
	}

	// A number and a list
	b := []primitive.Primitive{
		primitive.Number(5),
		primitive.List{
			primitive.Number(3),
			primitive.Number(4),
		},
	}

	// first one
	out = consFn(ENV, a)
	out, ok2 = out.(primitive.List)
	if !ok2 {
		t.Errorf("expected list")
	}
	if out.ToString() != "((3 4) 5)" {
		t.Fatalf("wrong result, got %v", out)
	}

	// second one
	out = consFn(ENV, b)
	out, ok2 = out.(primitive.List)
	if !ok2 {
		t.Errorf("expected list")
	}
	if out.ToString() != "(5 3 4)" {
		t.Fatalf("wrong result, got %v", out)
	}
}

func TestContains(t *testing.T) {

	// no arguments
	out := containsFn(ENV, []primitive.Primitive{})

	// Will lead to an error
	e, ok := out.(primitive.Error)
	if !ok {
		t.Fatalf("expected error, got %v", out)
	}
	if !strings.Contains(string(e), "argument") {
		t.Fatalf("got error, but wrong one")
	}

	// First argument must be a hash
	out = containsFn(ENV, []primitive.Primitive{
		primitive.String("foo"),
		primitive.String("bar"),
	})

	// Will lead to an error
	e, ok = out.(primitive.Error)
	if !ok {
		t.Fatalf("expected error, got %v", out)
	}
	if !strings.Contains(string(e), "not a hash") {
		t.Fatalf("got error, but wrong one %v", out)
	}

	// create a hash
	h := primitive.NewHash()
	h.Set("XXX", primitive.String("Last"))
	h.Set("Name", primitive.String("Steve"))
	h.Set("Age", primitive.Number(43))
	h.Set("Location", primitive.String("Helsinki"))

	// Should have Age
	res := containsFn(ENV, []primitive.Primitive{
		h,
		primitive.String("Age"),
	})

	// Will lead to a bool
	v, ok2 := res.(primitive.Bool)
	if !ok2 {
		t.Fatalf("expected bool, got %v", res)
	}
	if v != primitive.Bool(true) {
		t.Fatalf("failed to find expected key")
	}

	// Should have Age - as a symbol
	res = containsFn(ENV, []primitive.Primitive{
		h,
		primitive.Symbol("Age"),
	})

	// Will lead to a bool
	v, ok2 = res.(primitive.Bool)
	if !ok2 {
		t.Fatalf("expected bool, got %v", res)
	}
	if v != primitive.Bool(true) {
		t.Fatalf("failed to find expected key")
	}

	// Should NOT have Cake
	res = containsFn(ENV, []primitive.Primitive{
		h,
		primitive.String("Cake"),
	})

	// Will lead to a bool
	v, ok2 = res.(primitive.Bool)
	if !ok2 {
		t.Fatalf("expected bool, got %v", res)
	}
	if v != primitive.Bool(false) {
		t.Fatalf("unexpectedly found missing key")
	}
}

// We don't really test the contents here.
func TestDateTime(t *testing.T) {

	// No arguments
	dt := dateFn(ENV, []primitive.Primitive{})
	tm := timeFn(ENV, []primitive.Primitive{})

	// date should return a list
	out, ok := dt.(primitive.List)
	if !ok {
		t.Fatalf("expected list for (date), got %v", dt)
	}

	// "weekday", "day", "month", "year" == four entries
	if len(out) != 4 {
		t.Fatalf("date list had the wrong length, got %d: %v", len(out), out)
	}

	// time should return a list
	out, ok = tm.(primitive.List)
	if !ok {
		t.Fatalf("expected list for (time), got %v", tm)
	}

	// "hour", "minute", "seconds" == three entries
	if len(out) != 3 {
		t.Fatalf("time list had the wrong length, got %d: %v", len(out), out)
	}

}

// TestDirectory tests directory?
func TestDirectory(t *testing.T) {

	// No arguments
	out := directoryFn(ENV, []primitive.Primitive{})

	// Will lead to an error
	e, ok := out.(primitive.Error)
	if !ok {
		t.Fatalf("expected error, got %v", out)
	}
	if e != primitive.ArityError() {
		t.Fatalf("got error, but wrong one %v", out)
	}

	// One argument, wrong type
	out = directoryFn(ENV, []primitive.Primitive{primitive.Number(33)})

	// Will lead to an error
	e, ok = out.(primitive.Error)
	if !ok {
		t.Fatalf("expected error, got %v", out)
	}
	if !strings.Contains(string(e), "not a string") {
		t.Fatalf("got error, but wrong one %v", out)
	}

	// Create a temporary directory
	path, err := os.MkdirTemp("", "directory_")
	if err != nil {
		t.Fatalf("failed to create a temporary direcotry")
	}

	// directory? should return true
	res := directoryFn(ENV, []primitive.Primitive{primitive.String(path)})

	// So we need a boolean result
	r, ok2 := res.(primitive.Bool)
	if !ok2 {
		t.Fatalf("found wrong type in result")
	}

	if r.ToString() != "#t" {
		t.Fatalf("wrong result, got %v", r.ToString())
	}

	// Remove the file
	os.RemoveAll(path)

	// directory? should return false
	res = directoryFn(ENV, []primitive.Primitive{primitive.String(path)})

	// So we need a booelan result
	r, ok2 = res.(primitive.Bool)
	if !ok2 {
		t.Fatalf("found wrong type in result, got %v", res)
	}

	if r.ToString() != "#f" {
		t.Fatalf("wrong result, got %v", r.ToString())
	}
}

// TestDirectoryEntries tests directory:entries
func TestDirectoryEntries(t *testing.T) {

	// No arguments
	out := directoryEntriesFn(ENV, []primitive.Primitive{})

	// Will lead to an error
	e, ok := out.(primitive.Error)
	if !ok {
		t.Fatalf("expected error, got %v", out)
	}
	if !strings.Contains(string(e), "argument count") {
		t.Fatalf("got error, but wrong one %v", out)
	}

	// One argument, wrong type
	out = directoryEntriesFn(ENV, []primitive.Primitive{primitive.Number(33)})

	// Will lead to an error
	e, ok = out.(primitive.Error)
	if !ok {
		t.Fatalf("expected error, got %v", out)
	}
	if !strings.Contains(string(e), "not a string") {
		t.Fatalf("got error, but wrong one %v", out)
	}

	// Create a temporary directory
	path, err := os.MkdirTemp("", "directory_")
	if err != nil {
		t.Fatalf("failed to create a temporary directory")
	}

	// Populate two files.
	a := filepath.Join(path, "one.txt")
	b := filepath.Join(path, "two.foo")

	err = os.WriteFile(a, []byte("one.txt"), 0777)
	if err != nil {
		t.Fatalf("failed to write to file")
	}
	err = os.WriteFile(b, []byte("two.foo"), 0777)
	if err != nil {
		t.Fatalf("failed to write to file")
	}

	// Now we should find a list of two files if we walk the
	// temporary directory.
	res := directoryEntriesFn(ENV, []primitive.Primitive{
		primitive.String(path),
	})

	// Will lead to a list
	lst, ok2 := res.(primitive.List)
	if !ok2 {
		t.Fatalf("expected list, got %v", out)
	}
	if len(lst) != 3 {
		t.Fatalf("failed to find expected file-count, got %v", lst)
	}

	// Delete one of the two files, and ensure we still find results
	os.Remove(a)

	// walk again
	res = directoryEntriesFn(ENV, []primitive.Primitive{
		primitive.String(path),
	})

	// Will lead to a list
	lst, ok2 = res.(primitive.List)
	if !ok2 {
		t.Fatalf("expected list, got %v", out)
	}
	if len(lst) != 2 {
		t.Fatalf("failed to find expected file-count, got %v", lst)
	}

	// Finally cleanup
	os.RemoveAll(path)

	// walk a missing directory
	res = directoryEntriesFn(ENV, []primitive.Primitive{
		primitive.String("/ fdsf dsf /this doesnt exist"),
	})

	// Will lead to a list still
	lst, ok2 = res.(primitive.List)
	if !ok2 {
		t.Fatalf("expected list, got %v", out)
	}

	// but empty
	if len(lst) != 0 {
		t.Fatalf("failed to find expected file-count, got %v", lst)
	}

}

// TestDivide tests "*"
func TestDivide(t *testing.T) {

	// Test for "equality"
	//
	// Because floating points are hard
	almostEqual := func(a, b float64) bool {
		// Arbitrary equality threshold
		float64EqualityThreshold := float64(1.0 / 1000000)

		return math.Abs(a-b) <= float64EqualityThreshold
	}

	// No arguments
	out := divideFn(ENV, []primitive.Primitive{})

	// Will lead to an error
	e, ok := out.(primitive.Error)
	if !ok {
		t.Fatalf("expected error, got %v", out)
	}
	if e != primitive.ArityError() {
		t.Fatalf("got error, but wrong one %v", out)
	}

	// Argument which isn't a number
	out = divideFn(ENV, []primitive.Primitive{
		primitive.String("foo"),
	})

	// Will lead to an error
	e, ok = out.(primitive.Error)
	if !ok {
		t.Fatalf("expected error, got %v", out)
	}
	if !strings.Contains(string(e), "not a number") {
		t.Fatalf("got error, but wrong one %v", out)
	}

	// Argument which isn't a number
	out = divideFn(ENV, []primitive.Primitive{
		primitive.Number(32),
		primitive.String("foo"),
	})

	// Will lead to an error
	e, ok = out.(primitive.Error)
	if !ok {
		t.Fatalf("expected error, got %v", out)
	}
	if !strings.Contains(string(e), "not a number") {
		t.Fatalf("got error, but wrong one %v", out)
	}

	// Division by zero
	out = divideFn(ENV, []primitive.Primitive{
		primitive.Number(32),
		primitive.Number(0),
	})

	// Will lead to an error
	e, ok = out.(primitive.Error)
	if !ok {
		t.Fatalf("expected error, got %v", out)
	}
	if !strings.Contains(string(e), "division by zero") {
		t.Fatalf("got error, but wrong one %v", out)
	}

	//
	// Now a real one
	//
	out = divideFn(ENV, []primitive.Primitive{
		primitive.Number(12),
		primitive.Number(3),
	})

	// Will work
	n, ok2 := out.(primitive.Number)
	if !ok2 {
		t.Fatalf("expected number, got %v", out)
	}
	if n != 4 {
		t.Fatalf("got wrong result")
	}

	// Test only a single argument.
	out = divideFn(ENV, []primitive.Primitive{
		primitive.Number(12),
	})

	// Will work
	n, ok2 = out.(primitive.Number)
	if !ok2 {
		t.Fatalf("expected number, got %v", out)
	}
	if !almostEqual(float64(n), 0.083333) {
		t.Fatalf("got wrong result: %f", n)
	}

	// Another test
	out = divideFn(ENV, []primitive.Primitive{
		primitive.Number(-3),
	})

	// Will work
	n, ok2 = out.(primitive.Number)
	if !ok2 {
		t.Fatalf("expected number, got %v", out)
	}
	if !almostEqual(float64(n), -(1 / 3.0)) {
		t.Fatalf("got wrong result: %f", n)
	}

}

// TestEnsureHelpPresent ensures that all our built-in functions have
// help-text available
func TestEnsureHelpPresent(t *testing.T) {

	// create a new environment, and populate it
	e := env.New()
	PopulateEnvironment(e)

	// For each function
	items := e.Items()

	for name, val := range items {

		proc, ok := val.(*primitive.Procedure)
		if ok {

			t.Run("Testing "+name, func(t *testing.T) {

				if len(proc.Help) == 0 {
					t.Fatalf("help text is unset")
				}
			})
		}
	}
}

// TestEq tests "eq" (non-numerical equality)
func TestEq(t *testing.T) {

	// No arguments
	out := eqFn(ENV, []primitive.Primitive{})

	// Will lead to an error
	e, ok := out.(primitive.Error)
	if !ok {
		t.Fatalf("expected error, got %v", out)
	}
	if e != primitive.ArityError() {
		t.Fatalf("got error, but wrong one %v", out)
	}

	//
	// Now a real one: equal
	//
	out = eqFn(ENV, []primitive.Primitive{
		primitive.Number(9),
		primitive.Number(9),
	})

	// Will work
	n, ok2 := out.(primitive.Bool)
	if !ok2 {
		t.Fatalf("expected bool, got %v", out)
	}
	if n != true {
		t.Fatalf("got wrong result")
	}

	//
	// Now a real one: unequal values
	//
	out = eqFn(ENV, []primitive.Primitive{
		primitive.String("99"),
		primitive.String("9"),
	})

	// Will work
	n, ok2 = out.(primitive.Bool)
	if !ok2 {
		t.Fatalf("expected bool, got %v", out)
	}
	if n != false {
		t.Fatalf("got wrong result")
	}

	//
	// Now a real one: unequal types
	//
	out = eqFn(ENV, []primitive.Primitive{
		primitive.Number(9),
		primitive.String("9"),
	})

	// Will work
	n, ok2 = out.(primitive.Bool)
	if !ok2 {
		t.Fatalf("expected bool, got %v", out)
	}
	if n != false {
		t.Fatalf("got wrong result")
	}
}

// TestEquals tests "=" (numerical equality)
func TestEquals(t *testing.T) {

	// No arguments
	out := equalsFn(ENV, []primitive.Primitive{})

	// Will lead to an error
	e, ok := out.(primitive.Error)
	if !ok {
		t.Fatalf("expected error, got %v", out)
	}
	if e != primitive.ArityError() {
		t.Fatalf("got error, but wrong one %v", out)
	}

	//
	// Now a real one: equal
	//
	out = equalsFn(ENV, []primitive.Primitive{
		primitive.Number(9),
		primitive.Number(9),
	})

	// Will work
	n, ok2 := out.(primitive.Bool)
	if !ok2 {
		t.Fatalf("expected bool, got %v", out)
	}
	if n != true {
		t.Fatalf("got wrong result")
	}

	//
	// Now a real one: equal - but multiple values
	//
	out = equalsFn(ENV, []primitive.Primitive{
		primitive.Number(9),
		primitive.Number(9),
		primitive.Number(9),
		primitive.Number(9),
		primitive.Number(9),
		primitive.Number(9),
	})

	// Will work
	n, ok2 = out.(primitive.Bool)
	if !ok2 {
		t.Fatalf("expected bool, got %v", out)
	}
	if n != true {
		t.Fatalf("got wrong result")
	}

	//
	// Now a real one: unequal values
	//
	out = equalsFn(ENV, []primitive.Primitive{
		primitive.Number(99),
		primitive.Number(9),
		primitive.Number(99),
		primitive.Number(9),
	})

	// Will work
	n, ok2 = out.(primitive.Bool)
	if !ok2 {
		t.Fatalf("expected bool, got %v", out)
	}
	if n != false {
		t.Fatalf("got wrong result")
	}

	//
	// Now with wrong types - last one is wrong
	//
	out = equalsFn(ENV, []primitive.Primitive{
		primitive.Number(1),
		primitive.Number(2),
		primitive.Number(3),
		primitive.Number(4),
		primitive.String("5"),
	})

	e, ok = out.(primitive.Error)
	if !ok {
		t.Fatalf("expected error, got %v", out)
	}
	if !strings.Contains(string(e), "was not a number") {
		t.Fatalf("got error, but wrong one '%v'", e)
	}

	//
	// Now with wrong types
	//
	out = equalsFn(ENV, []primitive.Primitive{
		primitive.String("9"),
		primitive.Number(9),
	})

	// Will report an error
	e, ok = out.(primitive.Error)
	if !ok {
		t.Fatalf("expected error, got %v", out)
	}
	if !strings.Contains(string(e), "was not a number") {
		t.Fatalf("got error, but wrong one %v", out)
	}
}

// TestError tests error.
func TestError(t *testing.T) {

	// No arguments
	out := errorFn(ENV, []primitive.Primitive{})

	// Will lead to an error
	e, ok := out.(primitive.Error)
	if !ok {
		t.Fatalf("expected error, got %v", out)
	}
	if e != primitive.ArityError() {
		t.Fatalf("got error, but wrong one %v", out)
	}

	// calling with an arg
	out = errorFn(ENV, []primitive.Primitive{
		primitive.String("No Cheese Detected"),
	})

	// Will lead to an string
	e, ok = out.(primitive.Error)
	if !ok {
		t.Fatalf("expected error, got %v", out)
	}
	if e != "No Cheese Detected" {
		t.Fatalf("got wrong error %v", out)
	}
}

// TestExists tests exists?
func TestExists(t *testing.T) {

	// No arguments
	out := existsFn(ENV, []primitive.Primitive{})

	// Will lead to an error
	e, ok := out.(primitive.Error)
	if !ok {
		t.Fatalf("expected error, got %v", out)
	}
	if !strings.Contains(string(e), "argument count") {
		t.Fatalf("got error, but wrong one %v", out)
	}

	// One argument, wrong type
	out = existsFn(ENV, []primitive.Primitive{primitive.Number(33)})

	// Will lead to an error
	e, ok = out.(primitive.Error)
	if !ok {
		t.Fatalf("expected error, got %v", out)
	}
	if !strings.Contains(string(e), "not a string") {
		t.Fatalf("got error, but wrong one %v", out)
	}

	// Create a temporary file
	tmpfile, err := os.CreateTemp("", "exists")
	if err != nil {
		t.Fatalf("failed to create a temporary file")
	}

	// Exists should return true
	res := existsFn(ENV, []primitive.Primitive{primitive.String(tmpfile.Name())})

	// So we need a booelan result
	r, ok2 := res.(primitive.Bool)
	if !ok2 {
		t.Fatalf("found wrong type in result")
	}

	if r.ToString() != "#t" {
		t.Fatalf("wrong result, got %v", r.ToString())
	}

	// Remove the file
	os.Remove(tmpfile.Name())

	// Exists should return false
	res = existsFn(ENV, []primitive.Primitive{primitive.String(tmpfile.Name())})

	// So we need a booelan result
	r, ok2 = res.(primitive.Bool)
	if !ok2 {
		t.Fatalf("found wrong type in result")
	}

	if r.ToString() != "#f" {
		t.Fatalf("wrong result, got %v", r.ToString())
	}
}

// TestExpandString tests the escape-expansion of the various characters
func TestExpandString(t *testing.T) {

	type TC struct {
		in  string
		out string
	}

	tests := []TC{

		{in: "steve", out: "steve"},
		{in: "steve\\tkemp", out: "steve\tkemp"},
		{in: "steve\\rkemp", out: "steve\rkemp"},
		{in: "steve\\nkemp", out: "steve\nkemp"},
		{in: "steve\"kemp", out: "steve\"kemp"},
		{in: "steve\\\\kemp", out: "steve\\kemp"},
		{in: "steve\\bkemp", out: "steve\\bkemp"},
		{in: "steve\\ekemp", out: "steve" + string(rune(033)) + "kemp"},
	}

	for i, test := range tests {

		if expandStr(test.in) != test.out {
			t.Fatalf("%d: expected %s, got %s", i, test.out, expandStr(test.in))
		}
	}
}

// TestExplode tests explode?
func TestExplode(t *testing.T) {

	// No arguments
	out := explodeFn(ENV, []primitive.Primitive{})

	// Will lead to an error
	e, ok := out.(primitive.Error)
	if !ok {
		t.Fatalf("expected error, got %v", out)
	}
	if e != primitive.ArityError() {
		t.Fatalf("got error, but wrong one %v", out)
	}

	// An argument that isn't a string
	out = explodeFn(ENV, []primitive.Primitive{
		primitive.Number(3),
	})

	// Will lead to an error
	e, ok = out.(primitive.Error)
	if !ok {
		t.Fatalf("expected error, got %v", out)
	}
	if !strings.Contains(string(e), "not a string") {
		t.Fatalf("got error, but wrong one %v", out)
	}

	//
	// Now a proper string
	//
	out = explodeFn(ENV, []primitive.Primitive{
		primitive.String("fooπs"),
	})

	// Will lead to a list
	l, ok2 := out.(primitive.List)
	if !ok2 {
		t.Fatalf("expected list, got %v", out)
	}
	if len(l) != 5 {
		t.Fatalf("split list had the wrong length:%v", l)
	}
	if l.ToString() != "(f o o π s)" {
		t.Fatalf("got wrong result %v", out)
	}
}

// TestExpn tests "#"
func TestExpn(t *testing.T) {

	// No arguments
	out := expnFn(ENV, []primitive.Primitive{})

	// Will lead to an error
	e, ok := out.(primitive.Error)
	if !ok {
		t.Fatalf("expected error, got %v", out)
	}
	if e != primitive.ArityError() {
		t.Fatalf("got error, but wrong one %v", out)
	}

	// Argument which isn't a number
	out = expnFn(ENV, []primitive.Primitive{
		primitive.String("foo"),
		primitive.String("foo"),
	})

	// Will lead to an error
	e, ok = out.(primitive.Error)
	if !ok {
		t.Fatalf("expected error, got %v", out)
	}
	if !strings.Contains(string(e), "not a number") {
		t.Fatalf("got error, but wrong one %v", out)
	}

	// Argument which isn't a number
	out = expnFn(ENV, []primitive.Primitive{
		primitive.Number(32),
		primitive.String("foo"),
	})

	// Will lead to an error
	e, ok = out.(primitive.Error)
	if !ok {
		t.Fatalf("expected error, got %v", out)
	}
	if !strings.Contains(string(e), "not a number") {
		t.Fatalf("got error, but wrong one %v", out)
	}

	//
	// Now a real one
	//
	out = expnFn(ENV, []primitive.Primitive{
		primitive.Number(9),
		primitive.Number(0.5),
	})

	// Will work
	n, ok2 := out.(primitive.Number)
	if !ok2 {
		t.Fatalf("expected number, got %v", out)
	}
	if n != 3 {
		t.Fatalf("got wrong result")
	}
}

// TestFile tests file?
func TestFile(t *testing.T) {

	// No arguments
	out := fileFn(ENV, []primitive.Primitive{})

	// Will lead to an error
	e, ok := out.(primitive.Error)
	if !ok {
		t.Fatalf("expected error, got %v", out)
	}
	if !strings.Contains(string(e), "argument count") {
		t.Fatalf("got error, but wrong one %v", out)
	}

	// One argument, wrong type
	out = fileFn(ENV, []primitive.Primitive{primitive.Number(33)})

	// Will lead to an error
	e, ok = out.(primitive.Error)
	if !ok {
		t.Fatalf("expected error, got %v", out)
	}
	if !strings.Contains(string(e), "not a string") {
		t.Fatalf("got error, but wrong one %v", out)
	}

	// Create a temporary directory
	path, err := os.MkdirTemp("", "directory_")
	if err != nil {
		t.Fatalf("failed to create a temporary directory")
	}

	// file? should return false
	res := fileFn(ENV, []primitive.Primitive{primitive.String(path)})

	// So we need a boolean result
	r, ok2 := res.(primitive.Bool)
	if !ok2 {
		t.Fatalf("found wrong type in result")
	}

	if r.ToString() != "#f" {
		t.Fatalf("wrong result, got %v", r.ToString())
	}

	// Remove the file
	os.RemoveAll(path)

	// Now create a file.
	tmp, _ := os.CreateTemp("", "yal")
	err = os.WriteFile(tmp.Name(), []byte("I like cake"), 0777)
	if err != nil {
		t.Fatalf("failed to write to file")
	}
	defer os.Remove(tmp.Name())

	// file? should return true
	res = fileFn(ENV, []primitive.Primitive{primitive.String(tmp.Name())})

	// So we need a boolean result
	r, ok2 = res.(primitive.Bool)
	if !ok2 {
		t.Fatalf("found wrong type in result, got %v", res)
	}

	if r.ToString() != "#t" {
		t.Fatalf("wrong result, got %v", r.ToString())
	}

}

// TestFileLines tests file:lines
func TestFileLines(t *testing.T) {

	// calling with no argument
	out := fileLinesFn(ENV, []primitive.Primitive{})

	// Will lead to an error
	_, ok := out.(primitive.Error)
	if !ok {
		t.Fatalf("expected error, got %v", out)
	}

	// calling with the wrong argument type.
	out = fileLinesFn(ENV, []primitive.Primitive{
		primitive.Number(3)})

	// Will lead to an error
	_, ok = out.(primitive.Error)
	if !ok {
		t.Fatalf("expected error, got %v", out)
	}

	// Call with a file that doesn't exist
	out = fileLinesFn(ENV, []primitive.Primitive{
		primitive.String("path/not/found")})

	_, ok = out.(primitive.Error)
	if !ok {
		t.Fatalf("expected error, got %v", out)
	}

	// Create a temporary file, and read the contents
	tmp, _ := os.CreateTemp("", "yal")
	err := os.WriteFile(tmp.Name(), []byte("I like cake\nAnd pie."), 0777)
	if err != nil {
		t.Fatalf("failed to write to file")
	}
	defer os.Remove(tmp.Name())

	str := fileLinesFn(ENV, []primitive.Primitive{
		primitive.String(tmp.Name())})

	// Will lead to a list
	lst, ok2 := str.(primitive.List)
	if !ok2 {
		t.Fatalf("expected list, got %v", out)
	}

	if len(lst) != 2 {
		t.Fatalf("re-reading the temporary file gave bogus contents")
	}
}

// TestFileRead tests file:read
func TestFileRead(t *testing.T) {

	// calling with no argument
	out := fileReadFn(ENV, []primitive.Primitive{})

	// Will lead to an error
	_, ok := out.(primitive.Error)
	if !ok {
		t.Fatalf("expected error, got %v", out)
	}

	// Call with the wrong type
	out = fileReadFn(ENV, []primitive.Primitive{
		primitive.Number(3)})

	_, ok = out.(primitive.Error)
	if !ok {
		t.Fatalf("expected error, got %v", out)
	}

	// Call with a file that doesn't exist
	out = fileReadFn(ENV, []primitive.Primitive{
		primitive.String("path/not/found")})

	_, ok = out.(primitive.Error)
	if !ok {
		t.Fatalf("expected error, got %v", out)
	}

	// Create a temporary file, and read the contents
	tmp, _ := os.CreateTemp("", "yal")
	err := os.WriteFile(tmp.Name(), []byte("I like cake"), 0777)
	if err != nil {
		t.Fatalf("failed to write to file")
	}
	defer os.Remove(tmp.Name())

	str := fileReadFn(ENV, []primitive.Primitive{
		primitive.String(tmp.Name())})

	// Will lead to a string
	txt, ok2 := str.(primitive.String)
	if !ok2 {
		t.Fatalf("expected string, got %v", out)
	}

	if txt.ToString() != "I like cake" {
		t.Fatalf("re-reading the temporary file gave bogus contents")
	}
}

// TestFileStat tests file:stat
func TestFileStat(t *testing.T) {

	// calling with no argument
	out := fileStatFn(ENV, []primitive.Primitive{})

	// Will lead to an error
	_, ok := out.(primitive.Error)
	if !ok {
		t.Fatalf("expected error, got %v", out)
	}

	// Call with the wrong type
	out = fileStatFn(ENV, []primitive.Primitive{
		primitive.Number(3)})

	_, ok = out.(primitive.Error)
	if !ok {
		t.Fatalf("expected error, got %v", out)
	}

	// Call with a file that doesn't exist
	out = fileStatFn(ENV, []primitive.Primitive{
		primitive.String("path/not/found")})

	_, ok = out.(primitive.Nil)
	if !ok {
		t.Fatalf("expected nil, got %v", out)
	}

	// Create a temporary file
	tmp, _ := os.CreateTemp("", "yal")
	err := os.WriteFile(tmp.Name(), []byte("42 is the answer"), 0777)
	if err != nil {
		t.Fatalf("failed to write to file")
	}
	defer os.Remove(tmp.Name())

	// stat that
	str := fileStatFn(ENV, []primitive.Primitive{
		primitive.String(tmp.Name())})

	// Will lead to a list
	lst, ok2 := str.(primitive.List)
	if !ok2 {
		t.Fatalf("expected list, got %v", out)
	}

	// List will be: NAME SIZE ..
	if lst[1].ToString() != "16" {
		t.Fatalf("The size was wrong: %s", lst)
	}
}

// TestFileWrite tests file:write
func TestFileWrite(t *testing.T) {

	// calling with no argument
	out := fileWriteFn(ENV, []primitive.Primitive{})

	// Will lead to an error
	_, ok := out.(primitive.Error)
	if !ok {
		t.Fatalf("expected error, got %v", out)

	}

	// Call with the wrong type
	out = fileWriteFn(ENV, []primitive.Primitive{
		primitive.Number(3),
		primitive.String("cake")})

	_, ok = out.(primitive.Error)
	if !ok {
		t.Fatalf("expected error, got %v", out)
	}

	// Call with the wrong type
	out = fileWriteFn(ENV, []primitive.Primitive{
		primitive.String("/tmp/blah.txt"),
		primitive.Number(3)})

	_, ok = out.(primitive.Error)
	if !ok {
		t.Fatalf("expected error, got %v", out)
	}

	// Now we can try to write to something real
	// Create a temporary file to write to
	tmpfile, err := os.CreateTemp("", "exists")
	if err != nil {
		t.Fatalf("failed to create a temporary file")
	}
	defer os.Remove(tmpfile.Name())
	// Try to write there
	result := fileWriteFn(ENV, []primitive.Primitive{
		primitive.String(tmpfile.Name()),
		primitive.String("cake")})

	// Should be no errors
	_, ok = result.(primitive.Nil)
	if !ok {
		t.Fatalf("expected nil, got %v", out)
	}

	// Now create a temmporary directory, and use that
	// as the destination - which will fail once it is
	// removed
	path, err2 := os.MkdirTemp("", "directory_")
	if err2 != nil {
		t.Fatalf("failed to create temporary directory")
	}

	// path beneath the temporary directory
	target := filepath.Join(path, "foo.txt")

	// remove the direcotry
	os.RemoveAll(path)

	failure := fileWriteFn(ENV, []primitive.Primitive{
		primitive.String(target),
		primitive.String("cake")})

	// Should be no errors
	_, ok = failure.(primitive.Error)
	if !ok {
		t.Fatalf("expected failure writing beneath a directory which was removed, got %v", out)
	}

}

// TestGenSym tests gensym
func TestGenSym(t *testing.T) {

	// no arguments are required
	out := gensymFn(ENV, []primitive.Primitive{})

	// Will lead to a symbol
	_, ok := out.(primitive.Symbol)
	if !ok {
		t.Fatalf("expected symbol, got %v", out)
	}
}

// TestGet tests get
func TestGet(t *testing.T) {

	// no arguments
	out := getFn(ENV, []primitive.Primitive{})

	// Will lead to an error
	e, ok := out.(primitive.Error)
	if !ok {
		t.Fatalf("expected error, got %v", out)
	}
	if !strings.Contains(string(e), "argument") {
		t.Fatalf("got error, but wrong one")
	}

	// First argument must be a hash
	out = getFn(ENV, []primitive.Primitive{
		primitive.String("foo"),
		primitive.String("foo"),
	})

	// Will lead to an error
	e, ok = out.(primitive.Error)
	if !ok {
		t.Fatalf("expected error, got %v", out)
	}
	if !strings.Contains(string(e), "not a hash") {
		t.Fatalf("got error, but wrong one %v", out)
	}

	// create a hash
	h := primitive.NewHash()

	// Set a value
	h.Set("Name", primitive.String("STEVE"))

	// Now get it
	out2 := getFn(ENV, []primitive.Primitive{
		h,
		primitive.String("Name"),
	})

	// Will lead to a string
	s, ok2 := out2.(primitive.String)
	if !ok2 {
		t.Fatalf("expected string, got %v", out2)
	}
	if !strings.Contains(string(s), "STEVE") {
		t.Fatalf("got string, but wrong one %v", s)
	}
}

// TestGetenv tests getenv
func TestGetenv(t *testing.T) {

	// No arguments
	out := getenvFn(ENV, []primitive.Primitive{})

	// Will lead to an error
	e, ok := out.(primitive.Error)
	if !ok {
		t.Fatalf("expected error, got %v", out)
	}
	if e != primitive.ArityError() {
		t.Fatalf("got error, but wrong one %v", out)
	}

	// Argument that isn't a string
	out = getenvFn(ENV, []primitive.Primitive{
		primitive.Number(3),
	})

	// Will lead to an error
	e, ok = out.(primitive.Error)
	if !ok {
		t.Fatalf("expected error, got %v", out)
	}
	if !strings.Contains(string(e), "not a string") {
		t.Fatalf("got error, but wrong one %v", out)
	}

	// Valid result
	x := os.Getenv("USER")
	y := getenvFn(ENV, []primitive.Primitive{
		primitive.String("USER"),
	})

	yStr := string(y.(primitive.String))

	if yStr != x {
		t.Fatalf("getenv USER mismatch")
	}

}

// TestGlob tests glob
func TestGlob(t *testing.T) {

	// No arguments
	out := globFn(ENV, []primitive.Primitive{})

	// Will lead to an error
	e, ok := out.(primitive.Error)
	if !ok {
		t.Fatalf("expected error, got %v", out)
	}
	if !strings.Contains(string(e), "argument count") {
		t.Fatalf("got error, but wrong one %v", out)
	}

	// One argument, wrong type
	out = globFn(ENV, []primitive.Primitive{primitive.Number(33)})

	// Will lead to an error
	e, ok = out.(primitive.Error)
	if !ok {
		t.Fatalf("expected error, got %v", out)
	}
	if !strings.Contains(string(e), "not a string") {
		t.Fatalf("got error, but wrong one %v", out)
	}

	// Create a temporary directory
	path, err := os.MkdirTemp("", "directory_")
	if err != nil {
		t.Fatalf("failed to create a temporary directory")
	}

	// Populate two files.
	a := filepath.Join(path, "one.txt")
	b := filepath.Join(path, "two.foo")

	err = os.WriteFile(a, []byte("one.txt"), 0777)
	if err != nil {
		t.Fatalf("failed to write to file")
	}
	err = os.WriteFile(b, []byte("two.foo"), 0777)
	if err != nil {
		t.Fatalf("failed to write to file")
	}

	// Glob with two files
	files := globFn(ENV, []primitive.Primitive{primitive.String(path + "/*.*")})
	// Should get a list
	lst, ok2 := files.(primitive.List)
	if !ok2 {
		t.Fatalf("expected a list, got %v", files)
	}

	if len(lst) != 2 {
		t.Fatalf("expected two files, got %v", lst)
	}

	// Cleanup
	os.RemoveAll(path)

	// Finally an impossible glob
	out = globFn(ENV, []primitive.Primitive{primitive.String("[]")})

	// Will lead to an error
	e, ok = out.(primitive.Error)
	if !ok {
		t.Fatalf("expected error, got %v", out)
	}
	if !strings.Contains(string(e), "error in pattern") {
		t.Fatalf("got error, but wrong one %v", out)
	}

}

// TestHelp tests help
func TestHelp(t *testing.T) {
	// no arguments
	out := helpFn(ENV, []primitive.Primitive{})

	// Will lead to an error
	e, ok := out.(primitive.Error)
	if !ok {
		t.Fatalf("expected error, got %v", out)
	}
	if !strings.Contains(string(e), "argument") {
		t.Fatalf("got error, but wrong one")
	}

	// First argument must be a procedure
	out = helpFn(ENV, []primitive.Primitive{
		primitive.String("foo"),
	})

	// Will lead to an error
	e, ok = out.(primitive.Error)
	if !ok {
		t.Fatalf("expected error, got %v", out)
	}
	if !strings.Contains(string(e), "not a procedure") {
		t.Fatalf("got error, but wrong one %v", out)
	}

	//
	// create a new environment, and populate it
	//
	env := env.New()
	PopulateEnvironment(env)

	for _, name := range []string{"print", "sprintf", "length"} {

		// Load our standard library
		st := stdlib.Contents()
		std := string(st)

		// Create a new interpreter
		l := eval.New(std + "\n")

		l.Evaluate(env)

		fn, ok := env.Get(name)
		if !ok {
			t.Fatalf("failed to lookup function %s in environment", name)
		}

		result := helpFn(ENV, []primitive.Primitive{fn.(*primitive.Procedure)})

		txt, ok2 := result.(primitive.String)
		if !ok2 {
			t.Fatalf("expected a string, got %v", result)
		}
		if !strings.Contains(txt.ToString(), name) {
			t.Fatalf("got help text, but didn't find expected content: %v", result)
		}
	}
}

// TestInequality tests /=
func TestInequality(t *testing.T) {

	// No arguments
	out := inequalityFn(ENV, []primitive.Primitive{})

	// Will lead to an error
	e, ok := out.(primitive.Error)
	if !ok {
		t.Fatalf("expected error, got %v", out)
	}
	if e != primitive.ArityError() {
		t.Fatalf("got error, but wrong one %v", out)
	}

	//
	// Now a real one: unequal
	//
	out = inequalityFn(ENV, []primitive.Primitive{
		primitive.Number(9),
		primitive.Number(8),
	})

	// Will work
	n, ok2 := out.(primitive.Bool)
	if !ok2 {
		t.Fatalf("expected bool, got %v", out)
	}
	if n != true {
		t.Fatalf("got wrong result")
	}

	//
	// Now a real one: unequal - but multiple values
	//
	out = inequalityFn(ENV, []primitive.Primitive{
		primitive.Number(1),
		primitive.Number(2),
		primitive.Number(3),
		primitive.Number(4),
		primitive.Number(5),
		primitive.Number(6),
	})

	// Will work
	n, ok2 = out.(primitive.Bool)
	if !ok2 {
		t.Fatalf("expected bool, got %v", out)
	}
	if n != true {
		t.Fatalf("got wrong result")
	}

	//
	// Now a real one: unequal values
	//
	out = inequalityFn(ENV, []primitive.Primitive{
		primitive.Number(1),
		primitive.Number(2),
		primitive.Number(2),
		primitive.Number(1),
	})

	// Will work
	n, ok2 = out.(primitive.Bool)
	if !ok2 {
		t.Fatalf("expected bool, got %v", out)
	}
	if n != false {
		t.Fatalf("got wrong result")
	}

	//
	// Now with wrong types - last one is wrong
	//
	out = inequalityFn(ENV, []primitive.Primitive{
		primitive.Number(1),
		primitive.Number(2),
		primitive.Number(3),
		primitive.Number(4),
		primitive.String("5"),
	})

	e, ok = out.(primitive.Error)
	if !ok {
		t.Fatalf("expected error, got %v", out)
	}
	if !strings.Contains(string(e), "was not a number") {
		t.Fatalf("got error, but wrong one '%v'", e)
	}

	//
	// Now with wrong types
	//
	out = inequalityFn(ENV, []primitive.Primitive{
		primitive.String("9"),
		primitive.Number(9),
	})

	// Will report an error
	e, ok = out.(primitive.Error)
	if !ok {
		t.Fatalf("expected error, got %v", out)
	}
	if !strings.Contains(string(e), "was not a number") {
		t.Fatalf("got error, but wrong one %v", out)
	}
}

// TestJoin tests join
func TestJoin(t *testing.T) {

	// No arguments
	out := joinFn(ENV, []primitive.Primitive{})

	// Will lead to an error
	e, ok := out.(primitive.Error)
	if !ok {
		t.Fatalf("expected error, got %v", out)
	}
	if e != primitive.ArityError() {
		t.Fatalf("got error, but wrong one %v", out)
	}

	// Not a list
	out = joinFn(ENV, []primitive.Primitive{
		primitive.String("s"),
	})

	// Will lead to an error
	e, ok = out.(primitive.Error)
	if !ok {
		t.Fatalf("expected error, got %v", out)
	}
	if !strings.Contains(string(e), "not a list") {
		t.Fatalf("got error, but wrong one %v", out)
	}

	// Now a list
	out = joinFn(ENV, []primitive.Primitive{
		primitive.List{
			primitive.Number(3),
			primitive.Number(4),
		},
	})

	s, ok2 := out.(primitive.String)
	if !ok2 {
		t.Fatalf("expected string, got %v", s)
	}
	if s != "34" {
		t.Fatalf("got wrong result %v", s)
	}
}

// TestKeys tests keys
func TestKeys(t *testing.T) {

	// no arguments
	out := keysFn(ENV, []primitive.Primitive{})

	// Will lead to an error
	e, ok := out.(primitive.Error)
	if !ok {
		t.Fatalf("expected error, got %v", out)
	}
	if !strings.Contains(string(e), "argument") {
		t.Fatalf("got error, but wrong one")
	}

	// First argument must be a hash
	out = keysFn(ENV, []primitive.Primitive{
		primitive.String("foo"),
	})

	// Will lead to an error
	e, ok = out.(primitive.Error)
	if !ok {
		t.Fatalf("expected error, got %v", out)
	}
	if !strings.Contains(string(e), "not a hash") {
		t.Fatalf("got error, but wrong one %v", out)
	}

	// create a hash
	h := primitive.NewHash()
	h.Set("XXX", primitive.String("Last"))
	h.Set("Name", primitive.String("Steve"))
	h.Set("Age", primitive.Number(43))
	h.Set("Location", primitive.String("Helsinki"))

	// Get the keys
	res := keysFn(ENV, []primitive.Primitive{
		h,
	})

	// Will lead to a list
	_, ok2 := res.(primitive.List)
	if !ok2 {
		t.Fatalf("expected list, got %v", res)
	}

	// Sorted
	lst := res.(primitive.List)
	if lst[0].ToString() != "Age" {
		t.Fatalf("not a sorted list?")
	}
	if lst[1].ToString() != "Location" {
		t.Fatalf("not a sorted list?")
	}
	if lst[2].ToString() != "Name" {
		t.Fatalf("not a sorted list?")
	}
	if lst[3].ToString() != "XXX" {
		t.Fatalf("not a sorted list?")
	}
}

// TestList tests list
func TestList(t *testing.T) {

	// No arguments
	out := listFn(ENV, []primitive.Primitive{})

	// No error
	e, ok := out.(primitive.List)
	if !ok {
		t.Fatalf("expected error, got %v", out)
	}
	if e.ToString() != "()" {
		t.Fatalf("unexpected output %v", out)
	}

	// Two arguments
	out = listFn(ENV, []primitive.Primitive{
		primitive.Number(3),
		primitive.Number(43),
	})

	// No error
	e, ok = out.(primitive.List)
	if !ok {
		t.Fatalf("expected error, got %v", out)
	}
	if e.ToString() != "(3 43)" {
		t.Fatalf("unexpected output %v", out)
	}
}

// TestLt tests "<"
func TestLt(t *testing.T) {

	// No arguments
	out := ltFn(ENV, []primitive.Primitive{})

	// Will lead to an error
	e, ok := out.(primitive.Error)
	if !ok {
		t.Fatalf("expected error, got %v", out)
	}
	if e != primitive.ArityError() {
		t.Fatalf("got error, but wrong one %v", out)
	}

	// Argument which isn't a number
	out = ltFn(ENV, []primitive.Primitive{
		primitive.String("foo"),
		primitive.String("foo"),
	})

	// Will lead to an error
	e, ok = out.(primitive.Error)
	if !ok {
		t.Fatalf("expected error, got %v", out)
	}
	if !strings.Contains(string(e), "not a number") {
		t.Fatalf("got error, but wrong one %v", out)
	}

	// Argument which isn't a number
	out = ltFn(ENV, []primitive.Primitive{
		primitive.Number(32),
		primitive.String("foo"),
	})

	// Will lead to an error
	e, ok = out.(primitive.Error)
	if !ok {
		t.Fatalf("expected error, got %v", out)
	}
	if !strings.Contains(string(e), "not a number") {
		t.Fatalf("got error, but wrong one %v", out)
	}

	//
	// Now a real one
	//
	out = ltFn(ENV, []primitive.Primitive{
		primitive.Number(9),
		primitive.Number(100),
	})

	// Will work
	n, ok2 := out.(primitive.Bool)
	if !ok2 {
		t.Fatalf("expected bool, got %v", out)
	}
	if n != true {
		t.Fatalf("got wrong result")
	}
}

// TestMatches tests match
func TestMatches(t *testing.T) {

	// no arguments
	out := matchFn(ENV, []primitive.Primitive{})

	// Will lead to an error
	e, ok := out.(primitive.Error)
	if !ok {
		t.Fatalf("expected error, got %v", out)
	}
	if e != primitive.ArityError() {
		t.Fatalf("got error, but wrong one")
	}

	// First argument must be a string
	out = matchFn(ENV, []primitive.Primitive{
		primitive.Number(3),
		primitive.Number(4),
	})

	// Will lead to an error
	e, ok = out.(primitive.Error)
	if !ok {
		t.Fatalf("expected error, got %v", out)
	}
	if !strings.Contains(string(e), "not a string") {
		t.Fatalf("got error, but wrong one %v", out)
	}

	// Regexp must be valid
	out = matchFn(ENV, []primitive.Primitive{
		primitive.String("+"),
		primitive.Number(4),
	})

	// Will lead to an error
	e, ok = out.(primitive.Error)
	if !ok {
		t.Fatalf("expected error, got %v", out)
	}
	if !strings.Contains(string(e), "error parsing regexp") {
		t.Fatalf("got error, but wrong one %v", out)
	}

	// Now we have a valid call: no match
	fail := matchFn(ENV, []primitive.Primitive{
		primitive.String("foo"),
		primitive.String("bar"),
	})

	_, ok = fail.(primitive.Nil)
	if !ok {
		t.Fatalf("expected nil, got %v", out)
	}

	// Now we have a valid call: a match
	res := matchFn(ENV, []primitive.Primitive{
		primitive.String("[Ff]ood"),
		primitive.String("Food"),
	})

	// The list should have one entry
	lst, ok2 := res.(primitive.List)
	if !ok2 {
		t.Fatalf("expected a list, got %v", out)
	}
	if len(lst) != 1 {
		t.Fatalf("unexpected list size")
	}

	// Now we have a valid call: a match with capture group
	res = matchFn(ENV, []primitive.Primitive{
		primitive.String("([a-z]+)\\s*=\\s*([a-z]+)"),
		primitive.String("key = value"),
	})

	// The list should have three entries
	lst, ok2 = res.(primitive.List)
	if !ok2 {
		t.Fatalf("expected a list, got %v", out)
	}
	if len(lst) != 3 {
		t.Fatalf("unexpected list size")
	}

	if lst[0].ToString() != "key = value" {
		t.Fatalf("bogus match result")
	}
	if lst[1].ToString() != "key" {
		t.Fatalf("bogus match result")
	}
	if lst[2].ToString() != "value" {
		t.Fatalf("bogus match result")
	}
}

// TestMinus tests "-"
func TestMinus(t *testing.T) {

	// No arguments
	out := minusFn(ENV, []primitive.Primitive{})

	// Will lead to an error
	e, ok := out.(primitive.Error)
	if !ok {
		t.Fatalf("expected error, got %v", out)
	}
	if e != primitive.ArityError() {
		t.Fatalf("got error, but wrong one %v", out)
	}

	// Argument which isn't a number
	out = minusFn(ENV, []primitive.Primitive{
		primitive.String("foo"),
	})

	// Will lead to an error
	e, ok = out.(primitive.Error)
	if !ok {
		t.Fatalf("expected error, got %v", out)
	}
	if !strings.Contains(string(e), "not a number") {
		t.Fatalf("got error, but wrong one %v", out)
	}

	// Argument which isn't a number
	out = minusFn(ENV, []primitive.Primitive{
		primitive.Number(32),
		primitive.String("foo"),
	})

	// Will lead to an error
	e, ok = out.(primitive.Error)
	if !ok {
		t.Fatalf("expected error, got %v", out)
	}
	if !strings.Contains(string(e), "not a number") {
		t.Fatalf("got error, but wrong one %v", out)
	}

	//
	// Now a real one
	//
	out = minusFn(ENV, []primitive.Primitive{
		primitive.Number(10),
		primitive.Number(3),
	})

	// Will work
	n, ok2 := out.(primitive.Number)
	if !ok2 {
		t.Fatalf("expected number, got %v", out)
	}
	if n != 7 {
		t.Fatalf("got wrong result")
	}
}

// TestMod tests "%"
func TestMod(t *testing.T) {

	// No arguments
	out := modFn(ENV, []primitive.Primitive{})

	// Will lead to an error
	e, ok := out.(primitive.Error)
	if !ok {
		t.Fatalf("expected error, got %v", out)
	}
	if e != primitive.ArityError() {
		t.Fatalf("got error, but wrong one %v", out)
	}

	// Argument which isn't a number
	out = modFn(ENV, []primitive.Primitive{
		primitive.String("foo"),
		primitive.String("foo"),
	})

	// Will lead to an error
	e, ok = out.(primitive.Error)
	if !ok {
		t.Fatalf("expected error, got %v", out)
	}
	if !strings.Contains(string(e), "not a number") {
		t.Fatalf("got error, but wrong one %v", out)
	}

	// Argument which isn't a number
	out = modFn(ENV, []primitive.Primitive{
		primitive.Number(32),
		primitive.String("foo"),
	})

	// Will lead to an error
	e, ok = out.(primitive.Error)
	if !ok {
		t.Fatalf("expected error, got %v", out)
	}
	if !strings.Contains(string(e), "not a number") {
		t.Fatalf("got error, but wrong one %v", out)
	}

	//
	// Mod 0
	//
	out = modFn(ENV, []primitive.Primitive{
		primitive.Number(32),
		primitive.Number(0),
	})

	// Will lead to an error
	e, ok = out.(primitive.Error)
	if !ok {
		t.Fatalf("expected error, got %v", out)
	}
	if !strings.Contains(string(e), "division by zero") {
		t.Fatalf("got error, but wrong one %v", out)
	}

	//
	// Now a real one
	//
	out = modFn(ENV, []primitive.Primitive{
		primitive.Number(12),
		primitive.Number(3),
	})

	// Will work
	n, ok2 := out.(primitive.Number)
	if !ok2 {
		t.Fatalf("expected number, got %v", out)
	}
	if n != 0 {
		t.Fatalf("got wrong result")
	}
}

func TestMs(t *testing.T) {

	// No arguments
	out := msFn(ENV, []primitive.Primitive{})

	// Will lead to a number
	e, ok := out.(primitive.Number)
	if !ok {
		t.Fatalf("expected number, got %v", out)
	}

	// Get the current time
	tm := int(time.Now().UnixNano() / int64(time.Millisecond))

	if math.Abs(float64(tm-int(e))) > 10 {
		t.Fatalf("weird result; (ms) != ms - outside our bound of ten seconds inaccuracy")
	}

}

// TestMultiply tests "*"
func TestMultiply(t *testing.T) {

	// No arguments
	out := multiplyFn(ENV, []primitive.Primitive{})

	// Will lead to an error
	e, ok := out.(primitive.Error)
	if !ok {
		t.Fatalf("expected error, got %v", out)
	}
	if e != primitive.ArityError() {
		t.Fatalf("got error, but wrong one %v", out)
	}

	// Argument which isn't a number
	out = multiplyFn(ENV, []primitive.Primitive{
		primitive.String("foo"),
	})

	// Will lead to an error
	e, ok = out.(primitive.Error)
	if !ok {
		t.Fatalf("expected error, got %v", out)
	}
	if !strings.Contains(string(e), "not a number") {
		t.Fatalf("got error, but wrong one %v", out)
	}

	// Argument which isn't a number
	out = multiplyFn(ENV, []primitive.Primitive{
		primitive.Number(32),
		primitive.String("foo"),
	})

	// Will lead to an error
	e, ok = out.(primitive.Error)
	if !ok {
		t.Fatalf("expected error, got %v", out)
	}
	if !strings.Contains(string(e), "not a number") {
		t.Fatalf("got error, but wrong one %v", out)
	}

	//
	// Now a real one
	//
	out = multiplyFn(ENV, []primitive.Primitive{
		primitive.Number(10),
		primitive.Number(3),
	})

	// Will work
	n, ok2 := out.(primitive.Number)
	if !ok2 {
		t.Fatalf("expected number, got %v", out)
	}
	if n != 30 {
		t.Fatalf("got wrong result")
	}
}

// test nil?
func TestNil(t *testing.T) {

	// No arguments
	out := nilFn(ENV, []primitive.Primitive{})

	// Will lead to an error
	e, ok := out.(primitive.Error)
	if !ok {
		t.Fatalf("expected error, got %v", out)
	}
	if e != primitive.ArityError() {
		t.Fatalf("got error, but wrong one %v", out)
	}

	// nil is nil
	out = nilFn(ENV, []primitive.Primitive{
		primitive.Nil{},
	})

	// Will lead to a bool
	b, ok2 := out.(primitive.Bool)
	if !ok2 {
		t.Fatalf("unexpected type, expected bool, got %v", out)
	}
	if !b {
		t.Fatalf("wrong result")
	}

	// empty list is nil
	out = nilFn(ENV, []primitive.Primitive{
		primitive.List{},
	})

	// Will lead to a bool
	b, ok2 = out.(primitive.Bool)
	if !ok2 {
		t.Fatalf("unexpected type, expected bool, got %v", out)
	}
	if !b {
		t.Fatalf("wrong result")
	}

	// Finally a number is not a nil
	out = nilFn(ENV, []primitive.Primitive{
		primitive.Number(32),
	})

	// Will lead to a bool
	b, ok2 = out.(primitive.Bool)
	if !ok2 {
		t.Fatalf("unexpected type, expected bool, got %v", out)
	}
	if b {
		t.Fatalf("wrong result")
	}
}

func TestNow(t *testing.T) {

	// No arguments
	out := nowFn(ENV, []primitive.Primitive{})

	// Will lead to a number
	e, ok := out.(primitive.Number)
	if !ok {
		t.Fatalf("expected number, got %v", out)
	}

	// Get the current time
	tm := time.Now().Unix()

	if math.Abs(float64(tm-int64(e))) > 10 {
		t.Fatalf("weird result; (now) != now - outside our bound of ten seconds inaccuracy")
	}

}

func TestOrd(t *testing.T) {

	// no arguments
	out := ordFn(ENV, []primitive.Primitive{})

	// Will lead to an error
	e, ok := out.(primitive.Error)
	if !ok {
		t.Fatalf("expected error, got %v", out)
	}
	if e != primitive.ArityError() {
		t.Fatalf("got error, but wrong one:%s", out)
	}

	// First argument must be a string
	out = ordFn(ENV, []primitive.Primitive{
		primitive.Number(4),
	})

	// Will lead to an error
	e, ok = out.(primitive.Error)
	if !ok {
		t.Fatalf("expected error, got %v", out)
	}
	if !strings.Contains(string(e), "not a character/string") {
		t.Fatalf("got error, but wrong one %v", out)
	}

	// Now a valid call: * => 42
	val := ordFn(ENV, []primitive.Primitive{
		primitive.String("*"),
	})

	r, ok2 := val.(primitive.Number)
	if !ok2 {
		t.Fatalf("expected number, got %v", val)
	}
	if r.ToString() != "42" {
		t.Fatalf("got wrong result %v", r)
	}

	// Now a valid call: empty string => 0
	val = ordFn(ENV, []primitive.Primitive{
		primitive.String(""),
	})

	r, ok2 = val.(primitive.Number)
	if !ok2 {
		t.Fatalf("expected number, got %v", val)
	}
	if r.ToString() != "0" {
		t.Fatalf("got wrong result %v", r)
	}
}

func TestOs(t *testing.T) {

	// No arguments
	out := osFn(ENV, []primitive.Primitive{})

	// Will lead to a number
	e, ok := out.(primitive.String)
	if !ok {
		t.Fatalf("expected string, got %v", out)
	}

	if e.ToString() != runtime.GOOS {
		t.Fatalf("got wrong value for runtime OS")
	}
}

// TestPlus tests "+"
func TestPlus(t *testing.T) {

	// No arguments
	out := plusFn(ENV, []primitive.Primitive{})

	// Will lead to an error
	e, ok := out.(primitive.Error)
	if !ok {
		t.Fatalf("expected error, got %v", out)
	}
	if e != primitive.ArityError() {
		t.Fatalf("got error, but wrong one %v", out)
	}

	// Argument which isn't a number
	out = plusFn(ENV, []primitive.Primitive{
		primitive.String("foo"),
	})

	// Will lead to an error
	e, ok = out.(primitive.Error)
	if !ok {
		t.Fatalf("expected error, got %v", out)
	}
	if !strings.Contains(string(e), "not a number") {
		t.Fatalf("got error, but wrong one %v", out)
	}

	// Argument which isn't a number
	out = plusFn(ENV, []primitive.Primitive{
		primitive.Number(32),
		primitive.String("foo"),
	})

	// Will lead to an error
	e, ok = out.(primitive.Error)
	if !ok {
		t.Fatalf("expected error, got %v", out)
	}
	if !strings.Contains(string(e), "not a number") {
		t.Fatalf("got error, but wrong one %v", out)
	}

	//
	// Now a real one
	//
	out = plusFn(ENV, []primitive.Primitive{
		primitive.Number(10),
		primitive.Number(3),
	})

	// Will work
	n, ok2 := out.(primitive.Number)
	if !ok2 {
		t.Fatalf("expected number, got %v", out)
	}
	if n != 13 {
		t.Fatalf("got wrong result")
	}
}

func TestPrint(t *testing.T) {

	// No arguments
	out := printFn(ENV, []primitive.Primitive{})

	// Will lead to an error
	e, ok := out.(primitive.Error)
	if !ok {
		t.Fatalf("expected error, got %v", out)
	}
	if e != primitive.ArityError() {
		t.Fatalf("got error, but wrong one %v", out)
	}

	// One argument
	out = printFn(ENV, []primitive.Primitive{
		primitive.String("Hello!"),
	})

	e2, ok2 := out.(primitive.String)
	if !ok2 {
		t.Fatalf("expected string, got %v", out)
	}
	if e2 != "Hello!" {
		t.Fatalf("got error, but wrong one %v", e2)
	}

	// Two argument
	out = printFn(ENV, []primitive.Primitive{
		primitive.String("Hello %s!"),
		primitive.String("Steve"),
	})

	e2, ok2 = out.(primitive.String)
	if !ok2 {
		t.Fatalf("expected string, got %v", out)
	}
	if e2 != "Hello Steve!" {
		t.Fatalf("got error, but wrong one %v", e2)
	}
}

// TestRandom tests (random)
func TestRandom(t *testing.T) {

	// No arguments
	out := randomFn(ENV, []primitive.Primitive{})

	// Will lead to an error
	e, ok := out.(primitive.Error)
	if !ok {
		t.Fatalf("expected error, got %v", out)
	}
	if e != primitive.ArityError() {
		t.Fatalf("got error, but wrong one %v", out)
	}

	// One argument, of the wrong type
	out = randomFn(ENV, []primitive.Primitive{
		primitive.String("Hello!"),
	})

	e, ok = out.(primitive.Error)
	if !ok {
		t.Fatalf("expected error, got %v", out)
	}
	if !strings.Contains(string(e), "not a number") {
		t.Fatalf("got error, but wrong one %v", out)
	}

	// Now call with a number
	out = randomFn(ENV, []primitive.Primitive{
		primitive.Number(1),
	})
	_, ok2 := out.(primitive.Number)
	if !ok2 {
		t.Fatalf("expected string, got %v", out)
	}
}

// TestSet tests set
func TestSet(t *testing.T) {

	// no arguments
	out := setFn(ENV, []primitive.Primitive{})

	// Will lead to an error
	e, ok := out.(primitive.Error)
	if !ok {
		t.Fatalf("expected error, got %v", out)
	}
	if !strings.Contains(string(e), "argument") {
		t.Fatalf("got error, but wrong one")
	}

	// First argument must be a hash
	out = setFn(ENV, []primitive.Primitive{
		primitive.String("foo"),
		primitive.String("foo"),
		primitive.String("foo"),
	})

	// Will lead to an error
	e, ok = out.(primitive.Error)
	if !ok {
		t.Fatalf("expected error, got %v", out)
	}
	if !strings.Contains(string(e), "not a hash") {
		t.Fatalf("got error, but wrong one %v", out)
	}

	// create a hash
	h := primitive.NewHash()

	out2 := setFn(ENV, []primitive.Primitive{
		h,
		primitive.String("Name"),
		primitive.String("Steve"),
	})

	// Will lead to a string
	s, ok2 := out2.(primitive.String)
	if !ok2 {
		t.Fatalf("expected string, got %v", out2)
	}
	if !strings.Contains(string(s), "Steve") {
		t.Fatalf("got string, but wrong one %v", s)
	}

	// Now ensure the hash value was set
	v := h.Get("Name")
	if v.ToString() != "Steve" {
		t.Fatalf("The value wasn't set?")
	}
}

// TestSetup just instantiates the primitives in the environment
func TestSetup(t *testing.T) {

	// Create an empty environment
	e := env.New()

	// Before we start we have no functions
	_, ok := e.Get("print")
	if ok {
		t.Fatalf("didn't expect to get 'print' but did")
	}

	// Setup the builtins
	PopulateEnvironment(e)

	// Now we have functions
	_, ok = e.Get("print")
	if !ok {
		t.Fatalf("failed to find 'print' ")
	}
}

func TestSort(t *testing.T) {

	// No arguments
	out := sortFn(ENV, []primitive.Primitive{})

	// Will lead to an error
	e, ok := out.(primitive.Error)
	if !ok {
		t.Fatalf("expected error, got %v", out)
	}
	if e != primitive.ArityError() {
		t.Fatalf("got error, but wrong one %v", out)
	}

	// Not a list
	out = sortFn(ENV, []primitive.Primitive{
		primitive.Number(3),
	})

	e, ok = out.(primitive.Error)
	if !ok {
		t.Fatalf("expected error, got %v", out)
	}
	if !strings.Contains(string(e), "not a list") {
		t.Fatalf("got error, but wrong one %v", out)
	}

	//
	// Now we sort
	//
	out = sortFn(ENV, []primitive.Primitive{
		primitive.List{
			primitive.Number(30),
			primitive.Number(3),
			primitive.Number(-3),
		},
	})

	// Will lead to an error
	s, ok2 := out.(primitive.List)
	if !ok2 {
		t.Fatalf("expected list, got %v", out)
	}
	if s.ToString() != "(-3 3 30)" {
		t.Fatalf("got wrong result %v", s)
	}

	//
	// Now we sort a different range of things
	//
	out = sortFn(ENV, []primitive.Primitive{
		primitive.List{
			primitive.Bool(true),
			primitive.String("steve"),
			primitive.Number(3),
		},
	})

	s, ok2 = out.(primitive.List)
	if !ok2 {
		t.Fatalf("expected list, got %v", out)
	}
	if s.ToString() != "(#t 3 steve)" {
		t.Fatalf("got wrong result %v", s)
	}

}

func TestSplit(t *testing.T) {

	// No arguments
	out := splitFn(ENV, []primitive.Primitive{})

	// Will lead to an error
	e, ok := out.(primitive.Error)
	if !ok {
		t.Fatalf("expected error, got %v", out)
	}
	if e != primitive.ArityError() {
		t.Fatalf("got error, but wrong one %v", out)
	}

	// Arguments that aren't strings: 1
	out = splitFn(ENV, []primitive.Primitive{
		primitive.String("foo"),
		primitive.Number(3),
	})

	// Will lead to an error
	e, ok = out.(primitive.Error)
	if !ok {
		t.Fatalf("expected error, got %v", out)
	}
	if !strings.Contains(string(e), "not a string") {
		t.Fatalf("got error, but wrong one %v", out)
	}

	// Arguments that aren't strings: 2
	out = splitFn(ENV, []primitive.Primitive{
		primitive.Number(3),
		primitive.String("foo"),
	})

	// Will lead to an error
	e, ok = out.(primitive.Error)
	if !ok {
		t.Fatalf("expected error, got %v", out)
	}
	if !strings.Contains(string(e), "not a string") {
		t.Fatalf("got error, but wrong one %v", out)
	}

	//
	// Now a proper split
	//
	out = splitFn(ENV, []primitive.Primitive{
		primitive.String("fooπx"),
		primitive.String(""),
	})

	// Will lead to a list
	l, ok2 := out.(primitive.List)
	if !ok2 {
		t.Fatalf("expected list, got %v", out)
	}
	if len(l) != 5 {
		t.Fatalf("wrong length for result %v", l)
	}
	if l.ToString() != "(f o o π x)" {
		t.Fatalf("got wrong result %v", out)
	}

}

func TestSprintf(t *testing.T) {

	// No arguments
	out := sprintfFn(ENV, []primitive.Primitive{})

	// Will lead to an error
	e, ok := out.(primitive.Error)
	if !ok {
		t.Fatalf("expected error, got %v", out)
	}
	if e != primitive.ArityError() {
		t.Fatalf("got error, but wrong one %v", out)
	}

	// Two arguments
	out = sprintfFn(ENV, []primitive.Primitive{
		primitive.String("Hello\t\"%s\"\n\r!"),
		primitive.String("world"),
	})

	e2, ok2 := out.(primitive.String)
	if !ok2 {
		t.Fatalf("expected string, got %v", out)
	}
	if e2 != "Hello\t\"world\"\n\r!" {
		t.Fatalf("got wrong result %v", e2)
	}
}

func TestStr(t *testing.T) {

	// calling with no arguments will lead to an error
	fail := strFn(ENV, []primitive.Primitive{})

	// Will lead to an error
	_, ok := fail.(primitive.Error)
	if !ok {
		t.Fatalf("expected error, got %v", fail)
	}

	// calling with an arg
	out := strFn(ENV, []primitive.Primitive{
		primitive.Number(32),
	})

	// Will lead to an string
	e, ok := out.(primitive.String)
	if !ok {
		t.Fatalf("expected string, got %v", out)
	}
	if e != "32" {
		t.Fatalf("got wrong result %v", out)
	}
}

func TestType(t *testing.T) {

	// No arguments
	out := typeFn(ENV, []primitive.Primitive{})

	// Will lead to an error
	e, ok := out.(primitive.Error)
	if !ok {
		t.Fatalf("expected error, got %v", out)
	}
	if e != primitive.ArityError() {
		t.Fatalf("got error, but wrong one %v", out)
	}

	// calling with an arg
	out = typeFn(ENV, []primitive.Primitive{
		primitive.Number(32),
	})

	// Will lead to an string
	e2, ok2 := out.(primitive.String)
	if !ok2 {
		t.Fatalf("expected string, got %v", out)
	}
	if e2 != "number" {
		t.Fatalf("got wrong result %v", out)
	}
}
func TestVals(t *testing.T) {

	// no arguments
	out := valsFn(ENV, []primitive.Primitive{})

	// Will lead to an error
	e, ok := out.(primitive.Error)
	if !ok {
		t.Fatalf("expected error, got %v", out)
	}
	if !strings.Contains(string(e), "argument") {
		t.Fatalf("got error, but wrong one")
	}

	// First argument must be a hash
	out = valsFn(ENV, []primitive.Primitive{
		primitive.String("foo"),
	})

	// Will lead to an error
	e, ok = out.(primitive.Error)
	if !ok {
		t.Fatalf("expected error, got %v", out)
	}
	if !strings.Contains(string(e), "not a hash") {
		t.Fatalf("got error, but wrong one %v", out)
	}

	// create a hash
	h := primitive.NewHash()
	h.Set("XXX", primitive.String("Last"))
	h.Set("Name", primitive.String("Steve"))
	h.Set("Age", primitive.Number(43))
	h.Set("Location", primitive.String("Helsinki"))

	// Get the values
	res := valsFn(ENV, []primitive.Primitive{
		h,
	})

	// Will lead to a list
	_, ok2 := res.(primitive.List)
	if !ok2 {
		t.Fatalf("expected list, got %v", res)
	}

	lst := res.(primitive.List)
	if lst[0].ToString() != "43" {
		t.Fatalf("not a sorted list?")
	}
	if lst[1].ToString() != "Helsinki" {
		t.Fatalf("not a sorted list?")
	}
	if lst[2].ToString() != "Steve" {
		t.Fatalf("not a sorted list?")
	}
	if lst[3].ToString() != "Last" {
		t.Fatalf("not a sorted list?")
	}
}
