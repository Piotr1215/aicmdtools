package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/piotr1215/aicmdtools/internal/aicodereview"
)

const (
	version     = "v0.0.1"
	prompt_file = "code-review-prompt.txt"
)

func main() {
	// Define command-line flags
	versionFlag := flag.Bool("version", false, "Display version information")
	modelFlag := flag.Bool("model", false, "Display current model")
	fileFlag := flag.String("f", "", "File to review (required)")
	formatFlag := flag.String("format", "text", "Output format: text, json, markdown")
	focusFlag := flag.String("focus", "", "Focus areas (comma-separated): security,performance,style")
	maxSizeFlag := flag.Int64("max-size", 1024*1024, "Maximum file size in bytes (default: 1MB)")

	flag.Parse()

	// Handle version flag
	if *versionFlag {
		fmt.Printf("aicodereview version: %s\n", version)
		return
	}

	// Handle model flag
	if *modelFlag {
		// This would typically read from config
		fmt.Println("Current model: gpt-4 (from config)")
		fmt.Println("To change model, edit ~/.config/aicmdtools/config.yaml")
		return
	}

	// Validate file flag
	if *fileFlag == "" {
		fmt.Println("Error: file path is required")
		fmt.Println("\nUsage:")
		fmt.Println("  aicodereview -f <file_path> [options]")
		fmt.Println("\nOptions:")
		flag.PrintDefaults()
		fmt.Println("\nExamples:")
		fmt.Println("  aicodereview -f main.go")
		fmt.Println("  aicodereview -f main.go -format json")
		fmt.Println("  aicodereview -f main.go -focus security,performance")
		fmt.Println("  aicodereview -f main.go -format markdown -max-size 2097152")
		os.Exit(1)
	}

	// Parse output format
	var format aicodereview.OutputFormat
	switch *formatFlag {
	case "text":
		format = aicodereview.FormatText
	case "json":
		format = aicodereview.FormatJSON
	case "markdown":
		format = aicodereview.FormatMarkdown
	default:
		fmt.Printf("Error: invalid format '%s' (valid: text, json, markdown)\n", *formatFlag)
		os.Exit(1)
	}

	// Parse focus areas
	var focusAreas []string
	if *focusFlag != "" {
		// Simple comma-split (could be enhanced with better parsing)
		focusAreas = splitAndTrim(*focusFlag, ",")
	}

	// Build review options
	options := aicodereview.ReviewOptions{
		FilePath:     *fileFlag,
		Format:       format,
		FocusAreas:   focusAreas,
		MaxFileSize:  *maxSizeFlag,
		OutputWriter: os.Stdout,
	}

	// Validate options
	if err := aicodereview.ValidateOptions(options); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	// Execute review
	if err := aicodereview.Execute(prompt_file, options); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}

// splitAndTrim splits a string by delimiter and trims whitespace from each part
func splitAndTrim(s string, delimiter string) []string {
	if s == "" {
		return nil
	}

	parts := []string{}
	current := ""

	for _, ch := range s {
		if string(ch) == delimiter {
			if trimmed := trim(current); trimmed != "" {
				parts = append(parts, trimmed)
			}
			current = ""
		} else {
			current += string(ch)
		}
	}

	// Add last part
	if trimmed := trim(current); trimmed != "" {
		parts = append(parts, trimmed)
	}

	return parts
}

// trim removes leading and trailing whitespace
func trim(s string) string {
	start := 0
	end := len(s)

	// Trim leading whitespace
	for start < end && (s[start] == ' ' || s[start] == '\t' || s[start] == '\n' || s[start] == '\r') {
		start++
	}

	// Trim trailing whitespace
	for end > start && (s[end-1] == ' ' || s[end-1] == '\t' || s[end-1] == '\n' || s[end-1] == '\r') {
		end--
	}

	return s[start:end]
}
