# Convcommit

A CLI tool that simplifies the process of making conventional commits.

## What are Conventional Commits?

Conventional Commits is a specification for adding human and machine-readable meaning to commit messages. It provides an easy set of rules for creating an explicit commit history, which makes it easier to write automated tools on top of.

The commit message should be structured as follows:

```
<type>[optional scope]: <description>

[optional body]

[optional footer(s)]
```

## Features

- 🚀 Interactive commit generation with structured prompts
- 🌍 Global accessibility from any Git repository
- 🔄 One-line installer script
- 🎨 Color-coded prompts for better UX
- 📊 Feedback mechanism after 5 successful commits

## Installation

### One-line Installer

```bash
curl -sSL https://raw.githubusercontent.com/danfordChris/conventional_commit_CLI/main/installer.sh | bash
```

This script will:
- Download and install the binary globally
- Handle any dependencies
- Add it to your PATH (if needed)

### Manual Installation

If you prefer to install manually:

1. Clone the repository:
   ```bash
   git clone https://github.com/danfordChris/conventional_commit_CLI.git convcommit
   ```

2. Build the binary:
   ```bash
   cd convcommit/cmd/convcommit
   go mod tidy  # Ensure all dependencies are downloaded
   go build -o convcommit
   ```

   Note: You must build from the cmd/convcommit directory, not from the repository root.

3. Move the binary to a directory in your PATH:
   ```bash
   sudo mv convcommit /usr/local/bin/
   ```

## Usage

### Staging Changes

Before committing, you need to have changes staged in your Git repository. You can either:

1. Stage changes manually using Git commands:
   ```bash
   git add <files>  # Stage specific files
   git add .        # Stage all changes
   ```

2. Let convcommit prompt you to stage all changes automatically if no changes are staged when you run the command. The tool will provide helpful error messages if specific files can't be staged (e.g., if they don't exist or are ignored by .gitignore).

### Interactive Mode

Simply run the command without any arguments to enter interactive mode:

```bash
convcommit
```

You'll be prompted for:
1. Type (e.g., feat, fix, chore, etc.)
2. Scope (optional)
3. Title (mandatory, short line)
4. Description (optional, multiline)
5. Body (optional, multiline)
6. Footer (optional)

### Command-line Arguments

You can also use command-line arguments:

```bash
convcommit --type=feat --scope=ui --title="add login button" --push
```

Available options:
- `--type`, `-t`: Commit type (required)
- `--scope`, `-s`: Commit scope (optional)
- `--title`, `-d`: Commit title (required)
- `--body`, `-b`: Commit body (optional)
- `--breaking`, `-!`: Indicates a breaking change
- `--footer`, `-f`: Commit footer (optional)
- `--review`, `-r`: Review commit before pushing
- `--push`, `-p`: Push commit to remote repository
- `--suggestions`: Generate commit suggestions based on staged files

## Examples

### Interactive Mode

```
$ convcommit
Select commit type: feat
Enter scope (optional): login
Enter short title: add login API integration
Enter description (optional): 
Enter body (optional): 
Enter footer (optional): closes #12

✅ Commit Created:
feat(login): add login API integration

closes #12
```

### Smart Suggestions Mode

```
$ convcommit --suggestions

🔍 Detected changes:
- internal/auth/login.go

💡 Suggested type: feat
💡 Suggested scope: auth
💡 Suggested title: add auth feature
Use these suggestions? [Y/n]: Y

✅ Commit Created:
feat(auth): add auth feature
```

The `--suggestions` flag analyzes your staged files and suggests appropriate commit elements based on file paths and types. It will:

1. Check for staged files (or prompt you to stage changes if none are staged)
2. Analyze file paths to suggest appropriate type, scope, and title
3. Allow you to accept the suggestions or fall back to regular prompts

This is especially useful for quickly creating conventional commits without having to manually specify each element.

After your 5th commit, you'll be prompted for feedback:

```
How satisfied are you with this tool (1-5)? 4
Any suggestions? Add emoji support
✅ Feedback sent to jurvisdanford329@gmail.com
```

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the MIT License - see the LICENSE file for details.
       ___🎉 Happy Coding!
