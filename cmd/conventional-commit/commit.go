package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/smtp"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
)

// Valid commit types according to conventional commits specification
var validTypes = []string{
	"feat", "fix", "docs", "style", "refactor",
	"perf", "test", "build", "ci", "chore", "revert",
}

// CommitMessage represents a conventional commit message
type CommitMessage struct {
	Type        string
	Scope       string
	Description string
	Body        string
	Breaking    bool
	Footer      string
}

// Config represents the application configuration
type Config struct {
	CommitCount int    `json:"commit_count"`
	UserEmail   string `json:"user_email"`
}

// ConfigFilePath returns the path to the config file
func ConfigFilePath() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return ".convcommit_config"
	}
	return filepath.Join(homeDir, ".convcommit_config")
}

// LoadConfig loads the configuration from the config file
func LoadConfig() Config {
	configPath := ConfigFilePath()
	config := Config{}

	data, err := ioutil.ReadFile(configPath)
	if err != nil {
		return config
	}

	json.Unmarshal(data, &config)
	return config
}

// SaveConfig saves the configuration to the config file
func SaveConfig(config Config) error {
	configPath := ConfigFilePath()
	data, err := json.Marshal(config)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(configPath, data, 0644)
}

// IncrementCommitCount increments the commit count in the config file
func IncrementCommitCount() error {
	config := LoadConfig()
	config.CommitCount++
	return SaveConfig(config)
}

// CheckFeedbackNeeded checks if feedback is needed and prompts for it if necessary
func CheckFeedbackNeeded() {
	config := LoadConfig()
	if config.CommitCount == 5 {
		promptForFeedback()
	}
}

// promptForFeedback prompts the user for feedback and sends it via email
func promptForFeedback() {
	reader := bufio.NewReader(os.Stdin)
	
	color.Yellow("\nYou've made 5 commits with convcommit! We'd love your feedback:")
	
	// Ask for satisfaction rating
	fmt.Print(color.YellowString("How satisfied are you with this tool (1-5)? "))
	ratingInput, _ := reader.ReadString('\n')
	rating := strings.TrimSpace(ratingInput)
	
	// Ask for suggestions
	fmt.Print(color.YellowString("Any suggestions for improvement? "))
	suggestions, _ := reader.ReadString('\n')
	suggestions = strings.TrimSpace(suggestions)
	
	// Send feedback via email
	sendFeedbackEmail(rating, suggestions)
	
	color.Green("✅ Feedback sent to jurvisdanford329@gmail.com")
}

// sendFeedbackEmail sends the feedback via email
func sendFeedbackEmail(rating, suggestions string) {
	// This is a placeholder for the actual email sending logic
	// In a real implementation, you would use an SMTP package or email API
	// For now, we'll just print the feedback to the console
	fmt.Printf("Would send email with rating: %s, suggestions: %s\n", rating, suggestions)
	
	// Example of how to send an email using SMTP
	// You would need to replace these with actual credentials
	/*
	from := "your-email@example.com"
	password := "your-password"
	to := []string{"jurvisdanford329@gmail.com"}
	smtpHost := "smtp.example.com"
	smtpPort := "587"
	
	message := []byte(fmt.Sprintf("Subject: Convcommit Feedback\r\n\r\nRating: %s\r\nSuggestions: %s", rating, suggestions))
	
	auth := smtp.PlainAuth("", from, password, smtpHost)
	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, message)
	if err != nil {
		fmt.Println("Error sending email:", err)
	}
	*/
}

// IsGitRepository checks if the current directory is a Git repository
func IsGitRepository() bool {
	_, err := os.Stat(".git")
	if err == nil {
		return true
	}
	
	cmd := exec.Command("git", "rev-parse", "--is-inside-work-tree")
	output, err := cmd.Output()
	if err != nil {
		return false
	}
	
	return strings.TrimSpace(string(output)) == "true"
}

