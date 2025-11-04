package config

import (
	"log"
	"os"
	"os/user"
	"path/filepath"

	"github.com/piotr1215/aicmdtools/internal/utils"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Provider         string  `yaml:"provider"` // "openai" or "anthropic"
	Model            string  `yaml:"model"`
	Temperature      float64 `yaml:"temperature"`
	MaxTokens        int     `yaml:"max_tokens"`
	Safety           bool    `yaml:"safety"`
	OpenAI_APIKey    string  `yaml:"openai_api_key"`
	Anthropic_APIKey string  `yaml:"anthropic_api_key"`
}

func ReadAndParseConfig(configFilename, promptFilename string) (*Config, string, error) {
	configReader := &utils.FileReader{
		FilePathFunc: func() string { return ConfigFilePath(configFilename) },
	}
	configContent := configReader.ReadFile()
	conf := ParseConfig(configContent)

	promptReader := &utils.FileReader{
		FilePathFunc: func() string { return ConfigFilePath(promptFilename) },
	}
	prompt := promptReader.ReadFile()

	return &conf, prompt, nil
}
func ConfigFilePath(filename string) string {
	homeDir := os.Getenv("HOME")
	if homeDir == "" {
		usr, err := user.Current()
		if err != nil {
			log.Fatal(err)
		}
		homeDir = usr.HomeDir
	}

	configDir := filepath.Join(homeDir, ".config", "aicmdtools")
	return filepath.Join(configDir, filename)
}

func ParseConfig(configContent string) Config {
	var config Config
	err := yaml.Unmarshal([]byte(configContent), &config)
	if err != nil {
		log.Fatal(err)
	}
	return config
}
