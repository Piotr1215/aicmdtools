package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/piotr1215/aicmdtools/internal/aicmd"
	"github.com/piotr1215/aicmdtools/internal/config"
	"github.com/piotr1215/aicmdtools/internal/utils"
)

var version = "v0.0.9"

func shouldExecuteCommand(config *config.Config, reader io.Reader) bool {
	if !config.Safety {
		return true
	}

	fmt.Print("Execute the command? [Enter/n] ==> ")
	var answer string
	_, _ = fmt.Fscanln(reader, &answer)
	fmt.Printf("User input: %q\n", answer) // Add this line to print the user input

	return strings.ToUpper(answer) != "N"
}

func main() {
	versionFlag := flag.Bool("version", false, "Display version information")
	flag.Parse()

	if *versionFlag {
		fmt.Printf("Goai version: %s\n", version)
		return
	}

	configReader := &utils.FileReader{
		FilePathFunc: func() string { return config.ConfigFilePath("config.yaml") },
	}
	configContent := configReader.ReadFile()
	conf := config.ParseConfig(configContent)

	promptReader := &utils.FileReader{
		FilePathFunc: func() string { return config.ConfigFilePath("prompt.txt") },
	}
	prompt := promptReader.ReadFile()
	operating_system, shell := utils.DetectOSAndShell()
	prompt = utils.ReplacePlaceholders(prompt, operating_system, shell)

	client := aicmd.CreateOpenAIClient(conf)

	aiClient := &aicmd.GoaiClient{
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
	if conf.Safety {
		execute = shouldExecuteCommand(&conf, os.Stdin)
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