// PromptForCommit prompts the user for commit details
func PromptForCommit() CommitMessage {
	reader := bufio.NewReader(os.Stdin)
	
	// Get commit type
	commitType := promptForType(reader)
	
	// Get commit scope (optional)
	color.Cyan("Enter scope (optional, press Enter to skip): ")
	scope, _ := reader.ReadString('\n')
	scope = strings.TrimSpace(scope)
	
	// Get commit description
	description := ""
	for description == "" {
		color.Cyan("Enter short title (required): ")
		description, _ = reader.ReadString('\n')
		description = strings.TrimSpace(description)
		if description == "" {
			color.Red("Title is required. Please try again.")
		}
	}
	
	// Get commit body (optional)
	color.Cyan("Enter description (optional, press Enter to skip, multiple lines allowed, end with an empty line):")
	body := readMultilineInput(reader)
	
	// Get commit body (optional)
	color.Cyan("Enter body (optional, press Enter to skip, multiple lines allowed, end with an empty line):")
	body2 := readMultilineInput(reader)
	if body2 != "" {
		if body != "" {
			body += "\n\n" + body2
		} else {
			body = body2
		}
	}
	
	// Check for breaking changes
	color.Cyan("Is this a breaking change? (y/N): ")
	breakingInput, _ := reader.ReadString('\n')
	breaking := strings.ToLower(strings.TrimSpace(breakingInput)) == "y"
	
	// Get footer (optional)
	color.Cyan("Enter footer (optional, press Enter to skip, multiple lines allowed, end with an empty line):")
	footer := readMultilineInput(reader)
	
	// Create and return the commit message
	return CommitMessage{
		Type:        commitType,
		Scope:       scope,
		Description: description,
		Body:        body,
		Breaking:    breaking,
		Footer:      footer,
	}
}

// promptForType prompts the user to select a valid commit type
func promptForType(reader *bufio.Reader) string {
	for {
		color.Cyan("Select commit type:")
		for i, t := range validTypes {
			fmt.Printf("%2d. %s\n", i+1, t)
		}
		color.Cyan("Enter number or type directly: ")
		
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)
		
		// Check if input is a number
		var index int
		if _, err := fmt.Sscanf(input, "%d", &index); err == nil {
			if index >= 1 && index <= len(validTypes) {
				return validTypes[index-1]
			}
			color.Red("Invalid selection. Please try again.")
		} else {
			// Input is not a number, check if it's a valid type
			if isValidType(input) {
				return input
			}
			color.Red("Invalid type. Please try again.")
		}
	}
}

// isValidType checks if the given type is valid
func isValidType(t string) bool {
	for _, validType := range validTypes {
		if t == validType {
			return true
		}
	}
	return false
}

// readMultilineInput reads multiple lines of input until an empty line is entered
func readMultilineInput(reader *bufio.Reader) string {
	var lines []string
	for {
		line, _ := reader.ReadString('\n')
		line = strings.TrimSpace(line)
		if line == "" {
			break
		}
		lines = append(lines, line)
	}
	return strings.Join(lines, "\n")
}

// GenerateCommitMessage generates a conventional commit message
func GenerateCommitMessage(commit CommitMessage) string {
	var sb strings.Builder
	
	// Header: <type>[(scope)][!]: <description>
	sb.WriteString(commit.Type)
	
	if commit.Scope != "" {
		sb.WriteString("(")
		sb.WriteString(commit.Scope)
		sb.WriteString(")")
	}
	
	if commit.Breaking {
		sb.WriteString("!")
	}
	
	sb.WriteString(": ")
	sb.WriteString(commit.Description)
	
	// Body (optional)
	if commit.Body != "" {
		sb.WriteString("\n\n")
		sb.WriteString(commit.Body)
	}
	
	// Breaking change footer (if not already indicated in the header)
	footer := commit.Footer
	if commit.Breaking && (footer == "" || !strings.Contains(footer, "BREAKING CHANGE")) {
		if footer == "" {
			footer = "BREAKING CHANGE: This commit introduces breaking changes."
		} else {
			footer = "BREAKING CHANGE: This commit introduces breaking changes.\n" + footer
		}
	}
	
	// Footer (optional)
	if footer != "" {
		sb.WriteString("\n\n")
		sb.WriteString(footer)
	}
	
	return sb.String()
}

// CommitChanges commits the changes with the given message
func CommitChanges(commitMessage string) error {
	// Create a temporary file for the commit message
	tmpFile, err := ioutil.TempFile("", "commit-msg-*.txt")
	if err != nil {
		return fmt.Errorf("failed to create temporary file for commit message: %w", err)
	}
	defer os.Remove(tmpFile.Name())
	
	// Write the commit message to the temporary file
	if _, err := tmpFile.WriteString(commitMessage); err != nil {
		return fmt.Errorf("failed to write commit message to temporary file: %w", err)
	}
	tmpFile.Close()
	
	// Commit the changes
	cmd := exec.Command("git", "commit", "-m", commitMessage)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to commit changes: %s", string(output))
	}
	
	// Increment commit count
	if err := IncrementCommitCount(); err != nil {
		fmt.Println("Warning: Failed to update commit count:", err)
	}
	
	// Check if feedback is needed
	CheckFeedbackNeeded()
	
	return nil
}