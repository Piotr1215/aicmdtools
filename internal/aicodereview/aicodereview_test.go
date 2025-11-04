package aicodereview

import (
	"bytes"
	"context"
	"errors"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/piotr1215/aicmdtools/internal/config"
	"github.com/sashabaranov/go-openai"
	"github.com/stretchr/testify/assert"
)

// MockFileReader implements FileReader interface for testing
type MockFileReader struct {
	Content  string
	Size     int64
	ReadErr  error
	StatErr  error
	IsDir    bool
	FileName string
}

func (m *MockFileReader) ReadFile(path string) ([]byte, error) {
	if m.ReadErr != nil {
		return nil, m.ReadErr
	}
	return []byte(m.Content), nil
}

func (m *MockFileReader) Stat(path string) (os.FileInfo, error) {
	if m.StatErr != nil {
		return nil, m.StatErr
	}
	return &MockFileInfo{
		name:  m.FileName,
		size:  m.Size,
		isDir: m.IsDir,
	}, nil
}

// MockFileInfo implements os.FileInfo interface
type MockFileInfo struct {
	name  string
	size  int64
	isDir bool
}

func (m *MockFileInfo) Name() string       { return m.name }
func (m *MockFileInfo) Size() int64        { return m.size }
func (m *MockFileInfo) Mode() os.FileMode  { return 0644 }
func (m *MockFileInfo) ModTime() time.Time { return time.Now() }
func (m *MockFileInfo) IsDir() bool        { return m.isDir }
func (m *MockFileInfo) Sys() interface{}   { return nil }

// MockAIClient implements GAIClient interface for testing
type MockAIClient struct {
	Response *openai.ChatCompletionResponse
	Err      error
}

func (m *MockAIClient) ProcessCommand(userPrompt string, conf config.Config) (*openai.ChatCompletionResponse, error) {
	return m.Response, m.Err
}

func (m *MockAIClient) ProcessCommandWithContext(ctx context.Context, userPrompt string, conf config.Config) (*openai.ChatCompletionResponse, error) {
	return m.Response, m.Err
}

func TestReview_Success(t *testing.T) {
	// Setup
	mockClient := &MockAIClient{
		Response: &openai.ChatCompletionResponse{
			Choices: []openai.ChatCompletionChoice{
				{
					Message: openai.ChatCompletionMessage{
						Content: "=== CODE REVIEW RESULTS ===\n\nCRITICAL ISSUES:\nNo issues found.\n\nHIGH PRIORITY:\nNo issues found.\n\nMEDIUM PRIORITY:\nNo issues found.\n\nLOW PRIORITY:\nNo issues found.\n\nSUMMARY:\nCode looks good!",
					},
				},
			},
		},
		Err: nil,
	}

	mockFileReader := &MockFileReader{
		Content:  "package main\n\nfunc main() {\n\tprintln(\"Hello, World!\")\n}",
		Size:     100,
		IsDir:    false,
		FileName: "test.go",
	}

	reviewer := &Reviewer{
		Client:     mockClient,
		FileReader: mockFileReader,
	}

	var output bytes.Buffer
	options := ReviewOptions{
		FilePath:     "test.go",
		Format:       FormatText,
		MaxFileSize:  1024,
		OutputWriter: &output,
	}

	conf := config.Config{
		Model:       "gpt-4",
		Temperature: 0.7,
		MaxTokens:   2000,
	}

	// Execute
	err := reviewer.Review(options, conf)

	// Assert
	assert.NoError(t, err)
	assert.Contains(t, output.String(), "CODE REVIEW RESULTS")
}

func TestReview_EmptyFilePath(t *testing.T) {
	// Setup
	mockClient := &MockAIClient{}
	reviewer := &Reviewer{
		Client:     mockClient,
		FileReader: &MockFileReader{},
	}

	options := ReviewOptions{
		FilePath: "",
		Format:   FormatText,
	}

	conf := config.Config{}

	// Execute
	err := reviewer.Review(options, conf)

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "file path is required")
}

