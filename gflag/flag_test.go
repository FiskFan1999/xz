package gflag

import (
	"bytes"
	"testing"
)

func TestFlagSet_Bool(t *testing.T) {
	f := NewFlagSet("Bool", ContinueOnError)
	a := f.Bool("test-a", false, "")
	b := f.BoolP("test-b", "b", true, "")

	err := f.Parse([]string{"--test-a", "-b", "false"})
	if err != nil {
		t.Fatalf("f.Parse error %s", err)
	}

	if *a != true {
		t.Errorf("*a is %t; want %t", *a, true)
	}
	if *b != false {
		t.Errorf("*b is %t; want %t", *b, false)
	}

	if f.NArg() != 0 {
		t.Errorf("f.NArg() is %d; want %d", f.NArg(), 0)
	}
}

func TestFlagSet_Counter(t *testing.T) {
	f := NewFlagSet("Bool", ContinueOnError)
	a := f.Counter("test-a", 0, "")
	b := f.CounterP("test-b", "b", 0, "")
	err := f.Parse([]string{"--test-a=3", "-b", "5", "--test-a", "-b"})
	if err != nil {
		t.Fatalf("f.Parse error %s", err)
	}

	if *a != 4 {
		t.Errorf("*a is %d; want %d", *a, 4)
	}
	if *b != 6 {
		t.Errorf("*b is %d; want %d", *b, 6)
	}

	if f.NArg() != 0 {
		t.Errorf("f.NArg() is %d; want %d", f.NArg(), 0)
	}
}

func TestFlagSet_Int(t *testing.T) {
	f := NewFlagSet("Bool", ContinueOnError)
	a := f.Int("test-a", 0, "")
	b := f.IntP("test-b", "b", 0, "")
	c := f.Int("c", 0, "")
	err := f.Parse([]string{"--test-a=0x23", "foo", "-b", "077",
		"-c", "33", "bar"})
	if err != nil {
		t.Fatalf("f.Parse error %s", err)
	}

	if *a != 0x23 {
		t.Errorf("*a is %d; want %d", *a, 0x23)
	}
	if *b != 077 {
		t.Errorf("*b is %d; want %d", *b, 077)
	}
	if *c != 33 {
		t.Errorf("*c is %d; want %d", *c, 33)
	}

	if f.NArg() != 2 {
		t.Errorf("f.NArg() is %d; want %d", f.NArg(), 2)
	}

	for i, s := range []string{"foo", "bar"} {
		if f.Arg(i) != s {
			t.Errorf("f.Arg(%d) is %s; want %s", i, f.Arg(i), s)
		}
	}
}

func TestFlagSet_Usage(t *testing.T) {
	f := NewFlagSet("test", ContinueOnError)
	f.IntP("test-a", "a", 3, "tests a")
	f.CounterP("count-b", "b", 0, "counts b")
	buf := new(bytes.Buffer)
	f.SetOutput(buf)
	f.usage()
	t.Log(buf.String())
}