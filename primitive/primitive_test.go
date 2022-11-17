package primitive

import (
	"fmt"
	"strings"
	"testing"
)

func TestBool(t *testing.T) {

	true := Bool(true)
	false := Bool(false)

	if !true.IsSimpleType() {
		t.Fatalf("expected boolean to be a simple type")
	}

	if true.Type() != "boolean" {
		t.Fatalf("wrong type")
	}
	if true.ToString() != "#t" {
		t.Fatalf("bool->String had wrong result")
	}
	if false.ToString() != "#f" {
		t.Fatalf("bool->String had wrong result")
	}

	ti := true.ToInterface()
	fi := false.ToInterface()

	bTrue, bOK := ti.(bool)
	if !bOK {
		t.Fatalf("bool.ToInterface did not result in a bool")
	}
	if !bTrue {
		t.Fatalf("ToInterface resulted in the wrong result")
	}

	bFalse, bOK2 := fi.(bool)
	if !bOK2 {
		t.Fatalf("bool.ToInterface did not result in a bool")
	}
	if bFalse {
		t.Fatalf("ToInterface resulted in the wrong result")
	}

}

func TestCharacter(t *testing.T) {

	nl := Character("\n")
	ok := Character("o")
	empty := Character("")

	if !nl.IsSimpleType() {
		t.Fatalf("expected character to be a simple type")
	}

	if nl.Type() != "character" {
		t.Fatalf("wrong type")
	}
	if nl.ToString() != "\n" {
		t.Fatalf("char->String had wrong result")
	}
	if ok.ToString() != "o" {
		t.Fatalf("char->String had wrong result")
	}

	nli := nl.ToInterface()
	emptyi := empty.ToInterface()

	nliGo, nliOK := nli.(uint8)
	if !nliOK {
		t.Fatalf("character.ToInterface gave wrong type %T", nli)
	}
	if nliGo != '\n' {
		t.Fatalf("ToInterface resulted in the wrong result")
	}

	emptyGo, emptyOK := emptyi.(string)
	if !emptyOK {
		t.Fatalf("character.ToInterface gave wrong type %T", emptyi)
	}
	if emptyGo != "" {
		t.Fatalf("ToInterface resulted in the wrong result")
	}
}

func TestError(t *testing.T) {

	error := Error("no-cheese")

	if !error.IsSimpleType() {
		t.Fatalf("expected error to be a simple type")
	}

	if error.Type() != "error" {
		t.Fatalf("wrong type")
	}
	if error.ToString() != "ERROR{no-cheese}" {
		t.Fatalf("error->String had wrong result")
	}

	if !strings.Contains(ArityError().ToString(), "Arity") {
		t.Fatalf("arity-error is non-obvious")
	}

	if !strings.Contains(TypeError("xx").ToString(), "TypeError") {
		t.Fatalf("TypeError is non-obvious")
	}

	//
	// TODO: This is horrid
	//
	errGo := error.ToInterface()
	if !strings.Contains(fmt.Sprintf("%s", errGo), "cheese") {
		t.Fatalf("error.ToInterface is non-obvious")
	}

}

func TestIsNil(t *testing.T) {

	var n Nil

	if !n.IsSimpleType() {
		t.Fatalf("expected nil to be a simple type")
	}

	if n.Type() != "nil" {
		t.Fatalf("nil -> wrong type")
	}
	if n.ToString() != "nil" {
		t.Fatalf("nil->string wrong result")
	}

	i := n.ToInterface()
	if i != nil {
		t.Fatalf("nil.ToInterface gave wrong type")
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

	if lst.IsSimpleType() {
		t.Fatalf("Did not expect list to be a simple type")
	}

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

	if !i.IsSimpleType() {
		t.Fatalf("expected number to be a simple type")
	}

	if i.Type() != "number" {
		t.Fatalf("wrong type")
	}
	if i.ToString() != "3" {
		t.Fatalf("number->String had wrong result")
	}
	if !(strings.Contains(f.ToString(), "0.111")) {
		t.Fatalf("number->String (float) had wrong result:%s", f.ToString())
	}

	ii := i.ToInterface()
	fi := f.ToInterface()

	iGo, iOK := ii.(int)
	if !iOK {
		t.Fatalf("Int.ToInterface gave wrong type %T", ii)
	}
	if iGo != 3 {
		t.Fatalf("ToInterface resulted in the wrong result")
	}

	fGo, fOK := fi.(float64)
	if !fOK {
		t.Fatalf("Int.ToInterface gave wrong type")
	}
	if !(strings.Contains(fmt.Sprintf("%f", fGo), "0.111")) {
		t.Fatalf("ToInterface resulted in the wrong result")
	}

}

func TestString(t *testing.T) {

	str := String("i like cake")

	if !str.IsSimpleType() {
		t.Fatalf("expected string to be a simple type")
	}

	if str.Type() != "string" {
		t.Fatalf("wrong type")
	}
	if str.ToString() != "i like cake" {
		t.Fatalf("string->String had wrong result")
	}
}

func TestSymbol(t *testing.T) {

	sym := Symbol("pi")

	if sym.IsSimpleType() {
		t.Fatalf("did not expected symbol to be a simple type")
	}
	if sym.Type() != "symbol" {
		t.Fatalf("wrong type")
	}
	if sym.ToString() != "pi" {
		t.Fatalf("symbol->String had wrong result")
	}

	si := sym.ToInterface()

	sGo, sOK := si.(string)
	if !sOK {
		t.Fatalf("String.ToInterface gave wrong type")
	}
	if sGo != "pi" {
		t.Fatalf("ToInterface resulted in the wrong result")
	}

}
