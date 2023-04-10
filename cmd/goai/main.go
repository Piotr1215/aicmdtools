package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/piotr1215/goai"
)

func shouldExecuteCommand(config *goai.Config, reader io.Reader) bool {
	if !config.Safety {
		return true
	}

	fmt.Print("Execute the command? [Enter/n] ==> ")
	var answer string
	_, _ = fmt.Fscanln(reader, &answer)

	return strings.ToUpper(answer) != "N"
}

func main() {
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

	client := goai.CreateOpenAIClient(config)

	aiClient := &goai.GoaiClient{
		Client: client,
		Prompt: prompt,
	}

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
		execute = shouldExecuteCommand(&config, os.Stdin)
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
