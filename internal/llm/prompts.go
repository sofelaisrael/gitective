package llm

import "fmt"

func BuildCommitUserPrompt(diff string) string {
	return fmt.Sprintf("Generate a git commit message for the following changes. Use conventional commit format (feat:, fix:, chore:, etc). Output ONLY the commit message, nothing else.\n\nChanges:\n\n%s", diff)
}
