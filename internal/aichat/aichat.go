package aichat

import (
	"fmt"
	"log"
	"strings"

	"github.com/alecthomas/chroma/formatters"
	"github.com/alecthomas/chroma/lexers"
	"github.com/alecthomas/chroma/styles"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/mitchellh/go-wordwrap"
	"github.com/piotr1215/aicmdtools/internal/config"
	"github.com/piotr1215/aicmdtools/internal/nlp"
	"github.com/piotr1215/aicmdtools/internal/utils"
)

var prompt_file = "chat-prompt.txt"

type (
	errMsg error
)
type model struct {
	aiClient    *nlp.GoaiClient
	viewport    viewport.Model
	messages    []string
	textarea    textarea.Model
	senderStyle lipgloss.Style
	err         error
}

func initialModel() model {
	ta := textarea.New()
	ta.Placeholder = "Send a message..."
	ta.Focus()

	ta.Prompt = "┃ "
	ta.CharLimit = 500

	ta.SetWidth(300)
	ta.SetHeight(3)

	// Remove cursor line styling
	ta.FocusedStyle.CursorLine = lipgloss.NewStyle()

	ta.ShowLineNumbers = false

	vp := viewport.New(300, 5)
	vp.SetContent(`Welcome to the chat room!
Type a message and press Enter to send.`)

	ta.KeyMap.InsertNewline.SetEnabled(false)

	return model{
		aiClient:    Initialize(),
		textarea:    ta,
		messages:    []string{},
		viewport:    vp,
		senderStyle: lipgloss.NewStyle().Foreground(lipgloss.Color("5")),
		err:         nil,
	}
}

// Calculate the number of lines in a string
func countLines(str string) int {
	return strings.Count(str, "\n") + 1
}

// Calculate the new height for the viewport
func calculateNewHeight(messages []string) int {
	totalLines := 0
	for _, msg := range messages {
		totalLines += countLines(msg)
	}
	return totalLines
}
func (m model) Init() tea.Cmd {
	return textarea.Blink
}
func wrapLines(text string, width uint) string {
	return wordwrap.WrapString(text, width)
}

func highlightCode(code string, lang string) string {
	// Get the lexer and style
	lexer := lexers.Get(lang)
	style := styles.Get("monokai")

	// Create a new buffer to hold the highlighted text
	var highlightedCode strings.Builder

	// Format the code
	formatter := formatters.Get("terminal256")
	iterator, _ := lexer.Tokenise(nil, code)
	formatter.Format(&highlightedCode, style, iterator)

	return highlightedCode.String()
}
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		tiCmd tea.Cmd
		vpCmd tea.Cmd
	)

	m.textarea, tiCmd = m.textarea.Update(msg)
	m.viewport, vpCmd = m.viewport.Update(msg)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			fmt.Println(m.textarea.Value())
			return m, tea.Quit
		case tea.KeyEnter:
			// Append your message to the chat
			m.messages = append(m.messages, m.senderStyle.Render("You: ")+m.textarea.Value())

			// Get the AI's response and append it to the chat
			aiResponse, err := SendMessage(m.aiClient, m.textarea.Value())
			if err != nil {
				m.messages = append(m.messages, m.senderStyle.Render("Error: ")+err.Error())
			} else {
				m.messages = append(m.messages, m.senderStyle.Render("AI: ")+aiResponse)
			}

			// Wrap the lines to fit within the viewport width
			wrappedLines := wrapLines(strings.Join(m.messages, "\n"), 150)

			// Update the viewport content and reset the textarea
			m.viewport.SetContent(wrappedLines)

			// Calculate the new height and update the viewport
			newHeight := calculateNewHeight(strings.Split(wrappedLines, "\n"))
			oldYPosition := m.viewport.YPosition

			// Create a new viewport with the new dimensions
			m.viewport = viewport.New(150, newHeight)
			m.viewport.YPosition = oldYPosition

			// Set the content again for the new viewport
			m.viewport.SetContent(wrappedLines)

			m.textarea.Reset()
			m.viewport.GotoBottom()
		}
	// We handle errors just like any other message
	case errMsg:
		m.err = msg
		return m, nil
	}

	return m, tea.Batch(tiCmd, vpCmd)
}

func (m model) View() string {
	return fmt.Sprintf(
		"%s\n\n%s",
		m.viewport.View(),
		m.textarea.View(),
	) + "\n\n"
}

func Initialize() *nlp.GoaiClient {
	// Read and parse the configuration
	configReader := &utils.FileReader{
		FilePathFunc: func() string { return config.ConfigFilePath("config.yaml") },
	}
	configContent := configReader.ReadFile()
	conf := config.ParseConfig(configContent)

	// Read and parse the prompt
	promptReader := &utils.FileReader{
		FilePathFunc: func() string { return config.ConfigFilePath(prompt_file) },
	}
	prompt := promptReader.ReadFile()
	operating_system, shell := utils.DetectOSAndShell()
	prompt = utils.ReplacePlaceholders(prompt, operating_system, shell)

	// Initialize OpenAI client
	client := nlp.CreateOpenAIClient(conf)

	return &nlp.GoaiClient{
		Client: client,
		Prompt: prompt,
	}
}

func SendMessage(client *nlp.GoaiClient, userMessage string) (string, error) {
	userMessage = strings.TrimSpace(userMessage) // Remove trailing newline

	response, err := client.ProcessCommand(userMessage)
	if err != nil {
		return "", err
	}

	return response.Choices[0].Message.Content, nil
}

func Execute() error {
	p := tea.NewProgram(initialModel())

	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}

	return nil
}
