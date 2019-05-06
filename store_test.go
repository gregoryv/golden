package golden

import (
	"io/ioutil"
	"os"
	"path"
	"testing"
)

func TestAssert(t *testing.T) {
	got := doSomething()
	Assert(t, got)
}

func TestAssert_err(t *testing.T) {
	Assert(&noTest{}, "blah")
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
	store := &store{
		RootDir:   "/var/x",
		IndexFile: "",
		skip:      3,
	}
	mock := &noTest{}
	store.save(mock, []byte("hepp"))
	if mock.ok {
		t.Fail()
	}
}

func TestStore_load(t *testing.T) {
	dir, err := ioutil.TempDir("", "golden")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)
	store := &store{
		RootDir:   dir,
		IndexFile: path.Join(dir, "index.txt"),
		skip:      3,
	}
	got := store.load()
	if string(got) != "" {
		t.Fail()
	}
}

func doSomething() string { return "hello" }

type noTest struct {
	ok bool
}

func (t *noTest) Errorf(string, ...interface{}) { t.ok = false }
func (t *noTest) Helper()                       {}
func (t *noTest) Fatal(...interface{})          { t.ok = false }

var noop *noTest = &noTest{}
var t *noTest = &noTest{}

func Test_nosave(t *testing.T) {
	// Leave this test last as it sets a global
	*updateGolden = false
	SaveString(t, "hepp")
}
