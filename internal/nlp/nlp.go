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
	ProcessCommand(userPrompt string, conf config.Config) (*openai.ChatCompletionResponse, error)
}

type GoaiClient struct {
	Client *openai.Client
	Prompt string
}

func (g *GoaiClient) ProcessCommand(userPrompt string, conf config.Config) (*openai.ChatCompletionResponse, error) {

	// Map the model from the config to the OpenAI model
	response, err := g.Client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: conf.Model,
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
func CreateOpenAIClient(conf config.Config) *openai.Client {
	_ = godotenv.Load()

	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		apiKey = conf.OpenAI_APIKey
	}

	client := openai.NewClient(apiKey)

	return client
}
