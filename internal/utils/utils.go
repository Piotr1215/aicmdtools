package utils

import (
	"io"
	"log"
	"os"
	"runtime"
	"strings"
)

type FileOpener interface {
	Open(name string) (FileReaderCloser, error)
}

type osFileOpener struct{}

func (o *osFileOpener) Open(name string) (FileReaderCloser, error) {
	return os.Open(name)
}

type FileReaderCloser interface {
	io.Reader
	io.Closer
}

type FileReader struct {
	FilePathFunc func() string
	FileOpener   FileOpener
}

func (fr *FileReader) ReadFile() string {
	filePath := fr.FilePathFunc()

	if fr.FileOpener == nil {
		fr.FileOpener = &osFileOpener{}
	}

	file, err := fr.FileOpener.Open(filePath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	content, err := io.ReadAll(file)
	if err != nil {
		log.Fatal(err)
	}
	return string(content)
}

func DetectOSAndShell() (string, string) {
	os := runtime.GOOS
	var shell string
	switch os {
	case "windows":
		shell = "cmd"
	default:
		shell = "bash"
	}
	return os, shell
}

func ReplacePlaceholders(prompt, os, shell string) string {
	prompt = strings.ReplaceAll(prompt, "{os}", os)
	prompt = strings.ReplaceAll(prompt, "{shell}", shell)
	return prompt
}
