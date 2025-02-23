# lite-tie

`lite-tie` is a lightweight CLI tool designed to create and manage symbolic links for portable software in its own directory.

## Installation

Download the latest release from the [Releases](https://github.com/D7x7z49/lite-tie/releases) page and extract it to a directory in your `PATH`.

**Note for Windows Users**: Enabling **Developer Mode** (Settings > Update & Security > For developers) is recommended to create symbolic links without requiring admin privileges.

## Usage
```bash
lite-tie [command]
```

### Available Commands
- `add <exec_path> [--alias <name>]` - Add a symlink for a portable executable
- `list [<name>] [--simple]` - List all symlinks (or a specific one)
- `remove [<name> ...] [--silent] [--clean]` - Remove symlinks (silent or clean unavailable entries)

### Examples
```bash
# Add a symlink
lite-tie add ./myapp.exe --alias app

# List all symlinks
lite-tie list

# Remove specific symlinks silently
lite-tie remove app --silent

# Clean unavailable symlinks
lite-tie remove --clean
```

## Acknowledgments
Special thanks to **Grok 3**, created by xAI, for its invaluable assistance throughout the development and release process of `lite-tie`. Your guidance and insights made this project possible!