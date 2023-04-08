package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/piotr1215/goai"
	"github.com/sashabaranov/go-openai"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("No user prompt specified.")
		os.Exit(-1)
	}

	userPrompt := strings.Join(os.Args[1:], " ")

	// Initialize the GoAI client
	client := goai.CreateGoAIClient()

	// Use the GoAI client to get a response
	response, err := client.ProcessCommand(userPrompt)
	if err != nil {
		fmt.Printf("Error processing command: %v\n", err)
		return
	}

	configReader := &FileReader{
		filePathFunc: func() string { return configFilePath("config.yaml") },
	}
	configContent := configReader.ReadFile()
	config := parseConfig(configContent)

	promptReader := &FileReader{
		filePathFunc: func() string { return configFilePath("prompt.txt") },
	}
	prompt := promptReader.ReadFile()
	operating_system, shell := detectOSAndShell()
	prompt = replacePlaceholders(prompt, operating_system, shell)

	client := createOpenAIClient(config)

	if len(os.Args) < 2 {
		fmt.Println("No user prompt specified.")
		os.Exit(-1)
	}

	userPrompt := strings.Join(os.Args[1:], " ")

	response, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT4,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleSystem,
					Content: prompt,
				},
				{
					Role:    openai.ChatMessageRoleUser,
					Content: userPrompt,
				},
			},
		},
	)

	if err != nil {
		fmt.Printf("ChatCompletion error: %v\n", err)
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
