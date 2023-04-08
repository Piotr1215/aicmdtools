package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/piotr1215/goai"
)

func main() {
	// Get the configuration and prompt
	configReader := &goai.FileReader{
		FilePathFunc: func() string { return goai.ConfigFilePath("config.yaml") },
	}
	configContent := configReader.ReadFile()
	config := goai.ParseConfig(configContent)

	promptReader := &goai.FileReader{
		FilePathFunc: func() string { return goai.ConfigFilePath("prompt.txt") },
	}
	prompt := promptReader.ReadFile()
	operating_system, shell := goai.DetectOSAndShell()
	prompt = goai.ReplacePlaceholders(prompt, operating_system, shell)

	aiClient := goai.CreateGoAIClient() // Remove the argument

	if len(os.Args) < 2 {
		fmt.Println("No user prompt specified.")
		os.Exit(-1)
	}

	userPrompt := strings.Join(os.Args[1:], " ")

	response, err := aiClient.ProcessCommand(userPrompt)
	if err != nil {
		fmt.Printf("Error processing command: %v\n", err)
		return
	}

	command := response.Choices[0].Message.Content
	fmt.Printf("Command: %s\n", command)

	execute := true
	if config.Safety {
		fmt.Print("Execute the command? [Y/n] ==> ")
		var answer string
		_, _ = fmt.Scanln(&answer)
		if strings.ToUpper(answer) == "N" {
			execute = false
		}
	}

	if execute {
		var cmd *exec.Cmd
		// Use "sh -c" for Unix-like systems and "cmd /C" for Windows
		if runtime.GOOS == "windows" {
			cmd = exec.Command("cmd", "/C", command)
		} else {
			cmd = exec.Command("sh", "-c", command)
		}
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err = cmd.Run()
		if err != nil {
			log.Fatal(err)
		}
	}
}
