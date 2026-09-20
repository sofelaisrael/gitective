package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/sofelaisrael/gitective/internal/git"
)

var (
	treeLine  = lipgloss.NewStyle().Foreground(lipgloss.Color("#555555"))
	treeLabel = lipgloss.NewStyle().Foreground(lipgloss.Color("#888888"))
	hashStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#666666")).Bold(true)
	authorStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#AAAAAA"))
	timeStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("#666666"))
	msgStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("#FFFFFF"))
)

func RenderCommitTree(commits []git.Commit, styles []git.CommitStyle) string {
	var b strings.Builder

	for i, c := range commits {
		isLast := i == len(commits)-1
		connector := "│"
		if isLast {
			connector = " "
		}

		hash := c.Hash
		if len(hash) > 7 {
			hash = hash[:7]
		}

		color := StyleColor(styles[i].Style)
		branchStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(color))

		branch := branchStyle.Render("●")
		hashStr := hashStyle.Render(hash)
		authorStr := authorStyle.Render(c.Author)
		timeStr := timeStyle.Render("·")

		styleLabel := treeLabel.Render(fmt.Sprintf(" [%s]", styles[i].Style))

		msg := c.Message
		if len(msg) > 60 {
			msg = msg[:57] + "..."
		}
		msgStr := msgStyle.Render(msg)

		b.WriteString(treeLine.Render(connector + "  "))
		b.WriteString(branch + " ")
		b.WriteString(hashStr + " ")
		b.WriteString(authorStr + " ")
		b.WriteString(timeStr + " ")
		b.WriteString(styleLabel)
		b.WriteString("\n")
		b.WriteString(treeLine.Render(connector + "  "))
		b.WriteString("   " + msgStr)
		b.WriteString("\n")

		if !isLast {
			b.WriteString(treeLine.Render("│"))
			b.WriteString("\n")
		}
	}

	return b.String()
}
