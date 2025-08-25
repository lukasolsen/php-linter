package reporter

import (
	"fmt"

	"github.com/codevault-llc/php-lint/pkg/types"
	"github.com/fatih/color"
)

// Render formats and prints the final linting report to the console.
func Render(results map[string][]types.Issue) {
	totalIssues := 0
	filesWithIssues := 0

	// Pre-calculate stats
	for _, issues := range results {
		if len(issues) > 0 {
			filesWithIssues++
			totalIssues += len(issues)
		}
	}

	if totalIssues == 0 {
		color.New(color.FgGreen).Println("\n✅ No issues found. Your code is looking great!")
		return
	}

	// Create color functions
	errorColor := color.New(color.FgRed).Add(color.Bold)
	warningColor := color.New(color.FgYellow)
	filePathColor := color.New(color.Underline)
	ruleColor := color.New(color.Faint)

	fmt.Println() // Add a newline for spacing

	for file, issues := range results {
		if len(issues) == 0 {
			continue
		}

		filePathColor.Println(file)
		for _, issue := range issues {
			// In a real implementation, you'd populate the issue's Line and Col
			// from the AST node's position. For now, we'll omit them.
			warningColor.Printf("  %s ", "warning")
			fmt.Printf(" %s ", issue.Message)
			ruleColor.Printf(" (%s)\n", issue.RuleName)
		}
		fmt.Println()
	}

	summary := fmt.Sprintf("\n✖ %d problem(s) found in %d file(s).", totalIssues, filesWithIssues)
	errorColor.Println(summary)
}