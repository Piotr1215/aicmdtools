# AI Command Line Tools

AICmdTools is set of command-line tools that utilize the OpenAI ChatGPT API to:

- generate shell commands based on user input
- activate command line chat with ability to pass through system settings
- perform AI-powered code reviews with structured feedback

It's built with Go and designed for ease of use,
providing a simple and efficient way to interact with ChatGPT API. Project
inspired by https://github.com/wunderwuzzi23/yolo-ai-cmdbot.

## Features

- Simple command-line interface
- Supports multiple operating systems and shells
- Configurable safety feature to confirm command execution
- Easy installation and setup

## Installation

1.  Clone the repository:

```bash
git clone https://github.com/piotr1215/aicmdtools.git
cd aicmdtools
```

2.  Install all the command-line tools:

Use `just` command runner to install the commands and optionally copy the
configuration files:

```bash
just install
```

If running for the first time, bun the provided installation script to set up configuration files:

This will copy the `config.yaml` and `prompt.txt` files to the appropriate location in your home directory (e.g., `$HOME/.config/aicmdtools`).

```bash
just copy_files 
```

## Usage

There are 4 separate commands that you can use:

### 1. aicmd - Command Generation
Generate a shell command based on user input.

```bash
aicmd "create a new directory called my_project"
aicmd -model    # Display current model
aicmd -version  # Display version
```

### 2. aichat - Interactive Chat
Start an interactive chat session with the AI model.

```bash
aichat
aichat -version
```

### 3. aicompgraph - Diagram Generation
Generate PlantUML diagrams from YAML files (useful for Crossplane compositions).

```bash
aicompgraph -f composition.yaml
aicompgraph -version
```

### 4. aicodereview - Code Review
Perform AI-powered code review with structured feedback on security, performance, and best practices.

```bash
aicodereview -f main.go
aicodereview -f main.go -format json
aicodereview -f main.go -format markdown
aicodereview -f main.go -focus security,performance
aicodereview -f main.go -max-size 2097152
aicodereview -version
```

**Options:**
- `-f`: File path to review (required)
- `-format`: Output format - `text` (default), `json`, or `markdown`
- `-focus`: Comma-separated focus areas (e.g., `security,performance,style`)
- `-max-size`: Maximum file size in bytes (default: 1048576 = 1MB)

**Output Categories:**
- **CRITICAL**: Security vulnerabilities, data loss risks, crashes
- **HIGH**: Bugs, logic errors, resource leaks
- **MEDIUM**: Performance issues, code smells
- **LOW**: Style, maintainability, best practices

## Commands

- `-model`: Display the current model being used (supported by `aicmd`)
- `-version`: Display the current version (supported by all CLIs)

## Configuration

You can customize the behavior of AICmdTools by modifying the `config.yaml` file located in `$HOME/.config/aicmdtools`. The available options include:

- `openai_api_key`: Your OpenAI API key.
  > alternatively the api key can be passed via variable `$OPENAI_API_KEY`
- `safety`: If set to `true`, AICmdTools will prompt you to confirm before executing any generated command.
- `model`: any supported model that you have access to
  > to list all available models use `curl https://api.openai.com/v1/models \
-H "Authorization: Bearer $OPENAI_API_KEY"`

### Prompts

Each tool uses a dedicated prompt file that can be customized in the config folder (`$HOME/.config/aicmdtools`):

- `prompt.txt` - aicmd command generation
- `chat-prompt.txt` - aichat interactive chat
- `comp-graph-prompt.txt` - aicompgraph diagram generation
- `code-review-prompt.txt` - aicodereview code review

You can edit these files to adjust how each tool behaves and what kind of responses it generates.

## Contributing

Contributions are welcome! If you have any ideas for improvements or bug fixes, please submit a pull request or create an issue on the GitHub repository.

## License

AICmdTools is released under the MIT License. See the `LICENSE` file for more information.
