package golden

import (
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"sync"
)

type store struct {
	RootDir   string
	IndexFile string
	skip      int

	resetOnce sync.Once
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
