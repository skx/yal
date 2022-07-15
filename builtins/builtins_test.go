package builtins

import (
	"strings"
	"testing"

	"github.com/skx/yal/primitive"
)

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
	// Now a real one
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
