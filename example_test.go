package golden_test

import (
	"io/ioutil"

	"github.com/gregoryv/golden"
)

var somefile string

func init() {
	fh, _ := ioutil.TempFile("", "golden")
	fh.Close()
	somefile = fh.Name()
}

func ExampleAssertEquals() {
	got := doSomething()
	exp := "hello"
	golden.AssertEquals(t, got, exp)
}

func ExampleAssertWith() {
	got := doSomething()
	golden.AssertWith(t, got, somefile)
}

func ExampleAssert() {
	got := doSomething()
	// Assert and update if -update-golden flag is given
	golden.Assert(t, got)
}

func ExampleSaveString() {
	got := doSomething()
	exp := golden.LoadString()
	if got != exp {
		t.Errorf("Got %q, expected %q", got, exp)
	}
	// Save if -update-golden flag is given
	golden.SaveString(t, got)
}

func doSomething() string { return "hello" }

type noTest struct {
	ok bool
}

func (t *noTest) Errorf(string, ...interface{}) { t.ok = false }
func (t *noTest) Helper()                       {}
func (t *noTest) Fatal(...interface{})          { t.ok = false }

var t *noTest = &noTest{}
