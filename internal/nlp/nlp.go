package nlp

import (
	"context"
	"fmt"
	"os"

	"github.com/anthropics/anthropic-sdk-go"
	anthropicoption "github.com/anthropics/anthropic-sdk-go/option"
	"github.com/joho/godotenv"
	"github.com/piotr1215/aicmdtools/internal/config"
	"github.com/sashabaranov/go-openai"
)

type GAIClient interface {
	ProcessCommand(userPrompt string, conf config.Config) (*openai.ChatCompletionResponse, error)
	ProcessCommandWithContext(ctx context.Context, userPrompt string, conf config.Config) (*openai.ChatCompletionResponse, error)
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

func (g *GoaiClient) ProcessCommandWithContext(ctx context.Context, userPrompt string, conf config.Config) (*openai.ChatCompletionResponse, error) {
	// Map the model from the config to the OpenAI model
	response, err := g.Client.CreateChatCompletion(
		ctx,
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

// AnthropicClient wraps the Anthropic SDK client
type AnthropicClient struct {
	Client *anthropic.Client
	Prompt string
}

func (a *AnthropicClient) ProcessCommand(userPrompt string, conf config.Config) (*openai.ChatCompletionResponse, error) {
	return a.ProcessCommandWithContext(context.Background(), userPrompt, conf)
}

func (a *AnthropicClient) ProcessCommandWithContext(ctx context.Context, userPrompt string, conf config.Config) (*openai.ChatCompletionResponse, error) {
	// Create Anthropic message request with system prompt
	system := []anthropic.TextBlockParam{
		{Text: a.Prompt},
	}

	messages := []anthropic.MessageParam{
		anthropic.NewUserMessage(anthropic.NewTextBlock(userPrompt)),
	}

	message, err := a.Client.Messages.New(ctx, anthropic.MessageNewParams{
		Model:     anthropic.Model(conf.Model),
		MaxTokens: int64(conf.MaxTokens),
		System:    system,
		Messages:  messages,
	})

	if err != nil {
		return nil, fmt.Errorf("Anthropic API error: %v", err)
	}

	// Convert Anthropic response to OpenAI format for compatibility
	content := ""
	for _, block := range message.Content {
		// ContentBlockUnion is a struct, access Text field directly
		content += block.Text
	}

	response := &openai.ChatCompletionResponse{
		ID:    message.ID,
		Model: string(message.Model),
		Choices: []openai.ChatCompletionChoice{
			{
				Message: openai.ChatCompletionMessage{
					Role:    "assistant",
					Content: content,
				},
			},
		},
	}

	return response, nil
}

func CreateAnthropicClient(conf config.Config) *anthropic.Client {
	_ = godotenv.Load()

	apiKey := os.Getenv("ANTHROPIC_API_KEY")
	if apiKey == "" {
		apiKey = conf.Anthropic_APIKey
	}

	client := anthropic.NewClient(anthropicoption.WithAPIKey(apiKey))

	return &client
}
