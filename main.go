package goai

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/joho/godotenv"
	"github.com/sashabaranov/go-openai"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Model         string  `yaml:"model"`
	Temperature   float64 `yaml:"temperature"`
	MaxTokens     int     `yaml:"max_tokens"`
	Safety        bool    `yaml:"safety"`
	OpenAI_APIKey string  `yaml:"openai_api_key"`
}

func configFilePath(filename string) string {
	homeDir := os.Getenv("HOME")
	if homeDir == "" {
		usr, err := user.Current()
		if err != nil {
			log.Fatal(err)
		}
		homeDir = usr.HomeDir
	}

	configDir := filepath.Join(homeDir, ".config", "goai")
	return filepath.Join(configDir, filename)
}

type FileReader struct {
	filePathFunc func() string
}

func (fr *FileReader) ReadFile() string {
	filePath := fr.filePathFunc()

	file, err := os.Open(filePath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	content, err := io.ReadAll(file)
	if err != nil {
		log.Fatal(err)
	}
	return string(content)
}

func detectOSAndShell() (string, string) {
	os := runtime.GOOS
	var shell string
	switch os {
	case "windows":
		shell = "cmd"
	default:
		shell = "bash"
	}
	return os, shell
}

func replacePlaceholders(prompt, os, shell string) string {
	prompt = strings.ReplaceAll(prompt, "{os}", os)
	prompt = strings.ReplaceAll(prompt, "{shell}", shell)
	return prompt
}

func parseConfig(configContent string) Config {
	var config Config
	err := yaml.Unmarshal([]byte(configContent), &config)
	if err != nil {
		log.Fatal(err)
	}
	return config
}

func createOpenAIClient(config Config) *openai.Client {
	_ = godotenv.Load()

	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		apiKey = config.OpenAI_APIKey
	}

	return openai.NewClient(apiKey)
}

func main() {
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
