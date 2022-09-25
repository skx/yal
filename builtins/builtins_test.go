package builtins

import (
	"math"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/skx/yal/env"
	"github.com/skx/yal/primitive"
)

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

// TestPlus tests "+"
func TestPlus(t *testing.T) {

	// No arguments
	out := plusFn([]primitive.Primitive{})

	// Will lead to an error
	e, ok := out.(primitive.Error)
	if !ok {
		t.Fatalf("expected error, got %v", out)
	}
	if !strings.Contains(string(e), "invalid argument count") {
		t.Fatalf("got error, but wrong one %v", out)
	}

	// Argument which isn't a number
	out = plusFn([]primitive.Primitive{
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
	out = plusFn([]primitive.Primitive{
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
	out = plusFn([]primitive.Primitive{
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

// TestMinus tests "-"
func TestMinus(t *testing.T) {

	// No arguments
	out := minusFn([]primitive.Primitive{})

	// Will lead to an error
	e, ok := out.(primitive.Error)
	if !ok {
		t.Fatalf("expected error, got %v", out)
	}
	if !strings.Contains(string(e), "invalid argument count") {
		t.Fatalf("got error, but wrong one %v", out)
	}

	// Argument which isn't a number
	out = minusFn([]primitive.Primitive{
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
	out = minusFn([]primitive.Primitive{
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
	out = minusFn([]primitive.Primitive{
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

// TestMultiply tests "*"
func TestMultiply(t *testing.T) {

	// No arguments
	out := multiplyFn([]primitive.Primitive{})

	// Will lead to an error
	e, ok := out.(primitive.Error)
	if !ok {
		t.Fatalf("expected error, got %v", out)
	}
	if !strings.Contains(string(e), "invalid argument count") {
		t.Fatalf("got error, but wrong one %v", out)
	}

	// Argument which isn't a number
	out = multiplyFn([]primitive.Primitive{
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
	out = multiplyFn([]primitive.Primitive{
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
	out = multiplyFn([]primitive.Primitive{
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

// TestDivide tests "*"
func TestDivide(t *testing.T) {

	// No arguments
	out := divideFn([]primitive.Primitive{})

	// Will lead to an error
	e, ok := out.(primitive.Error)
	if !ok {
		t.Fatalf("expected error, got %v", out)
	}
	if !strings.Contains(string(e), "invalid argument count") {
		t.Fatalf("got error, but wrong one %v", out)
	}

	// Argument which isn't a number
	out = divideFn([]primitive.Primitive{
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
	out = divideFn([]primitive.Primitive{
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
	out = divideFn([]primitive.Primitive{
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
	out = divideFn([]primitive.Primitive{
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
}

// TestMod tests "%"
func TestMod(t *testing.T) {

	// No arguments
	out := modFn([]primitive.Primitive{})

	// Will lead to an error
	e, ok := out.(primitive.Error)
	if !ok {
		t.Fatalf("expected error, got %v", out)
	}
	if !strings.Contains(string(e), "number of arguments") {
		t.Fatalf("got error, but wrong one %v", out)
	}

	// Argument which isn't a number
	out = modFn([]primitive.Primitive{
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
	out = modFn([]primitive.Primitive{
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
	out = modFn([]primitive.Primitive{
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
	out = modFn([]primitive.Primitive{
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

// TestExpn tests "#"
func TestExpn(t *testing.T) {

	// No arguments
	out := expnFn([]primitive.Primitive{})

	// Will lead to an error
	e, ok := out.(primitive.Error)
	if !ok {
		t.Fatalf("expected error, got %v", out)
	}
	if !strings.Contains(string(e), "number of arguments") {
		t.Fatalf("got error, but wrong one %v", out)
	}

	// Argument which isn't a number
	out = expnFn([]primitive.Primitive{
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
	out = expnFn([]primitive.Primitive{
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
	out = expnFn([]primitive.Primitive{
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

// TestEq tests "eq"
func TestEq(t *testing.T) {

	// No arguments
	out := eqFn([]primitive.Primitive{})

	// Will lead to an error
	e, ok := out.(primitive.Error)
	if !ok {
		t.Fatalf("expected error, got %v", out)
	}
	if !strings.Contains(string(e), "number of arguments") {
		t.Fatalf("got error, but wrong one %v", out)
	}

	//
	// Now a real one: equal
	//
	out = eqFn([]primitive.Primitive{
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
	out = eqFn([]primitive.Primitive{
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
	out = eqFn([]primitive.Primitive{
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

// TestEquals tests "="
func TestEquals(t *testing.T) {

	// No arguments
	out := equalsFn([]primitive.Primitive{})

	// Will lead to an error
	e, ok := out.(primitive.Error)
	if !ok {
		t.Fatalf("expected error, got %v", out)
	}
	if !strings.Contains(string(e), "number of arguments") {
		t.Fatalf("got error, but wrong one %v", out)
	}

	//
	// Now a real one: equal
	//
	out = equalsFn([]primitive.Primitive{
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
	out = equalsFn([]primitive.Primitive{
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
	// Now with wrong types
	//
	out = equalsFn([]primitive.Primitive{
		primitive.Number(9),
		primitive.String("9"),
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
	out = equalsFn([]primitive.Primitive{
		primitive.String("9"),
		primitive.Number(9),
	})

	// Will work
	e, ok = out.(primitive.Error)
	if !ok {
		t.Fatalf("expected error, got %v", out)
	}
	if !strings.Contains(string(e), "was not a number") {
		t.Fatalf("got error, but wrong one %v", out)
	}
}

// TestLt tests "<"
func TestLt(t *testing.T) {

	// No arguments
	out := ltFn([]primitive.Primitive{})

	// Will lead to an error
	e, ok := out.(primitive.Error)
	if !ok {
		t.Fatalf("expected error, got %v", out)
	}
	if !strings.Contains(string(e), "number of arguments") {
		t.Fatalf("got error, but wrong one %v", out)
	}

	// Argument which isn't a number
	out = ltFn([]primitive.Primitive{
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
	out = ltFn([]primitive.Primitive{
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
	out = ltFn([]primitive.Primitive{
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

func TestList(t *testing.T) {

	// No arguments
	out := listFn([]primitive.Primitive{})

	// No error
	e, ok := out.(primitive.List)
	if !ok {
		t.Fatalf("expected error, got %v", out)
	}
	if e.ToString() != "()" {
		t.Fatalf("unexpected output %v", out)
	}

	// Two arguments
	out = listFn([]primitive.Primitive{
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

// Test (car
func TestCar(t *testing.T) {

	// No arguments
	out := carFn([]primitive.Primitive{})

	// Will lead to an error
	e, ok := out.(primitive.Error)
	if !ok {
		t.Fatalf("expected error, got %v", out)
	}
	if !strings.Contains(string(e), "number of arguments") {
		t.Fatalf("got error, but wrong one %v", out)
	}

	// One argument
	out = carFn([]primitive.Primitive{
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
	out = carFn([]primitive.Primitive{
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
	out = carFn([]primitive.Primitive{
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
	out := cdrFn([]primitive.Primitive{})

	// Will lead to an error
	e, ok := out.(primitive.Error)
	if !ok {
		t.Fatalf("expected error, got %v", out)
	}
	if !strings.Contains(string(e), "number of arguments") {
		t.Fatalf("got error, but wrong one %v", out)
	}

	// One argument
	out = cdrFn([]primitive.Primitive{
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
	out = cdrFn([]primitive.Primitive{
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
	out = cdrFn([]primitive.Primitive{
		primitive.List{},
	})

	// No error
	_, ok3 := out.(primitive.Nil)
	if !ok3 {
		t.Fatalf("expected nil, got %v", out)
	}
}

func TestStr(t *testing.T) {

	// calling with an arg
	out := strFn([]primitive.Primitive{
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
	out := typeFn([]primitive.Primitive{})

	// Will lead to an error
	e, ok := out.(primitive.Error)
	if !ok {
		t.Fatalf("expected error, got %v", out)
	}
	if !strings.Contains(string(e), "number of arguments") {
		t.Fatalf("got error, but wrong one %v", out)
	}

	// calling with an arg
	out = typeFn([]primitive.Primitive{
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

func TestError(t *testing.T) {

	// No arguments
	out := errorFn([]primitive.Primitive{})

	// Will lead to an error
	e, ok := out.(primitive.Error)
	if !ok {
		t.Fatalf("expected error, got %v", out)
	}
	if !strings.Contains(string(e), "number of arguments") {
		t.Fatalf("got error, but wrong one %v", out)
	}

	// calling with an arg
	out = errorFn([]primitive.Primitive{
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

// test nil?
func TestNil(t *testing.T) {

	// No arguments
	out := nilFn([]primitive.Primitive{})

	// Will lead to an error
	e, ok := out.(primitive.Error)
	if !ok {
		t.Fatalf("expected error, got %v", out)
	}
	if !strings.Contains(string(e), "number of arguments") {
		t.Fatalf("got error, but wrong one %v", out)
	}

	// nil is nil
	out = nilFn([]primitive.Primitive{
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
	out = nilFn([]primitive.Primitive{
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
	out = nilFn([]primitive.Primitive{
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

func TestCons(t *testing.T) {

	// No arguments
	out := consFn([]primitive.Primitive{})

	// Will lead to an error
	e, ok := out.(primitive.Error)
	if !ok {
		t.Fatalf("expected error, got %v", out)
	}
	if !strings.Contains(string(e), "wrong number of arguments") {
		t.Fatalf("got error, but wrong one %v", out)
	}

	// one argument, string -> list
	out = consFn([]primitive.Primitive{
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
	out = consFn([]primitive.Primitive{
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
	out = consFn(a)
	out, ok2 = out.(primitive.List)
	if !ok2 {
		t.Errorf("expected list")
	}
	if out.ToString() != "((3 4) 5)" {
		t.Fatalf("wrong result, got %v", out)
	}

	// second one
	out = consFn(b)
	out, ok2 = out.(primitive.List)
	if !ok2 {
		t.Errorf("expected list")
	}
	if out.ToString() != "(5 3 4)" {
		t.Fatalf("wrong result, got %v", out)
	}
}

func TestPrint(t *testing.T) {

	// No arguments
	out := printFn([]primitive.Primitive{})

	// Will lead to an error
	e, ok := out.(primitive.Error)
	if !ok {
		t.Fatalf("expected error, got %v", out)
	}
	if !strings.Contains(string(e), "wrong number of arguments") {
		t.Fatalf("got error, but wrong one %v", out)
	}

	// One argument
	out = printFn([]primitive.Primitive{
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
	out = printFn([]primitive.Primitive{
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

func TestSprintf(t *testing.T) {

	// No arguments
	out := sprintfFn([]primitive.Primitive{})

	// Will lead to an error
	e, ok := out.(primitive.Error)
	if !ok {
		t.Fatalf("expected error, got %v", out)
	}
	if !strings.Contains(string(e), "wrong number of arguments") {
		t.Fatalf("got error, but wrong one %v", out)
	}

	// Two arguments
	out = sprintfFn([]primitive.Primitive{
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

func TestJoin(t *testing.T) {

	// No arguments
	out := joinFn([]primitive.Primitive{})

	// Will lead to an error
	e, ok := out.(primitive.Error)
	if !ok {
		t.Fatalf("expected error, got %v", out)
	}
	if !strings.Contains(string(e), "invalid argument count") {
		t.Fatalf("got error, but wrong one %v", out)
	}

	// Not a list
	out = joinFn([]primitive.Primitive{
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
	out = joinFn([]primitive.Primitive{
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

func TestSplit(t *testing.T) {

	// No arguments
	out := splitFn([]primitive.Primitive{})

	// Will lead to an error
	e, ok := out.(primitive.Error)
	if !ok {
		t.Fatalf("expected error, got %v", out)
	}
	if !strings.Contains(string(e), "invalid argument count") {
		t.Fatalf("got error, but wrong one %v", out)
	}

	// Arguments that aren't strings: 1
	out = splitFn([]primitive.Primitive{
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
	out = splitFn([]primitive.Primitive{
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
	out = splitFn([]primitive.Primitive{
		primitive.String("foo"),
		primitive.String(""),
	})

	// Will lead to a list
	l, ok2 := out.(primitive.List)
	if !ok2 {
		t.Fatalf("expected list, got %v", out)
	}
	if l.ToString() != "(f o o)" {
		t.Fatalf("got wrong result %v", out)
	}

}

func TestSort(t *testing.T) {

	// No arguments
	out := sortFn([]primitive.Primitive{})

	// Will lead to an error
	e, ok := out.(primitive.Error)
	if !ok {
		t.Fatalf("expected error, got %v", out)
	}
	if !strings.Contains(string(e), "invalid argument count") {
		t.Fatalf("got error, but wrong one %v", out)
	}

	// Not a list
	out = sortFn([]primitive.Primitive{
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
	out = sortFn([]primitive.Primitive{
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
	out = sortFn([]primitive.Primitive{
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

func TestGetenv(t *testing.T) {

	// No arguments
	out := getenvFn([]primitive.Primitive{})

	// Will lead to an error
	e, ok := out.(primitive.Error)
	if !ok {
		t.Fatalf("expected error, got %v", out)
	}
	if !strings.Contains(string(e), "invalid argument count") {
		t.Fatalf("got error, but wrong one %v", out)
	}

	// Argument that isn't a string
	out = getenvFn([]primitive.Primitive{
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
	y := getenvFn([]primitive.Primitive{
		primitive.String("USER"),
	})

	yStr := string(y.(primitive.String))

	if yStr != x {
		t.Fatalf("getenv USER mismatch")
	}

}

func TestNow(t *testing.T) {

	// No arguments
	out := nowFn([]primitive.Primitive{})

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

func TestMs(t *testing.T) {

	// No arguments
	out := msFn([]primitive.Primitive{})

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

func TestGet(t *testing.T) {

	// no arguments
	out := getFn([]primitive.Primitive{})

	// Will lead to an error
	e, ok := out.(primitive.Error)
	if !ok {
		t.Fatalf("expected error, got %v", out)
	}
	if !strings.Contains(string(e), "argument") {
		t.Fatalf("got error, but wrong one")
	}

	// First argument must be a hash
	out = getFn([]primitive.Primitive{
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
	out2 := getFn([]primitive.Primitive{
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

func TestKeys(t *testing.T) {

	// no arguments
	out := keysFn([]primitive.Primitive{})

	// Will lead to an error
	e, ok := out.(primitive.Error)
	if !ok {
		t.Fatalf("expected error, got %v", out)
	}
	if !strings.Contains(string(e), "argument") {
		t.Fatalf("got error, but wrong one")
	}

	// First argument must be a hash
	out = keysFn([]primitive.Primitive{
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
	res := keysFn([]primitive.Primitive{
		h,
	})

	// Will lead to a string
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

func TestSet(t *testing.T) {

	// no arguments
	out := setFn([]primitive.Primitive{})

	// Will lead to an error
	e, ok := out.(primitive.Error)
	if !ok {
		t.Fatalf("expected error, got %v", out)
	}
	if !strings.Contains(string(e), "argument") {
		t.Fatalf("got error, but wrong one")
	}

	// First argument must be a hash
	out = setFn([]primitive.Primitive{
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

	out2 := setFn([]primitive.Primitive{
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

func TestMatches(t *testing.T) {

	// no arguments
	out := matchFn([]primitive.Primitive{})

	// Will lead to an error
	e, ok := out.(primitive.Error)
	if !ok {
		t.Fatalf("expected error, got %v", out)
	}
	if !strings.Contains(string(e), "invalid argument count") {
		t.Fatalf("got error, but wrong one")
	}

	// First argument must be a string
	out = matchFn([]primitive.Primitive{
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
	out = matchFn([]primitive.Primitive{
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
	fail := matchFn([]primitive.Primitive{
		primitive.String("foo"),
		primitive.String("bar"),
	})

	_, ok = fail.(primitive.Nil)
	if !ok {
		t.Fatalf("expected nil, got %v", out)
	}

	// Now we have a valid call: a match
	res := matchFn([]primitive.Primitive{
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
	res = matchFn([]primitive.Primitive{
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

func TestOrd(t *testing.T) {

	// no arguments
	out := ordFn([]primitive.Primitive{})

	// Will lead to an error
	e, ok := out.(primitive.Error)
	if !ok {
		t.Fatalf("expected error, got %v", out)
	}
	if !strings.Contains(string(e), "wrong number of arguments") {
		t.Fatalf("got error, but wrong one:%s", out)
	}

	// First argument must be a string
	out = ordFn([]primitive.Primitive{
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

	// Now a valid call: * => 42
	val := ordFn([]primitive.Primitive{
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
	val = ordFn([]primitive.Primitive{
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

func TestChr(t *testing.T) {

	// no arguments
	out := chrFn([]primitive.Primitive{})

	// Will lead to an error
	e, ok := out.(primitive.Error)
	if !ok {
		t.Fatalf("expected error, got %v", out)
	}
	if !strings.Contains(string(e), "wrong number of arguments") {
		t.Fatalf("got error, but wrong one:%s", out)
	}

	// First argument must be a number
	out = chrFn([]primitive.Primitive{
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
	val := chrFn([]primitive.Primitive{
		primitive.Number(42),
	})

	r, ok2 := val.(primitive.String)
	if !ok2 {
		t.Fatalf("expected string, got %v", val)
	}
	if r.ToString() != "*" {
		t.Fatalf("got wrong result %v", r)
	}

}
