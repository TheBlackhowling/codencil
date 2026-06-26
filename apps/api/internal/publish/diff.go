package publish

import "strings"

// LineDiffKind describes how a line changed between two markdown snapshots.
type LineDiffKind int

const (
	LineUnchanged LineDiffKind = iota
	LineAdded
	LineRemoved
)

// LineChange is one line-level change between old and new text.
type LineChange struct {
	Kind       LineDiffKind
	LineNumber int // 1-based line number in the new text (0 if removed only)
	Text       string
}

// DiffLines compares two markdown strings line-by-line (no anchor logic yet).
func DiffLines(oldMarkdown, newMarkdown string) []LineChange {
	oldLines := splitLines(oldMarkdown)
	newLines := splitLines(newMarkdown)

	// Simple longest-common-subsequence style diff via Myers is overkill for scaffold;
	// use a straightforward two-pointer merge for equal prefixes/suffixes + middle changes.
	prefix := 0
	for prefix < len(oldLines) && prefix < len(newLines) && oldLines[prefix] == newLines[prefix] {
		prefix++
	}

	suffix := 0
	for suffix < len(oldLines)-prefix && suffix < len(newLines)-prefix &&
		oldLines[len(oldLines)-1-suffix] == newLines[len(newLines)-1-suffix] {
		suffix++
	}

	var changes []LineChange
	for i := 0; i < prefix; i++ {
		changes = append(changes, LineChange{
			Kind:       LineUnchanged,
			LineNumber: i + 1,
			Text:       newLines[i],
		})
	}

	oldMid := oldLines[prefix : len(oldLines)-suffix]
	newMid := newLines[prefix : len(newLines)-suffix]

	// Middle: emit removals then additions (scaffold — refined in Phase 3 anchor remap).
	for _, line := range oldMid {
		changes = append(changes, LineChange{Kind: LineRemoved, LineNumber: 0, Text: line})
	}
	for i, line := range newMid {
		changes = append(changes, LineChange{
			Kind:       LineAdded,
			LineNumber: prefix + i + 1,
			Text:       line,
		})
	}

	for i := len(newLines) - suffix; i < len(newLines); i++ {
		changes = append(changes, LineChange{
			Kind:       LineUnchanged,
			LineNumber: i + 1,
			Text:       newLines[i],
		})
	}

	return changes
}

func splitLines(s string) []string {
	if s == "" {
		return nil
	}
	return strings.Split(s, "\n")
}
