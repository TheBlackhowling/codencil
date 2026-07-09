package publish

import "strings"

const (
	AnchorStatusActive   = "active"
	AnchorStatusShifted  = "shifted"
	AnchorStatusOrphaned = "orphaned"
)

// AnchorInput is one anchor from the previous published version.
type AnchorInput struct {
	AnchorID    string
	ThreadID    string
	StartLine   int
	EndLine     int
	QuotedText  string
	ReviewState string
}

// RemappedAnchor is the anchor placement on a new published version.
type RemappedAnchor struct {
	AnchorID     string
	ThreadID     string
	StartLine    int
	EndLine      int
	QuotedText   string
	AnchorStatus string
	ReviewState  string
}

// RemapAnchors relocates anchors from oldMarkdown to newMarkdown using quoted_text matching.
func RemapAnchors(oldMarkdown, newMarkdown string, anchors []AnchorInput) []RemappedAnchor {
	out := make([]RemappedAnchor, 0, len(anchors))
	for _, anchor := range anchors {
		out = append(out, remapOne(oldMarkdown, newMarkdown, anchor))
	}
	return out
}

func remapOne(oldMarkdown, newMarkdown string, anchor AnchorInput) RemappedAnchor {
	quoted := strings.TrimSpace(anchor.QuotedText)
	base := RemappedAnchor{
		AnchorID:    anchor.AnchorID,
		ThreadID:    anchor.ThreadID,
		StartLine:   anchor.StartLine,
		EndLine:     anchor.EndLine,
		QuotedText:  anchor.QuotedText,
		ReviewState: anchor.ReviewState,
	}

	if quoted == "" {
		base.AnchorStatus = AnchorStatusOrphaned
		return base
	}

	start, end, ok := findQuotedTextLines(newMarkdown, quoted)
	if !ok {
		base.AnchorStatus = AnchorStatusOrphaned
		return base
	}

	status := AnchorStatusActive
	if start != anchor.StartLine || end != anchor.EndLine {
		status = AnchorStatusShifted
	}

	base.StartLine = start
	base.EndLine = end
	base.QuotedText = quoted
	base.AnchorStatus = status
	return base
}

func findQuotedTextLines(markdown, quoted string) (startLine, endLine int, ok bool) {
	idx := strings.Index(markdown, quoted)
	if idx < 0 {
		return 0, 0, false
	}
	before := markdown[:idx]
	startLine = strings.Count(before, "\n") + 1
	endLine = startLine + strings.Count(quoted, "\n")
	return startLine, endLine, true
}
