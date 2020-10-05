package golden

import (
	"io/ioutil"
	"os"
	"path"
	"testing"
)

func TestAssert_BDD_helpers(t *testing.T) {
	bdd := &bdd{}
	bdd.method_helper(t)
}

type bdd struct{}

func (b *bdd) method_helper(t *testing.T) {
	t.Helper()
	Assert(noop, "blah")
	// only tested when flags -args -update-golden is given
	// the test is for cleanFilename func
	if err := os.RemoveAll("testdata/golden.bdd.method_helper"); err != nil {
		t.Error(err)
	}
}

func TestAssertWith(t *testing.T) {
	got := doSomething()
	fh, err := ioutil.TempFile("", "golden")
	if err != nil {
		t.Fatal(err)
	}
	fh.Close()
	defer os.RemoveAll(fh.Name())
	*updateGolden = false
	mock := &noTest{ok: true}
	AssertWith(mock, got, fh.Name())
	if mock.ok {
		t.Error("Assert should fail")
	}
	*updateGolden = true
	mock = &noTest{ok: true}
	AssertWith(mock, got, fh.Name())
	if !mock.ok {
		t.Error("Assert should not fail when updating")
	}
}

func TestAssert(t *testing.T) {
	got := doSomething()
	mock := &noTest{ok: true}
	Assert(mock, got)
	if !mock.ok {
		t.Error("Assert should be ok")
	}
}

func TestAssert_err(t *testing.T) {
	mock := &noTest{ok: true}
	*updateGolden = false // global, other tests can affect it
	Assert(mock, "blah")
	if mock.ok {
		t.Error("Assert should have failed")
	}
	os.RemoveAll("testdata/golden.TestAssert_err")
}

func TestLoad(t *testing.T) {
	got := doSomething()
	exp := Load()
	if got != string(exp) {
		t.Errorf("Got %q, expected %q", got, exp)
	}
	Save(t, []byte(got))
}

func TestLoadString(t *testing.T) {
	got := doSomething()
	exp := LoadString()
	if got != exp {
		t.Errorf("Got %q, expected %q", got, exp)
	}
	SaveString(t, got)
}

func Test_fail(t *testing.T) {
	Store := &Store{
		RootDir:   "/var/x",
		IndexFile: "",
		skip:      3,
	}
	mock := &noTest{ok: true}
	*updateGolden = true
	Store.Save(mock, []byte("hepp"))
	if mock.ok {
		t.Fail()
	}
}

func TestStore_Load(t *testing.T) {
	dir, err := ioutil.TempDir("", "golden")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)
	Store := &Store{
		RootDir:   dir,
		IndexFile: path.Join(dir, "index.txt"),
		skip:      3,
	}
	got := Store.Load()
	if string(got) != "" {
		t.Fail()
	}
}

func doSomething() string { return "hello" }

type noTest struct {
	ok     bool
	format string
	v      []interface{}
}

func (t *noTest) Helper() {}
func (t *noTest) Errorf(f string, v ...interface{}) {
	t.ok = false
	t.format = f
	t.v = v
}
func (t *noTest) Fatal(...interface{}) { t.ok = false }

var noop *noTest = &noTest{}
var t *noTest = &noTest{}

func Test_nosave(t *testing.T) {
	*updateGolden = false
	SaveString(t, "hepp")
	*updateGolden = true
}
