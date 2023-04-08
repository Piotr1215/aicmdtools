package goai

import (
	"io"
	"log"
	"os"
	"runtime"
	"strings"
)

type FileReader struct {
	filePathFunc func() string
}

func (fr *FileReader) ReadFile() string {
	filePath := fr.filePathFunc()

	file, err := os.Open(filePath)
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

func detectOSAndShell() (string, string) {
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

func replacePlaceholders(prompt, os, shell string) string {
	prompt = strings.ReplaceAll(prompt, "{os}", os)
	prompt = strings.ReplaceAll(prompt, "{shell}", shell)
	return prompt
}
