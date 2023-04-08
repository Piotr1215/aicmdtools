package goai

import (
	"context"
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/sashabaranov/go-openai"
)

type GoAIClient interface {
	ProcessCommand(userPrompt string) (*openai.ChatCompletionResponse, error)
}

type goaiClient struct {
	client *openai.Client
	prompt string
}

func (g *goaiClient) ProcessCommand(userPrompt string) (*openai.ChatCompletionResponse, error) {

	response, err := g.client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT4,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleSystem,
					Content: g.prompt,
				},
				{
					Role:    openai.ChatMessageRoleUser,
					Content: userPrompt,
				},
			},
		},
	)

	if err != nil {
		return nil, fmt.Errorf("ChatCompletion error: %v", err)
	}

	return &response, nil
}
func createOpenAIClient(config Config) *openai.Client {
	_ = godotenv.Load()

	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		apiKey = config.OpenAI_APIKey
	}

	return openai.NewClient(apiKey)
}

func CreateGoAIClient() GoAIClient {
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

	return &goaiClient{
		client: client,
		prompt: prompt,
	}
}
