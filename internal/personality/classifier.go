package personality

import (
	"strings"

	commit "github.com/sofelaisrael/gitective/internal/commit"
)

type Personality string

const (
	PanicPatch     Personality = "PANIC_PATCH"
	NuclearJanitor Personality = "NUCLEAR_JANITOR"
	MadScientist   Personality = "MAD_SCIENTIST"
	SilentCommit   Personality = "SILENT_COMMIT"
	Shipper        Personality = "SHIPPER"
)

func Classify(facts commit.CommitFacts) Personality {
	if facts.FilesChanged > 15 && facts.LinesDeleted > 500 && !facts.IsFix && len(facts.Message) > 0 && (strings.Contains(strings.ToLower(facts.Message), "cleanup") || strings.Contains(strings.ToLower(facts.Message), "clean")) {
		return NuclearJanitor
	} else if facts.IsFix && facts.LinesAdded < 10 && facts.FilesChanged <= 3 {
		return PanicPatch
	} else if facts.LinesAdded > 200 && facts.FilesChanged > 5 {
		return MadScientist
	} else if len(strings.TrimSpace(facts.Message)) < 5 {
		return SilentCommit
	} else {
		return Shipper
	}
}
