package nlp

import (
	"context"
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/piotr1215/aicmdtools/internal/config"
	"github.com/sashabaranov/go-openai"
)

type GAIClient interface {
	ProcessCommand(userPrompt string) (*openai.ChatCompletionResponse, error)
}

type GoaiClient struct {
	Client *openai.Client
	Prompt string
}

func (g *GoaiClient) ProcessCommand(userPrompt string) (*openai.ChatCompletionResponse, error) {

	response, err := g.Client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT4,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleSystem,
					Content: g.Prompt,
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
func CreateOpenAIClient(config config.Config) *openai.Client {
	_ = godotenv.Load()

	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		apiKey = config.OpenAI_APIKey
	}

	return openai.NewClient(apiKey)
}