func TestReview_FileNotFound(t *testing.T) {
	// Setup
	mockClient := &MockAIClient{}
	mockFileReader := &MockFileReader{
		StatErr: errors.New("file not found"),
	}

	reviewer := &Reviewer{
		Client:     mockClient,
		FileReader: mockFileReader,
	}

	options := ReviewOptions{
		FilePath: "nonexistent.go",
		Format:   FormatText,
	}

	conf := config.Config{}

	// Execute
	err := reviewer.Review(options, conf)

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to stat file")
}

func TestReview_DirectoryNotSupported(t *testing.T) {
	// Setup
	mockClient := &MockAIClient{}
	mockFileReader := &MockFileReader{
		Size:  0,
		IsDir: true,
	}

	reviewer := &Reviewer{
		Client:     mockClient,
		FileReader: mockFileReader,
	}

	options := ReviewOptions{
		FilePath: "/some/directory",
		Format:   FormatText,
	}

	conf := config.Config{}

	// Execute
	err := reviewer.Review(options, conf)

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "directory review not yet supported")
}

func TestReview_FileTooLarge(t *testing.T) {
	// Setup
	mockClient := &MockAIClient{}
	mockFileReader := &MockFileReader{
		Size:  2 * 1024 * 1024, // 2MB
		IsDir: false,
	}

	reviewer := &Reviewer{
		Client:     mockClient,
		FileReader: mockFileReader,
	}

	options := ReviewOptions{
		FilePath:    "large.go",
		Format:      FormatText,
		MaxFileSize: 1024 * 1024, // 1MB limit
	}

	conf := config.Config{}

	// Execute
	err := reviewer.Review(options, conf)

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "exceeds maximum allowed size")
}

func TestReview_ReadFileError(t *testing.T) {
	// Setup
	mockClient := &MockAIClient{}
	mockFileReader := &MockFileReader{
		Size:    100,
		IsDir:   false,
		ReadErr: errors.New("permission denied"),
	}

	reviewer := &Reviewer{
		Client:     mockClient,
		FileReader: mockFileReader,
	}

	options := ReviewOptions{
		FilePath: "test.go",
		Format:   FormatText,
	}

	conf := config.Config{}

	// Execute
	err := reviewer.Review(options, conf)

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to read file")
}

func TestReview_AIClientError(t *testing.T) {
	// Setup
	mockClient := &MockAIClient{
		Response: nil,
		Err:      errors.New("API rate limit exceeded"),
	}

	mockFileReader := &MockFileReader{
		Content:  "package main",
		Size:     100,
		IsDir:    false,
		FileName: "test.go",
	}

	reviewer := &Reviewer{
		Client:     mockClient,
		FileReader: mockFileReader,
	}

	var output bytes.Buffer
	options := ReviewOptions{
		FilePath:     "test.go",
		Format:       FormatText,
		OutputWriter: &output,
	}

	conf := config.Config{}

	// Execute
	err := reviewer.Review(options, conf)

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to get AI review")
}

func TestReview_NoResponseChoices(t *testing.T) {
	// Setup
	mockClient := &MockAIClient{
		Response: &openai.ChatCompletionResponse{
			Choices: []openai.ChatCompletionChoice{},
		},
		Err: nil,
	}

	mockFileReader := &MockFileReader{
		Content:  "package main",
		Size:     100,
		IsDir:    false,
		FileName: "test.go",
	}

	reviewer := &Reviewer{
		Client:     mockClient,
		FileReader: mockFileReader,
	}

	var output bytes.Buffer
	options := ReviewOptions{
		FilePath:     "test.go",
		Format:       FormatText,
		OutputWriter: &output,
	}

	conf := config.Config{}

	// Execute
	err := reviewer.Review(options, conf)

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "no response from AI")
}

