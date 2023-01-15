//go:build !windows
// +build !windows

package builtins

import (
	"os"
	"strings"
	"testing"

	"github.com/skx/yal/primitive"
)

// TestShell tests shell - but only on Linux/Unix
func TestShell(t *testing.T) {

	// calling with no argument
	out := shellFn(nil, ENV, []primitive.Primitive{})

	// Will lead to an error
	_, ok := out.(primitive.Error)
	if !ok {
		t.Fatalf("expected error, got %v", out)
	}

	// One argument, but the wrong type
	out = shellFn(nil, ENV, []primitive.Primitive{
		primitive.Number(3),
	})

	var e primitive.Primitive
	e, ok = out.(primitive.Error)
	if !ok {
		t.Fatalf("expected error, got %v", out)
	}
	if !strings.Contains(e.ToString(), "not a list") {
		t.Fatalf("got error, but wrong one %v", out)
	}

	// Echo command to execute.
	cmd := primitive.List{}
	cmd = append(cmd, primitive.String("echo"))
	cmd = append(cmd, primitive.String("foo"))
	cmd = append(cmd, primitive.String("bar"))

	// Run the command
	res := shellFn(nil, ENV, []primitive.Primitive{cmd})

	// Response should be a list
	lst, ok2 := res.(primitive.List)
	if !ok2 {
		t.Fatalf("expected (shell) to return a list, got %v", res)
	}

	// with two entries
	if len(lst) != 2 {
		t.Fatalf("expected (shell) result to have two entries, got %v", lst)
	}

	//
	// Now: run a command that will fail
	//
	fail := primitive.List{}
	fail = append(fail, primitive.String("/fdsf/fdsf/-path-not/exists"))

	// Run the command
	out = shellFn(nil, ENV, []primitive.Primitive{fail})

	// Will lead to an error
	_, ok = out.(primitive.Error)
	if !ok {
		t.Fatalf("expected error, got %v", out)
	}

	//
	// Now: Pretend we're running under a fuzzer
	//
	// Preserve any previous content of $FUZZ
	//
	old := os.Getenv("FUZZ")
	os.Setenv("FUZZ", "FUZZ")
	res = shellFn(nil, ENV, []primitive.Primitive{cmd})
	os.Setenv("FUZZ", old)

	// Response should still be a list
	lst, ok2 = res.(primitive.List)
	if !ok2 {
		t.Fatalf("expected (shell) to return a list, got %v", res)
	}

	// with zero entries
	if len(lst) != 0 {
		t.Fatalf("expected (shell) result to have zero entries, got %v", lst)
	}

}
