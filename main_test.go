package goai

import (
	"io/ioutil"
	"log"
	"os"
	"testing"

	"gopkg.in/yaml.v3"
)

type TestConfigReader struct {
	configFilePath string
}

func createTempConfigFile(t *testing.T, content string) string {
	tmpFile, err := ioutil.TempFile("", "test_config.yaml")
	if err != nil {
		t.Fatal(err)
	}

	_, err = tmpFile.WriteString(content)
	if err != nil {
		t.Fatal(err)
	}
	tmpFile.Close()

	return tmpFile.Name()
}
func (tcr *TestConfigReader) ReadConfig() Config {
	file, err := os.Open(tcr.configFilePath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	var config Config
	err = yaml.NewDecoder(file).Decode(&config)
	if err != nil {
		log.Fatal(err)
	}
	return config
}

func TestReadPrompt(t *testing.T) {
	testPromptContent := "Test prompt content\n"

	tempPromptFilePath := createTempConfigFile(t, testPromptContent)
	defer os.Remove(tempPromptFilePath)

	promptReader := &FileReader{
		FilePathFunc: func() string { return tempPromptFilePath },
	}

	result := promptReader.ReadFile()
	if result != testPromptContent {
		t.Errorf("ReadFile() returned %s, expected %s", result, testPromptContent)
	}
}

func TestReadConfig(t *testing.T) {
	testConfigContent := "model: text-davinci-002\n" +
		"temperature: 0.5\n" +
		"max_tokens: 100\n" +
		"safety: true\n" +
		"openai_api_key: test_api_key\n"

	tempConfigFilePath := createTempConfigFile(t, testConfigContent)
	defer os.Remove(tempConfigFilePath)

	configReader := &FileReader{
		FilePathFunc: func() string { return tempConfigFilePath },
	}

	result := configReader.ReadFile()
	if result != testConfigContent {
		t.Errorf("ReadFile() returned %s, expected %s", result, testConfigContent)
	}
}
