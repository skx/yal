package env

import "testing"

// TestGetSet tests get/set on a variable
func TestGetSet(t *testing.T) {

	e := New()

	// by default the environment is empty
	_, ok := e.Get("FOO")
	if ok {
		t.Fatalf("fetching missing variable shouldn't work")
	}

	// Now set
	e.Set("FOO", "BAR")
	out, ok2 := e.Get("FOO")
	if !ok2 {
		t.Fatalf("fetching variable shouldn't fail")
	}
	if out.(string) != "BAR" {
		t.Fatalf("variable had wrong value")
	}

}

func TestItems(t *testing.T) {

	// parent
	p := New()
	p.Set("FOO", "BAR")

	// child
	c := NewEnvironment(p)

	items := c.Items()
	if len(items) != 1 {
		t.Fatalf("wrong number of items found")
	}

	// set in the child
	c.Set("FOO", "BART")
	items = c.Items()
	if len(items) != 1 {
		t.Fatalf("wrong number of items found")
	}

	if items["FOO"] != "BART" {
		t.Fatalf("wrong value in items")
	}

	// set in parent
	p.Set("NAME", "STEVE")

	items = c.Items()
	if len(items) != 2 {
		t.Fatalf("wrong item count")
	}
}

func TestScopedSet(t *testing.T) {

	// parent
	p := New()
	p.Set("FOO", "BAR")

	// child
	c := NewEnvironment(p)

	// Child should be able to reach parent variable
	val, ok := c.Get("FOO")
	if !ok {
		t.Fatalf("failed to get variable in parent scope")
	}
	if val.(string) != "BAR" {
		t.Fatalf("got variable; wrong value")
	}

}
