/* package golden enables reading and writing golden files in testdata

   func TestMeShort(t *testing.T) {
       got := doSomething()
       // Does got equal the content of the golden file and update
       // golden file if -update-golden flag is given.
       golden.Assert(t, got)
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
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
)

var defaultStore *store

func init() {
	defaultStore = newStore()
}

// AssertWith compares got with the contents of filename.  If
// -update-golden flag is given got is saved into filename.
func AssertWith(t T, got, filename string) {
	t.Helper()
	body, _ := ioutil.ReadFile(filename)
	exp := string(body)
	if got != exp {
		t.Errorf("Got ----\n%s\nexpected ----\n%s\n", got, exp)
	}
	if !*updateGolden {
		return
	}
	ioutil.WriteFile(filename, []byte(got), 0644)
}

// Assert compares got to the contents of the default golden file
// found in testdata/ matching the name of the test calling the
// assert.
func Assert(t T, got string) {
	t.Helper()
	exp := string(defaultStore.load())
	if got != exp {
		t.Errorf("Got ----\n%s\nexpected ----\n%s\n", got, exp)
	}
	defaultStore.save(t, []byte(got))
}

// Load returns the content of a stored golden file, defaults to empty slice.
func Load() []byte {
	return defaultStore.load()
}

// Save saves the data as a golden file using the callers func name
func Save(t T, data []byte) {
	t.Helper()
	defaultStore.save(t, data)
}

// LoadString loads the golden string from file using the default store
func LoadString() string {
	return string(defaultStore.load())
}

func SaveString(t T, data string) {
	t.Helper()
	defaultStore.save(t, []byte(data))
}

type store struct {
	RootDir   string
	IndexFile string
	skip      int
}

// newStore returns a Store initialized with testdata as RootDir and
// golden.files as IndexFile
func newStore() *store {
	return &store{
		RootDir:   "testdata",
		IndexFile: filepath.Join("testdata", "golden.files"),
		skip:      3,
	}
}

func (s *store) save(t T, data []byte) {
	if !*updateGolden {
		return
	}
	t.Helper()
	once.Do(s.resetGoldenFiles(t))
	filename, file := s.filenameFromCaller(s.skip)
	err := ioutil.WriteFile(file, data, 0644)
	if err != nil {
		t.Fatal(err)
	}
	// Append the output to list of golden files so it's easy to spot
	// when a file should be removed.
	flag := os.O_APPEND | os.O_CREATE | os.O_WRONLY
	f, err := os.OpenFile(s.IndexFile, flag, 0644)
	if err != nil {
		t.Fatal(err)
		return
	}
	f.Write([]byte(filename + "\n"))
	f.Close()
}

var updateGolden = flag.Bool("update-golden", false, "Update golden files")

var once sync.Once

func (s *store) resetGoldenFiles(t T) func() {
	t.Helper()
	return func() {
		t.Helper()
		// No error checking here
		os.MkdirAll(s.RootDir, 0755)
		os.RemoveAll(s.IndexFile)
	}
}

func (s *store) load() []byte {
	_, file := s.filenameFromCaller(s.skip)
	body, _ := ioutil.ReadFile(file)
	return body
}

func (s *store) filenameFromCaller(skip int) (filename, file string) {
	pc, _, _, _ := runtime.Caller(skip)
	fullName := runtime.FuncForPC(pc).Name()
	fullName = cleanFilename(fullName)
	filename = filepath.Base(fullName)
	file = path.Join(s.RootDir, filename)
	return
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

// T defines parts of testing.T needed in this package
type T interface {
	Errorf(string, ...interface{})
	Helper()
	Fatal(...interface{})
}
