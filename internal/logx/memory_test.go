package logx

import "testing"

func TestMemory_Name(t *testing.T) {
	mem := NewMemory("log")
	if mem.Name != "log" {
		t.Errorf("Name: unexpected name %q", mem.Name)
	}
}

func TestMemory_Write(t *testing.T) {
	mem := NewMemory("log")
	if len(mem.Lines) != 0 {
		t.Fatalf("Write: unexpected line count %v", len(mem.Lines))
	}

	n, err := mem.Write([]byte("hello world"))
	if err != nil {
		t.Fatalf("Write: unexpected error %v", err)
	}
	if n != 11 {
		t.Errorf("Write: unexpected byte count %v", n)
	}

	if len(mem.Lines) != 1 {
		t.Fatalf("Write: unexpected line count %v", len(mem.Lines))
	}
	if mem.Lines[0] != "hello world" {
		t.Errorf("Write: unexpected line #0 %q", mem.Lines[0])
	}
}

func TestMemory_Has(t *testing.T) {
	mem := NewMemory("log")
	if mem.Has("hello world") {
		t.Error("Has: unexpected true")
	}
	_, _ = mem.Write([]byte("hello world"))
	if !mem.Has("hello world") {
		t.Error("Has: unexpected false")
	}
	_, _ = mem.Write([]byte("one two three four"))
	if !mem.Has("one two") {
		t.Error("Has: one two: unexpected false")
	}
	if !mem.Has("two three") {
		t.Error("Has: two three: unexpected false")
	}
	if mem.Has("one three") {
		t.Error("Has: one three: unexpected true")
	}
	if !mem.Has("one", "three") {
		t.Error("Has: [one three]: unexpected false")
	}
	if !mem.Has("one", "three", "four") {
		t.Error("Has: [one three four]: unexpected false")
	}
	if !mem.Has("four", "three") {
		t.Error("Has: [four three]: unexpected false")
	}
}
