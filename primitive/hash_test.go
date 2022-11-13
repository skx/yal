package primitive

import "testing"

func TestHash(t *testing.T) {

	// Create a hash
	h := NewHash()

	out := h.Get("NAME")
	_, ok := out.(Nil)
	if !ok {
		t.Fatalf("expected nil getting hash value that is absent")
	}

	h.Set("NAME", String("ME"))
	valid := h.Get("NAME")
	if valid.ToString() != "ME" {
		t.Fatalf("got wrong value")
	}

	if h.Type() != "hash" {
		t.Fatalf("Wrong type for hash")
	}

	if h.ToString() != "{\n\tNAME => ME\n}" {
		t.Fatalf("string value of hash was wrong")
	}
}

func TestHashStruct(t *testing.T) {

	// Create a hash
	h := NewHash()

	// mark it as a struct
	h.SetStruct("pie")

	out := h.Get("NAME")
	_, ok := out.(Nil)
	if !ok {
		t.Fatalf("expected nil getting hash value that is absent")
	}

	h.Set("NAME", String("ME"))
	valid := h.Get("NAME")
	if valid.ToString() != "ME" {
		t.Fatalf("got wrong value")
	}

	if h.Type() != "struct-pie" {
		t.Fatalf("Wrong type for struct")
	}

	if h.GetStruct() != "pie" {
		t.Fatalf("Wrong struct-type for struct")
	}

}
