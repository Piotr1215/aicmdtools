package utils

import (
	"io/ioutil"
	"strings"
	"testing"
)

type stubFileOpener struct{}

func (s *stubFileOpener) Open(name string) (FileReaderCloser, error) {
	content := "This is a test file."
	return ioutil.NopCloser(strings.NewReader(content)), nil
}

func TestReadFile(t *testing.T) {
	fileReader := &FileReader{
		FilePathFunc: func() string { return "testfile.txt" },
		FileOpener:   &stubFileOpener{},
	}

	content := fileReader.ReadFile()
	expectedContent := "This is a test file."

	if content != expectedContent {
		t.Errorf("Expected content: %s, got: %s", expectedContent, content)
	}
}
