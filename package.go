/*
Package golden enables reading and writing golden files in testdata.

Within your test, check complex output

   func TestMeShort(t *testing.T) {
       complex := doSomething()
       // Does got equal the content of the golden file and update
       // golden file if -update-golden flag is given.
       golden.Assert(t, complex)
   }

Golden file is saved in testdata/package.TestMeShort and an entry is added to
testdata/golden.files

To update the golden files use

    go test -args -update-golden

*/
package golden

import (
	"flag"
	"io/ioutil"
	"strings"
)

var (
	defaultStore *Store = NewStore()

	updateGolden = flag.Bool("update-golden", false, "Update golden files")
)

// AssertWith compares got with the contents of filename.  If
// -update-golden flag is given got is saved into filename.
func AssertWith(t T, got, filename string) {
	t.Helper()
	if *updateGolden {
		ioutil.WriteFile(filename, []byte(got), 0644)
	}
	body, _ := ioutil.ReadFile(filename)
	exp := string(body)
	if got != exp {
		t.Errorf("Got ----\n%s\nexpected ----\n%s\n", got, exp)
	}
}

// T defines parts of testing.T needed in this package
type T interface {
	Errorf(string, ...interface{})
	Helper()
	Fatal(...interface{})
}

// Assert compares got to the contents of the default golden file
// found in testdata/ matching the name of the test calling the
// assert.
func Assert(t T, got string) {
	t.Helper()
	defaultStore.Save(t, []byte(got))
	exp := string(defaultStore.Load())
	if got != exp {
		t.Errorf("Got ----\n%s\nexpected ----\n%s\n", got, exp)
	}
}

// Load returns the content of a stored golden file, defaults to empty slice.
func Load() []byte {
	return defaultStore.Load()
}

// Save saves the data as a golden file using the callers func name
func Save(t T, data []byte) {
	t.Helper()
	defaultStore.Save(t, data)
}

// LoadString loads the golden string from file using the default store
func LoadString() string {
	return string(defaultStore.Load())
}

func SaveString(t T, data string) {
	t.Helper()
	defaultStore.Save(t, []byte(data))
}

func cleanFilename(filename string) string {
	return strings.Map(
		func(r rune) rune {
			switch r {
			case '(', '*', ')':
				return -1
			}
			return r
		},
		filename,
	)
}