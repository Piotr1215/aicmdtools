package aicodereview

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/piotr1215/aicmdtools/internal/config"
	"github.com/piotr1215/aicmdtools/internal/nlp"
)

// OutputFormat represents the format of the review output
type OutputFormat string

const (
	FormatText     OutputFormat = "text"
	FormatJSON     OutputFormat = "json"
	FormatMarkdown OutputFormat = "markdown"
)

// ReviewOptions contains options for code review
type ReviewOptions struct {
	FilePath     string
	Format       OutputFormat
	FocusAreas   []string // e.g., "security", "performance", "style"
	MaxFileSize  int64    // Maximum file size in bytes to review
	OutputWriter io.Writer
}

// FileReader interface for dependency injection
type FileReader interface {
	ReadFile(path string) ([]byte, error)
	Stat(path string) (os.FileInfo, error)
}

// DefaultFileReader implements FileReader using os package
type DefaultFileReader struct{}

func (d *DefaultFileReader) ReadFile(path string) ([]byte, error) {
	return os.ReadFile(path)
}

func (d *DefaultFileReader) Stat(path string) (os.FileInfo, error) {
	return os.Stat(path)
}

// Reviewer handles code review operations
type Reviewer struct {
	Client     nlp.GAIClient
	FileReader FileReader
}

// NewReviewer creates a new Reviewer instance
func NewReviewer(client nlp.GAIClient) *Reviewer {
	return &Reviewer{
		Client:     client,
		FileReader: &DefaultFileReader{},
	}
}

// Execute performs the code review
func Execute(promptFile string, options ReviewOptions) error {
	// Read and parse config
	conf, prompt, err := config.ReadAndParseConfig("config.yaml", promptFile)
	if err != nil {
		return fmt.Errorf("failed to read config: %w", err)
	}

	// Create AI client
	client := nlp.CreateOpenAIClient(*conf)
	aiClient := nlp.GoaiClient{Client: client, Prompt: prompt}

	// Create reviewer
	reviewer := NewReviewer(&aiClient)

	// Set default output writer if not provided
	if options.OutputWriter == nil {
		options.OutputWriter = os.Stdout
	}

	// Perform review
	return reviewer.Review(options, *conf)
}

// Review performs the actual code review
func (r *Reviewer) Review(options ReviewOptions, conf config.Config) error {
	// Validate file path
	if options.FilePath == "" {
		return fmt.Errorf("file path is required")
	}

	// Check file exists and get info
	fileInfo, err := r.FileReader.Stat(options.FilePath)
	if err != nil {
		return fmt.Errorf("failed to stat file: %w", err)
	}

	// Check if it's a directory (not supported yet)
	if fileInfo.IsDir() {
		return fmt.Errorf("directory review not yet supported, please provide a file path")
	}

	// Check file size
	maxSize := options.MaxFileSize
	if maxSize == 0 {
		maxSize = 1024 * 1024 // Default 1MB
	}
	if fileInfo.Size() > maxSize {
		return fmt.Errorf("file size (%d bytes) exceeds maximum allowed size (%d bytes)", fileInfo.Size(), maxSize)
	}

	// Read file content
	content, err := r.FileReader.ReadFile(options.FilePath)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	// Build user prompt
	userPrompt := r.buildUserPrompt(options.FilePath, string(content), options.FocusAreas)

	// Call AI for review
	ctx := context.Background()
	response, err := r.Client.ProcessCommandWithContext(ctx, userPrompt, conf)
	if err != nil {
		return fmt.Errorf("failed to get AI review: %w", err)
	}

	// Extract review result
	if len(response.Choices) == 0 {
		return fmt.Errorf("no response from AI")
	}
	reviewResult := response.Choices[0].Message.Content

	// Format and output result
	return r.outputReview(reviewResult, options)
}

// buildUserPrompt constructs the user prompt for code review
func (r *Reviewer) buildUserPrompt(filePath string, content string, focusAreas []string) string {
	var sb strings.Builder

	// Add file info
	ext := filepath.Ext(filePath)
	fileName := filepath.Base(filePath)
	sb.WriteString(fmt.Sprintf("File: %s\n", fileName))
	sb.WriteString(fmt.Sprintf("Type: %s\n\n", r.detectLanguage(ext)))

	// Add focus areas if specified
	if len(focusAreas) > 0 {
		sb.WriteString("Focus Areas: ")
		sb.WriteString(strings.Join(focusAreas, ", "))
		sb.WriteString("\n\n")
	}

	// Add code content
	sb.WriteString("Code to review:\n\n")
	sb.WriteString("```")
	if len(ext) > 0 {
		sb.WriteString(ext[1:]) // Remove the leading dot
	}
	sb.WriteString("\n")
	sb.WriteString(content)
	sb.WriteString("\n```\n")

	return sb.String()
}

