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
func CreateOpenAIClient(config Config) *openai.Client {
	_ = godotenv.Load()

	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		apiKey = config.OpenAI_APIKey
	}

	return openai.NewClient(apiKey)
}

func CreateGoAIClient() GoAIClient {
	configReader := &FileReader{
		FilePathFunc: func() string { return ConfigFilePath("config.yaml") },
	}
	configContent := configReader.ReadFile()
	config := ParseConfig(configContent)

	promptReader := &FileReader{
		FilePathFunc: func() string { return ConfigFilePath("prompt.txt") },
	}
	prompt := promptReader.ReadFile()
	operating_system, shell := DetectOSAndShell()
	prompt = ReplacePlaceholders(prompt, operating_system, shell)

	client := CreateOpenAIClient(config)

	return &goaiClient{
		client: client,
		prompt: prompt,
	}
}
