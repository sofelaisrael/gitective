# gitective 🔍

> **Your Git history tells a story. You've just never heard it speak.**

`gitective` reads your Git history and translates raw commit data — timestamps, file changes, line counts, messages — into something human: a *personality classification* for each commit, rendered in themed terminal output.

Three commits at 2am that each touch 15 files? That's not "feat: update" — that's a **panic patch**. A 4,000-line commit with a single-word message? That's a **silent commit**. Your Git log already knows who you are. `gitective` just makes it say it out loud.

---

## What Does It Actually Do?

`gitective` detects the Git repository in your current directory, reads your commit history, extracts behavioral *facts* from each commit, and classifies your development pattern into one of several personalities:

```
// EMERGENCY PATCH
Operator: you
Breach: fix: broken auth flow after midnight
Sectors: 3  +9/-2
Action: hot-patched at 02:00
Status: CONTAINED.
```

That's the **Cyberpunk** theme. The same commit in the **Renaissance** theme:

```
⚜ THE CHRONICLES [THE HASTY MENDER]

Upon this day, keeper you hath wrought 'fix: broken auth flow after midnight'.
Scrolls: 3  +9/-2
Thus endureth.
```

It's not generating funny messages for the sake of it. It's *inferring behavior* from the data you already committed.

---

## The Personalities

Every commit gets classified based on *what it actually did* — not what the message says:

| Personality | What It Means |
|---|---|
| **PanicPatch** | Small fix, late hours, tight files. You woke up to a broken build. |
| **NuclearJanitor** | Massive cleanup — 15+ files, 500+ lines deleted. The codebase needed surgery. |
| **MadScientist** | Big feature, 5+ files, 200+ lines added. No docs, no tests. Just vibes. |
| **Architect** | Big feature *with* docs *and* tests. The responsible one. |
| **SilentCommit** | Commit message under 5 characters. Git history's cryptic moods. |
| **Janitor** | Moderate cleanup across many files. Steady, unglamorous maintenance. |
| **Shipper** | The default. You shipped something. Keep going. |

The classification isn't random — it looks at hour of day, lines added/deleted, files changed, and message content. A fix at 3am is a `PanicPatch` whether you admit it or not.

---

## Installation

```bash
go install github.com/sofelaisrael/gitective/cmd/gitective@latest
```

Requires **Go 1.24+**.

Or build from source:

```bash
git clone https://github.com/sofelaisrael/gitective.git
cd gitective
go build -o gitective ./cmd/gitective
```

---

## Usage

Run it from inside any Git repository:

```bash
cd your-project
gitective
```

Or without installing:

```bash
go run ./cmd/gitective
```

It will:
1. Detect the repository in your current directory
2. Read your commit history
3. Extract facts from your latest commit
4. Render the themed personality output

---

## Themes

`gitective` ships with two visual themes that completely change the output tone:

### Cyberpunk
Neon-drenched, operator-style reports. Commits become "breach reports," file changes are "sectors," and everything reads like a terminal from *Blade Runner*.

```
// FREIGHT RUN
Operator: you
Payload: feat: add dark mode toggle
Sectors: 4  +187/-32
Status: DELIVERED.
```

### Renaissance
Classical, chronicle-style prose. Your commits become entries in a medieval manuscript, and you are the keeper.

```
⚜ THE CHRONICLES [THE CARAVAN]

Upon this day, keeper you hath wrought 'feat: add dark mode toggle'.
Scrolls: 4  +187/-32
Thus endureth.
```

Both themes use the same underlying personality classification — just different vocabulary. Pick the one that matches your mood.

---

## Architecture

```
cmd/gitective/
  main.go                  CLI entrypoint — detect repo, read history, render

internal/
  analysis/
    activity.go            Commit activity pattern analysis
    timing.go              Commit timing analysis
  commit/
    facts.go               Extract behavioral facts from a commit
  git/
    repo.go                Repository detection, commit reading
  personality/
    classifier.go          Pattern → personality classification engine
  themes/
    theme.go               Theme interface
    cyberpunk.go           Cyberpunk visual theme
    renaissance.go         Renaissance visual theme
  ui/
    banner.go              ASCII art banners (go-figure)
    spinner.go             Terminal spinners (pterm)
```

**Dependencies:**
- [`charmbracelet/lipgloss`](https://github.com/charmbracelet/lipgloss) — terminal styling
- [`go-figure`](https://github.com/common-nighthawk/go-figure) — ASCII art banners
- [`pterm`](https://github.com/pterm/pterm) — terminal UI (spinners, colors)

---

## Why?

Because `git log --oneline` tells you *what* happened. `gitective` tells you *why it happened that way*.

I built this because I noticed something about my own commit history: the timestamps and file counts were screaming at me. My best code came in focused 3-file commits during the day. My worst code came in 15-file "just fix it" commits at 2am. The data was already there — I just needed something to read it back to me.

`gitective` is that mirror. It's a small tool, but the insight is real: **how you commit says as much as what you commit**.

---

## Contributing

This is an early prototype — 24 commits over 4 days. The personality classifications are intentionally opinionated and a little silly. That's the point.

If you want to add a personality, a new theme, or improve the classification logic, open a PR. The codebase is small, the architecture is flat, and everything lives under `internal/`.

```bash
# Run in development
go run ./cmd/gitective

# Build
go build -o gitective ./cmd/gitective

# Run tests
go test ./...
```

---

## License

MIT — use it, fork it, make it weirder.