// detectLanguage detects programming language from file extension
func (r *Reviewer) detectLanguage(ext string) string {
	languages := map[string]string{
		".go":    "Go",
		".js":    "JavaScript",
		".ts":    "TypeScript",
		".py":    "Python",
		".java":  "Java",
		".c":     "C",
		".cpp":   "C++",
		".cs":    "C#",
		".rb":    "Ruby",
		".php":   "PHP",
		".rs":    "Rust",
		".swift": "Swift",
		".kt":    "Kotlin",
		".sh":    "Shell",
		".bash":  "Bash",
		".yaml":  "YAML",
		".yml":   "YAML",
		".json":  "JSON",
		".xml":   "XML",
		".html":  "HTML",
		".css":   "CSS",
		".sql":   "SQL",
	}

	if lang, ok := languages[ext]; ok {
		return lang
	}
	return "Unknown"
}

// outputReview formats and outputs the review result
func (r *Reviewer) outputReview(reviewResult string, options ReviewOptions) error {
	var output string

	switch options.Format {
	case FormatJSON:
		output = r.formatAsJSON(reviewResult)
	case FormatMarkdown:
		output = r.formatAsMarkdown(reviewResult)
	case FormatText:
		fallthrough
	default:
		output = reviewResult
	}

	// Write output
	_, err := fmt.Fprintln(options.OutputWriter, output)
	return err
}

// formatAsJSON converts review result to JSON format
func (r *Reviewer) formatAsJSON(reviewResult string) string {
	// Simple JSON wrapper - could be enhanced to parse the review result
	// and create structured JSON
	escaped := strings.ReplaceAll(reviewResult, "\"", "\\\"")
	escaped = strings.ReplaceAll(escaped, "\n", "\\n")
	return fmt.Sprintf(`{"review": "%s"}`, escaped)
}

// formatAsMarkdown converts review result to Markdown format
func (r *Reviewer) formatAsMarkdown(reviewResult string) string {
	// Add markdown headers and formatting
	var sb strings.Builder
	sb.WriteString("# Code Review Results\n\n")

	// The AI response is already well-formatted, just wrap it
	lines := strings.Split(reviewResult, "\n")
	for _, line := range lines {
		// Convert === headers to markdown headers
		if strings.HasPrefix(line, "===") && strings.HasSuffix(line, "===") {
			title := strings.Trim(line, "= ")
			sb.WriteString(fmt.Sprintf("## %s\n", title))
		} else if strings.HasSuffix(line, ":") && isHeaderLine(line) {
			// Category headers (e.g., "CRITICAL ISSUES:", "HIGH PRIORITY:", "SUMMARY:")
			sb.WriteString(fmt.Sprintf("### %s\n", line))
		} else {
			sb.WriteString(line)
			sb.WriteString("\n")
		}
	}

	return sb.String()
}

// isHeaderLine checks if a line looks like a header (all caps or starts with caps word)
func isHeaderLine(line string) bool {
	line = strings.TrimSpace(line)
	if len(line) == 0 {
		return false
	}

	// Remove trailing colon for checking
	checkLine := strings.TrimSuffix(line, ":")

	// Check if it's mostly uppercase letters and spaces
	upperCount := 0
	letterCount := 0
	for _, ch := range checkLine {
		if ch >= 'A' && ch <= 'Z' {
			upperCount++
			letterCount++
		} else if ch >= 'a' && ch <= 'z' {
			letterCount++
		}
	}

	// If at least 70% of letters are uppercase, it's likely a header
	if letterCount > 0 && float64(upperCount)/float64(letterCount) >= 0.7 {
		return true
	}

	return false
}

// ValidateOptions validates review options
func ValidateOptions(options ReviewOptions) error {
	if options.FilePath == "" {
		return fmt.Errorf("file path cannot be empty")
	}

	validFormats := map[OutputFormat]bool{
		FormatText:     true,
		FormatJSON:     true,
		FormatMarkdown: true,
	}

	if !validFormats[options.Format] {
		return fmt.Errorf("invalid output format: %s (valid: text, json, markdown)", options.Format)
	}

	return nil
}
