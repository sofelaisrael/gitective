package llm

import "fmt"

func BuildSummarizePrompt(diff string) string {
	return fmt.Sprintf(`Summarize these code changes as a git commit message.

Rules:
- First line: conventional commit prefix (feat:, fix:, chore:, refactor:, etc) + summary under 50 chars
- Blank line
- 1-2 sentence body: what changed and why
- MAX 30 words total
- Do NOT describe every file. Summarize the overall change.
- Output ONLY the commit message, nothing else

Changes:

%s`, diff)
}

func BuildStyledCommitPrompt(diff string) string {
	return fmt.Sprintf(`Write a git commit message for these changes.

CRITICAL RULES:
- MAXIMUM 25 words. This is a hard limit. Do NOT exceed it.
- First line: conventional commit prefix (feat:, fix:, chore:, refactor:, etc) + short summary
- Blank line, then 1 sentence explaining what and why
- Do NOT describe every file change. Summarize the overall change.
- Do NOT expand short text into long prose. Keep it tight.
- Output ONLY the commit message. Nothing else.

Changes:

%s`, diff)
}
