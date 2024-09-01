# AI Command Line Tools

AICmdTools is set of command-line tools that utilize the OpenAI CahtGPT API to:

- generate shell commands based on user input.
- activate command line chat with ability to pass through system settings.

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

There are 3 separate commands that you can use:

- `aicmd`: Generate a shell command based on user input.
> Example: aicmdtools "create a new directory called my_project"

- `aichat`: Start a chat with the AI model
- `aicompgraph`: Generate plantuml diagrams based YAML files (useful for Crossplane diagrams)

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

### Prompt

It is possible to edit the `promt.txt` file in the config folder and make aicmdtools
behave in a different way if you want to adjust the prompt further.

## Contributing

Contributions are welcome! If you have any ideas for improvements or bug fixes, please submit a pull request or create an issue on the GitHub repository.

## License

AICmdTools is released under the MIT License. See the `LICENSE` file for more information.
