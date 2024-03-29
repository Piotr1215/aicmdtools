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

    git clone https://github.com/piotr1215/aicmdtools.git
    cd aicmdtools

2.  Install the command-line tool:

    go install github.com/piotr1215/aicmdtools/cmd/aicmdtools@latest

3.  Run the provided installation script to set up configuration files:

        ./install.sh

    This script will copy the `config.yaml` and `prompt.txt` files to the appropriate location in your home directory (e.g., `$HOME/.config/aicmdtools`).

## Usage

To use AICmdTools, simply run the following command:

    aicmdtools "<your_input_here>"

Replace `<your_input_here>` with your desired input. For example:

    aicmdtools "create a new directory called my_project"

AICmdTools will generate a shell command based on the input, and if the safety feature is enabled in the configuration, it will prompt you to confirm whether you want to execute the command. If you confirm, the command will be executed in your shell.

## Configuration

You can customize the behavior of AICmdTools by modifying the `config.yaml` file located in `$HOME/.config/aicmdtools`. The available options include:

- `openai_api_key`: Your OpenAI API key.
  > alternatively the api key can be passed via variable `$OPENAI_API_KEY`
- `safety`: If set to `true`, AICmdTools will prompt you to confirm before executing any generated command.
- `model`: any model that you have access to
  > to list all available models use `curl https://api.openai.com/v1/models \
-H "Authorization: Bearer $OPENAI_API_KEY" `

### Prompt

It is possible to edit the `promt.txt` file in the config folder and make aicmdtools
behave in a different way if you want to adjust the prompt further.

## Contributing

Contributions are welcome! If you have any ideas for improvements or bug fixes, please submit a pull request or create an issue on the GitHub repository.

## License

AICmdTools is released under the MIT License. See the `LICENSE` file for more information.
