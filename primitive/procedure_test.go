package primitive

import ( "testing"
	"github.com/skx/yal/env"
)

func TestProcedure(t *testing.T) {

	// built-in
	b := Procedure{
		F: func(e *env.Environment,args []Primitive) Primitive {
			return Nil{}
		},
	}

	// lisp
	l := Procedure{
		Args: []Symbol{
			Symbol("A"),
			Symbol("B"),
		},
		Body: List{
			Symbol("+"),
			Symbol("A"),
			Symbol("B"),
		},
	}

	// macro
	m := Procedure{
		Macro: true,
		Args: []Symbol{
			Symbol("A"),
			Symbol("B"),
		},
		Body: List{
			Symbol("+"),
			Symbol("A"),
			Symbol("B"),
		},
	}

	if b.Type() != "procedure(golang)" {
		t.Fatalf("wrong type for builtin")
	}
	if b.ToString() != "#built-in-function" {
		t.Fatalf("wrong string-type for builtin, got %s", b.ToString())
	}

	if l.Type() != "procedure(lisp)" {
		t.Fatalf("wrong type for lisp proc")
	}
	if l.ToString() != "(lambda (A B) (+ A B))" {
		t.Fatalf("wrong string-type for lisp-proc, got %s", l.ToString())
	}

	if m.Type() != "macro" {
		t.Fatalf("wrong type for lisp macro")
	}
	if m.ToString() != "(macro (A B) (+ A B))" {
		t.Fatalf("wrong string-type for macro, got %s", m.ToString())
	}
}
