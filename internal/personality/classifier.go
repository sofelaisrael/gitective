package personality

import (
	"strings"
	"time"

	commit "github.com/sofelaisrael/gitective/internal/commit"
)

type Personality string

const (
	PanicPatch     Personality = "PANIC_PATCH"
	NuclearJanitor Personality = "NUCLEAR_JANITOR"
	MadScientist   Personality = "MAD_SCIENTIST"
	SilentCommit   Personality = "SILENT_COMMIT"
	Shipper        Personality = "SHIPPER"
	Architect      Personality = "ARCHITECT"
	Janitor        Personality = "JANITOR"
)

func Classify(facts commit.CommitFacts) Personality {
	msg := strings.ToLower(facts.Message)
	isCleanup := strings.Contains(msg, "cleanup") || strings.Contains(msg, "clean")
	isFeat := strings.Contains(msg, "feat")
	isDocs := strings.Contains(msg, "docs")
	isTest := strings.Contains(msg, "test")
	if isCleanup && facts.FilesChanged > 15 && facts.LinesDeleted > 500 {
		return NuclearJanitor
	}
	if isCleanup && facts.FilesChanged > 10 {
		return Janitor
	}
	if facts.IsFix && facts.FilesChanged <= 3 && facts.LinesAdded < 10 {
		h := facts.Timestamp.In(time.Local).Hour()
		if h < 6 || h >= 22 {
			return PanicPatch
		}
		return PanicPatch
	}
	if isFeat && facts.LinesAdded > 200 && facts.FilesChanged > 5 && isDocs && isTest {
		return Architect
	}
	if isFeat && facts.LinesAdded > 200 && facts.FilesChanged > 5 {
		return MadScientist
	}
	if len(strings.TrimSpace(facts.Message)) < 5 {
		return SilentCommit
	}
	return Shipper
}
