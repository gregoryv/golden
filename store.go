package golden

import (
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"sync"
)

// Store defines a location and index of golden files.
type Store struct {
	RootDir   string
	IndexFile string
	skip      int

	resetOnce sync.Once
}

// NewStore returns a Store initialized with testdata as RootDir and
// golden.files as IndexFile
func NewStore() *Store {
	return &Store{
		RootDir:   "testdata",
		IndexFile: filepath.Join("testdata", "golden.files"),
		skip:      3,
	}
}

// Load returns body of a golden file based on caller func name.
// Empty if no golden file exists.
func (s *Store) Load() []byte {
	_, file := s.filenameFromCaller(s.skip)
	body, _ := ioutil.ReadFile(file)
	return body
}

func (s *Store) filenameFromCaller(skip int) (filename, file string) {
	pc, _, _, _ := runtime.Caller(skip)
	fullName := runtime.FuncForPC(pc).Name()
	fullName = cleanFilename(fullName)
	filename = filepath.Base(fullName)
	file = path.Join(s.RootDir, filename)
	return
}

// Save writes the given data to RootDir with name from caller func
func (s *Store) Save(t T, data []byte) {
	if !*updateGolden {
		return
	}
	t.Helper()
	s.resetOnce.Do(s.resetGoldenFiles(t))
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

func (s *Store) resetGoldenFiles(t T) func() {
	t.Helper()
	return func() {
		t.Helper()
		// No error checking here
		os.MkdirAll(s.RootDir, 0755)
		os.RemoveAll(s.IndexFile)
	}
}
