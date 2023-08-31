package aichat

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/piotr1215/aicmdtools/internal/config"
	"github.com/piotr1215/aicmdtools/internal/nlp"
	"github.com/piotr1215/aicmdtools/internal/utils"
)

var prompt_file = "chat-prompt.txt"

func Initialize() *nlp.GoaiClient {
	configReader := &utils.FileReader{
		FilePathFunc: func() string { return config.ConfigFilePath("config.yaml") },
	}
	configContent := configReader.ReadFile()
	conf := config.ParseConfig(configContent)

	promptReader := &utils.FileReader{
		FilePathFunc: func() string { return config.ConfigFilePath(prompt_file) },
	}
	prompt := promptReader.ReadFile()
	operating_system, shell := utils.DetectOSAndShell()
	prompt = utils.ReplacePlaceholders(prompt, operating_system, shell)

	client := nlp.CreateOpenAIClient(conf)

	return &nlp.GoaiClient{
		Client: client,
		Prompt: prompt,
	}
}

func SendMessage(client *nlp.GoaiClient, userMessage string) (string, error) {
	userMessage = strings.TrimSpace(userMessage)
	response, err := client.ProcessCommand(userMessage)
	if err != nil {
		return "", err
	}
	return response.Choices[0].Message.Content, nil
}

func Execute() error {
	client := Initialize()
	reader := bufio.NewReader(os.Stdin)

	for {
		// Color "You:" in green
		fmt.Print("\033[32mYou:\033[0m ")
		userMessage, _ := reader.ReadString('\n')
		userMessage = strings.TrimSpace(userMessage)

		if userMessage == "quit" {
			break
		}

		aiResponse, err := SendMessage(client, userMessage)
		if err != nil {
			// Color "Error:" in red
			fmt.Println("\033[31mError:\033[0m", err)
			continue
		}

		// Color "AI:" in blue
		fmt.Println("\033[34mAI:\033[0m", aiResponse)
	}

	return nil
}
