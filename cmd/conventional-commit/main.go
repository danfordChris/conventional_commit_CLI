package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var (
	commitType  string
	scope       string
	title       string
	description string
	body        string
	breaking    bool
	footer      string
	review      bool
	push        bool
	interactive bool
)

func main() {
	// Create the root command
	rootCmd := &cobra.Command{
		Use:   "convcommit",
		Short: "A CLI tool for creating conventional commits",
		Long: `Convcommit is a CLI tool that simplifies the process of making conventional commits.
It guides you through the commit process and ensures your commits follow the conventional commits specification.`,
		Run: func(cmd *cobra.Command, args []string) {
			// Check if the current directory is a Git repository
			if !IsGitRepository() {
				color.Red("Error: Not a Git repository. Please run this command inside a Git repository.")
				os.Exit(1)
			}

			var commit CommitMessage

			// Check if running in interactive mode
			if interactive || (commitType == "" && title == "") {
				// Run in interactive mode
				commit = PromptForCommit()
			} else {
				// Run in command-line mode
				if commitType == "" {
					color.Red("Error: --type is required in non-interactive mode.")
					os.Exit(1)
				}
				if title == "" {
					color.Red("Error: --title is required in non-interactive mode.")
					os.Exit(1)
				}
				if !isValidType(commitType) {
					color.Red("Error: Invalid commit type. Valid types are: %s", strings.Join(validTypes, ", "))
					os.Exit(1)
				}

				commit = CommitMessage{
					Type:        commitType,
					Scope:       scope,
					Description: title,
					Body:        body,
					Breaking:    breaking,
					Footer:      footer,
				}

				// If description is provided, add it to the body
				if description != "" {
					if commit.Body != "" {
						commit.Body = description + "\n\n" + commit.Body
					} else {
						commit.Body = description
					}
				}
			}

			// Generate the commit message
			commitMessage := GenerateCommitMessage(commit)

			// Print the generated commit message
			color.Green("\n✅ Commit Created:")
			fmt.Println(commitMessage)

			// Commit the changes
			err := CommitChanges(commitMessage)
			if err != nil {
				color.Red("Error: %s", err)
				os.Exit(1)
			}

			color.Green("Changes committed successfully.")

			// If push is requested, push the changes to the remote repository
			if push {
				color.Yellow("Pushing changes to remote repository...")
				pushCmd := exec.Command("git", "push")
				output, err := pushCmd.CombinedOutput()
				if err != nil {
					color.Red("Error: Failed to push changes: %s", string(output))
					os.Exit(1)
				}
				color.Green("Changes pushed successfully.")
			}
		},
	}

	// Add flags to the root command
	rootCmd.Flags().StringVarP(&commitType, "type", "t", "", "Commit type (required)")
	rootCmd.Flags().StringVarP(&scope, "scope", "s", "", "Commit scope (optional)")
	rootCmd.Flags().StringVarP(&title, "title", "m", "", "Commit title (required)")
	rootCmd.Flags().StringVarP(&description, "description", "d", "", "Commit description (optional)")
	rootCmd.Flags().StringVarP(&body, "body", "b", "", "Commit body (optional)")
	rootCmd.Flags().BoolVarP(&breaking, "breaking", "!", false, "Indicates a breaking change")
	rootCmd.Flags().StringVarP(&footer, "footer", "f", "", "Commit footer (optional)")
	rootCmd.Flags().BoolVarP(&review, "review", "r", false, "Review commit before pushing")
	rootCmd.Flags().BoolVarP(&push, "push", "p", false, "Push commit to remote repository")
	rootCmd.Flags().BoolVarP(&interactive, "interactive", "i", false, "Run in interactive mode")

	// Add version command
	versionCmd := &cobra.Command{
		Use:   "version",
		Short: "Print the version number",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Convcommit v1.0.0")
		},
	}
	rootCmd.AddCommand(versionCmd)

	// Execute the root command
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
