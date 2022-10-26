package primitive

import (
	"strings"
	"testing"
)

func TestBool(t *testing.T) {

	true := Bool(true)
	false := Bool(false)

	if true.Type() != "boolean" {
		t.Fatalf("wrong type")
	}
	if true.ToString() != "#t" {
		t.Fatalf("bool->String had wrong result")
	}
	if false.ToString() != "#f" {
		t.Fatalf("bool->String had wrong result")
	}
}

func TestCharacter(t *testing.T) {

	nl := Character("\n")
	ok := Character("o")

	if nl.Type() != "character" {
		t.Fatalf("wrong type")
	}
	if nl.ToString() != "\n" {
		t.Fatalf("char->String had wrong result")
	}
	if ok.ToString() != "o" {
		t.Fatalf("char->String had wrong result")
	}
}

func TestError(t *testing.T) {

	error := Error("no-cheese")

	if error.Type() != "error" {
		t.Fatalf("wrong type")
	}
	if error.ToString() != "ERROR{no-cheese}" {
		t.Fatalf("error->String had wrong result")
	}

	if !strings.Contains(ArityError().ToString(), "Arity") {
		t.Fatalf("arity-error is non-obvious")
	}

}

func TestIsNil(t *testing.T) {

	var n Nil

	if n.Type() != "nil" {
		t.Fatalf("nil -> wrong type")
	}
	if n.ToString() != "nil" {
		t.Fatalf("nil->string wrong result")
	}

	var s String
	var f Number
	var b Bool

	if !IsNil(n) {
		t.Fatalf("nil is supposed to be nil")
	}
	if IsNil(s) {
		t.Fatalf("a string is not nil")
	}
	if IsNil(f) {
		t.Fatalf("a number is not nil")
	}
	if IsNil(b) {
		t.Fatalf("a bool is not nil")
	}
}

func TestList(t *testing.T) {

	lst := List([]Primitive{
		Error("no-cheese"),
		Number(3),
	})

	if lst.Type() != "list" {
		t.Fatalf("wrong type")
	}
	if lst.ToString() != "(ERROR{no-cheese} 3)" {
		t.Fatalf("list->String had wrong result:%s", lst.ToString())
	}
}

func TestNumber(t *testing.T) {

	i := Number(3)
	f := Number(1.0 / 9)

	if i.Type() != "number" {
		t.Fatalf("wrong type")
	}
	if i.ToString() != "3" {
		t.Fatalf("number->String had wrong result")
	}
	if !(strings.Contains(f.ToString(), "0.111")) {
		t.Fatalf("number->String (float) had wrong result:%s", f.ToString())
	}
}

func TestString(t *testing.T) {

	str := String("i like cake")

	if str.Type() != "string" {
		t.Fatalf("wrong type")
	}
	if str.ToString() != "i like cake" {
		t.Fatalf("string->String had wrong result")
	}
}

func TestSymbol(t *testing.T) {

	sym := Symbol("pi")

	if sym.Type() != "symbol" {
		t.Fatalf("wrong type")
	}
	if sym.ToString() != "pi" {
		t.Fatalf("symbol->String had wrong result")
	}
}