func TestBuildUserPrompt(t *testing.T) {
	// Setup
	reviewer := &Reviewer{}

	// Execute
	prompt := reviewer.buildUserPrompt(
		"test.go",
		"package main\n\nfunc main() {}",
		[]string{"security", "performance"},
	)

	// Assert
	assert.Contains(t, prompt, "File: test.go")
	assert.Contains(t, prompt, "Type: Go")
	assert.Contains(t, prompt, "Focus Areas: security, performance")
	assert.Contains(t, prompt, "package main")
	assert.Contains(t, prompt, "```go")
}

func TestBuildUserPrompt_NoFocusAreas(t *testing.T) {
	// Setup
	reviewer := &Reviewer{}

	// Execute
	prompt := reviewer.buildUserPrompt(
		"script.py",
		"print('hello')",
		nil,
	)

	// Assert
	assert.Contains(t, prompt, "File: script.py")
	assert.Contains(t, prompt, "Type: Python")
	assert.NotContains(t, prompt, "Focus Areas:")
	assert.Contains(t, prompt, "print('hello')")
}

func TestDetectLanguage(t *testing.T) {
	reviewer := &Reviewer{}

	tests := []struct {
		ext      string
		expected string
	}{
		{".go", "Go"},
		{".js", "JavaScript"},
		{".ts", "TypeScript"},
		{".py", "Python"},
		{".java", "Java"},
		{".rs", "Rust"},
		{".rb", "Ruby"},
		{".php", "PHP"},
		{".cpp", "C++"},
		{".cs", "C#"},
		{".sh", "Shell"},
		{".yaml", "YAML"},
		{".json", "JSON"},
		{".unknown", "Unknown"},
		{"", "Unknown"},
	}

	for _, tt := range tests {
		t.Run(tt.ext, func(t *testing.T) {
			result := reviewer.detectLanguage(tt.ext)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestFormatAsJSON(t *testing.T) {
	reviewer := &Reviewer{}

	reviewResult := "Code looks good!\nNo issues found."
	result := reviewer.formatAsJSON(reviewResult)

	assert.Contains(t, result, `"review"`)
	assert.Contains(t, result, "Code looks good!")
	// Check that special characters are escaped
	assert.Contains(t, result, "\\n")
}

func TestFormatAsMarkdown(t *testing.T) {
	reviewer := &Reviewer{}

	reviewResult := "=== CODE REVIEW RESULTS ===\n\nCRITICAL ISSUES:\nNo issues found."
	result := reviewer.formatAsMarkdown(reviewResult)

	assert.Contains(t, result, "# Code Review Results")
	assert.Contains(t, result, "## CODE REVIEW RESULTS")
	assert.Contains(t, result, "### CRITICAL ISSUES:")
}

func TestValidateOptions_EmptyFilePath(t *testing.T) {
	options := ReviewOptions{
		FilePath: "",
		Format:   FormatText,
	}

	err := ValidateOptions(options)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "file path cannot be empty")
}

func TestValidateOptions_InvalidFormat(t *testing.T) {
	options := ReviewOptions{
		FilePath: "test.go",
		Format:   OutputFormat("invalid"),
	}

	err := ValidateOptions(options)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid output format")
}

func TestValidateOptions_ValidOptions(t *testing.T) {
	tests := []struct {
		name   string
		format OutputFormat
	}{
		{"text format", FormatText},
		{"json format", FormatJSON},
		{"markdown format", FormatMarkdown},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			options := ReviewOptions{
				FilePath: "test.go",
				Format:   tt.format,
			}

			err := ValidateOptions(options)
			assert.NoError(t, err)
		})
	}
}

func TestReview_JSONFormat(t *testing.T) {
	// Setup
	mockClient := &MockAIClient{
		Response: &openai.ChatCompletionResponse{
			Choices: []openai.ChatCompletionChoice{
				{
					Message: openai.ChatCompletionMessage{
						Content: "Code review result",
					},
				},
			},
		},
	}

	mockFileReader := &MockFileReader{
		Content:  "package main",
		Size:     100,
		IsDir:    false,
		FileName: "test.go",
	}

	reviewer := &Reviewer{
		Client:     mockClient,
		FileReader: mockFileReader,
	}

	var output bytes.Buffer
	options := ReviewOptions{
		FilePath:     "test.go",
		Format:       FormatJSON,
		OutputWriter: &output,
	}

	conf := config.Config{}

	// Execute
	err := reviewer.Review(options, conf)

	// Assert
	assert.NoError(t, err)
	result := output.String()
	assert.Contains(t, result, `"review"`)
	assert.True(t, strings.HasPrefix(result, "{"))
}

func TestReview_MarkdownFormat(t *testing.T) {
	// Setup
	mockClient := &MockAIClient{
		Response: &openai.ChatCompletionResponse{
			Choices: []openai.ChatCompletionChoice{
				{
					Message: openai.ChatCompletionMessage{
						Content: "=== CODE REVIEW ===\n\nCRITICAL ISSUES:\nNone",
					},
				},
			},
		},
	}

	mockFileReader := &MockFileReader{
		Content:  "package main",
		Size:     100,
		IsDir:    false,
		FileName: "test.go",
	}

	reviewer := &Reviewer{
		Client:     mockClient,
		FileReader: mockFileReader,
	}

	var output bytes.Buffer
	options := ReviewOptions{
		FilePath:     "test.go",
		Format:       FormatMarkdown,
		OutputWriter: &output,
	}

	conf := config.Config{}

	// Execute
	err := reviewer.Review(options, conf)

	// Assert
	assert.NoError(t, err)
	result := output.String()
	assert.Contains(t, result, "# Code Review Results")
	assert.Contains(t, result, "##")
}

func TestReview_WithFocusAreas(t *testing.T) {
	// Setup
	mockClient := &MockAIClient{
		Response: &openai.ChatCompletionResponse{
			Choices: []openai.ChatCompletionChoice{
				{
					Message: openai.ChatCompletionMessage{
						Content: "Security review complete",
					},
				},
			},
		},
	}

	mockFileReader := &MockFileReader{
		Content:  "package main",
		Size:     100,
		IsDir:    false,
		FileName: "test.go",
	}

	reviewer := &Reviewer{
		Client:     mockClient,
		FileReader: mockFileReader,
	}

	var output bytes.Buffer
	options := ReviewOptions{
		FilePath:     "test.go",
		Format:       FormatText,
		FocusAreas:   []string{"security", "performance"},
		OutputWriter: &output,
	}

	conf := config.Config{}

	// Execute
	err := reviewer.Review(options, conf)

	// Assert
	assert.NoError(t, err)
	assert.Contains(t, output.String(), "Security review complete")
}

func TestNewReviewer(t *testing.T) {
	// Setup
	mockClient := &MockAIClient{}

	// Execute
	reviewer := NewReviewer(mockClient)

	// Assert
	assert.NotNil(t, reviewer)
	assert.NotNil(t, reviewer.FileReader)
	assert.Equal(t, mockClient, reviewer.Client)
}

func TestDefaultFileReader_ReadFile(t *testing.T) {
	// This test requires an actual file
	// We'll create a temporary file for testing
	content := "test content"
	tmpFile := "/tmp/test_aicodereview.txt"

	err := os.WriteFile(tmpFile, []byte(content), 0644)
	assert.NoError(t, err)
	defer os.Remove(tmpFile)

	reader := &DefaultFileReader{}
	data, err := reader.ReadFile(tmpFile)

	assert.NoError(t, err)
	assert.Equal(t, content, string(data))
}

func TestDefaultFileReader_Stat(t *testing.T) {
	// Create a temporary file
	tmpFile := "/tmp/test_aicodereview_stat.txt"
	err := os.WriteFile(tmpFile, []byte("test"), 0644)
	assert.NoError(t, err)
	defer os.Remove(tmpFile)

	reader := &DefaultFileReader{}
	info, err := reader.Stat(tmpFile)

	assert.NoError(t, err)
	assert.NotNil(t, info)
	assert.False(t, info.IsDir())
	assert.Greater(t, info.Size(), int64(0))
}
