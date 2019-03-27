/* package golden enables reading and writing golden files in testdata

   func TestMe(t *testing.T) {
       got := doSomething()
       exp := golden.LoadString()
       if got != exp {
           t.Fail()
       }
       golden.SaveString(t, got)
   }

Golden file is saved in testdata/package.TestMe and an entry is added to
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
	"sync"
)

var defaultStore *store

func init() {
	defaultStore = newStore()
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
	filename = filepath.Base(fullName)
	file = path.Join(s.RootDir, filename)
	return
}

// T defines parts of testing.T needed in this package
type T interface {
	Helper()
	Fatal(...interface{})
}
