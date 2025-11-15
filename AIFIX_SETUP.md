# aifix - Shell Integration Setup

`aifix` can automatically detect command errors when integrated with your shell.

## Quick Setup (ZSH)

1. **Generate the integration code:**
   ```bash
   aifix -init-shell zsh
   ```

2. **Add to your ~/.zshrc:**
   ```bash
   aifix -init-shell zsh >> ~/.zshrc
   ```

3. **Reload your shell:**
   ```bash
   source ~/.zshrc
   ```

4. **Test it:**
   ```bash
   $ ls -pap
   exa: Unknown argument -p
   $ fix
   ```

## Manual Setup (if you prefer to review first)

Run `aifix -init-shell zsh` and manually copy the output to your `~/.zshrc`.

The integration adds:
- Error capture to `/tmp/aifix_error_$$`
- Command tracking to `/tmp/aifix_cmd_$$`
- Alias `fix` for quick access
- Auto-cleanup on shell exit

## How It Works

1. **`preexec`** - Captures command before execution
2. **`precmd`** - Checks exit code after execution
3. **`exec 2>`** - Redirects stderr to capture file
4. **`aifix`** - Reads captured error and analyzes it

## Usage After Setup

### Automatic (with integration):
```bash
$ go build
# error: undefined: fmt.Println
$ fix
→ Instant AI-powered fix suggestion
```

### Manual (always works):
```bash
$ aifix "your error message here"
```

## Troubleshooting

**Q: "aifix says no error detected"**
A: Make sure you added the integration to ~/.zshrc and reloaded your shell

**Q: "stderr capture interferes with other tools"**
A: Comment out the `exec 2>` line in your ~/.zshrc

**Q: "I want to disable it temporarily"**
A: `unset AIFIX_ERROR_FILE AIFIX_CMD_FILE AIFIX_LAST_EXIT`

## Other Shells

- **Bash:** `aifix -init-shell bash >> ~/.bashrc`
- **Fish:** `aifix -init-shell fish >> ~/.config/fish/config.fish`

## Uninstall Integration

Remove the aifix section from your shell config file and reload:
```bash
# Edit ~/.zshrc and remove the aifix section
source ~/.zshrc
```
