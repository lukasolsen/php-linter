package reporter

import (
	"fmt"

	"github.com/codevault-llc/php-lint/pkg/types"
	"github.com/fatih/color"
)

// Render formats and prints the final linting report to the console.
func Render(issues []types.Issue) {
	groupedIssues := make(map[string][]types.Issue)
	for _, issue := range issues {
		groupedIssues[issue.Source] = append(groupedIssues[issue.Source], issue)
	}

	// Create color functions
	errorColor := color.New(color.FgRed).Add(color.Bold)
	warningColor := color.New(color.FgYellow)
	filePathColor := color.New(color.Underline)
	ruleColor := color.New(color.Faint)

	fmt.Println() // Add a newline for spacing

	for _, issueList := range groupedIssues {
		if len(issueList) == 0 {
			continue
		}

		for _, issue := range issueList {
			filePathColor.Println(issue.Source)
			for _, issue := range issueList {
				// In a real implementation, you'd populate the issue's Line and Col
				// from the AST node's position. For now, we'll omit them.
			warningColor.Printf("  %s ", "warning")
			fmt.Printf(" %s ", issue.Message)
			ruleColor.Printf(" (%s)\n", issue.RuleName)
		}
		fmt.Println()
	}

	summary := fmt.Sprintf("\n✖ %d problem(s) found in %d file(s).", len(groupedIssues), len(issues))
	errorColor.Println(summary)
}
}