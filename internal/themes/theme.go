package themes

import "github.com/sofelaisrael/gitective/internal/commit"

type Theme interface {
	Name() string
	Render(facts commit.CommitFacts) string
}
