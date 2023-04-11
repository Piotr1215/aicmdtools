package utils

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

const changelogScript = `#!/usr/bin/env bash

# Set the default number of commits to 10
num_commits="${1:-10}"

# Get the remote URL for the current Git repository
remote_url=$(git remote get-url origin | sed -E 's/.*github.com.(.*)\.git/\1/' | sed 's/:/\//')

# Get the git log
git_log_output=$(git log --oneline --decorate=short -n "$num_commits" 2>&1)
exit_status=$?

# If git log exits with a non-zero status, exit the script with an error
if [ $exit_status -ne 0 ]; then
	echo "ERROR: ${git_log_output}"
	exit 1
fi

# Process the git log output
processed_git_log=$(echo "${git_log_output}" | awk -v remote_url="${remote_url}" '{ printf "- [https://github.com/%s/commit/%s](%s) - %s\n", remote_url, substr($1, 1, 7), substr($1, 1, 40), substr($0, index($0,$2)) }')

# Print the processed git log
echo -e "${processed_git_log}"
`

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

type CmdFactoryFunc func(name string, arg ...string) *exec.Cmd

func GenerateChangelog(cmdFactory CmdFactoryFunc) (string, error) {
	// Create a temporary file to store the shell script
	tmpFile, err := ioutil.TempFile("", "changelog-*.sh")
	if err != nil {
		return "", fmt.Errorf("error creating temporary file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	// Write the shell script to the temporary file
	if _, err := tmpFile.WriteString(changelogScript); err != nil {
		return "", fmt.Errorf("error writing shell script to temporary file: %v", err)
	}

	if err := tmpFile.Close(); err != nil {
		return "", fmt.Errorf("error closing temporary file: %v", err)
	}

	// Make the temporary file executable
	if err := os.Chmod(tmpFile.Name(), 0755); err != nil {
		return "", fmt.Errorf("error setting temporary file permissions: %v", err)
	}

	// Run the shell script
	cmd := cmdFactory(tmpFile.Name())
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err = cmd.Run()
	if err != nil {
		return "", fmt.Errorf("changelog generation error: %v, stderr: %s", err, stderr.String())
	}

	return stdout.String(), nil
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
