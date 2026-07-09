package publish

import "testing"

func TestRemapAnchors_unchanged(t *testing.T) {
	old := "# Title\n\nKeep this line.\n"
	newMd := "# Title\n\nKeep this line.\n"
	anchors := []AnchorInput{{
		AnchorID: "a1", ThreadID: "t1",
		StartLine: 3, EndLine: 3,
		QuotedText: "Keep this line.",
		ReviewState: "open",
	}}

	got := RemapAnchors(old, newMd, anchors)
	if len(got) != 1 {
		t.Fatalf("len = %d", len(got))
	}
	if got[0].AnchorStatus != AnchorStatusActive {
		t.Fatalf("status = %q", got[0].AnchorStatus)
	}
	if got[0].StartLine != 3 || got[0].EndLine != 3 {
		t.Fatalf("lines = %d-%d", got[0].StartLine, got[0].EndLine)
	}
}

func TestRemapAnchors_shiftedWhenLineInsertedAbove(t *testing.T) {
	old := "a\nb\nc\n"
	newMd := "a\ninserted\nb\nc\n"
	anchors := []AnchorInput{{
		AnchorID: "a1", ThreadID: "t1",
		StartLine: 2, EndLine: 2,
		QuotedText: "b",
		ReviewState: "open",
	}}

	got := RemapAnchors(old, newMd, anchors)
	if got[0].AnchorStatus != AnchorStatusShifted {
		t.Fatalf("status = %q", got[0].AnchorStatus)
	}
	if got[0].StartLine != 3 || got[0].EndLine != 3 {
		t.Fatalf("lines = %d-%d", got[0].StartLine, got[0].EndLine)
	}
}

func TestRemapAnchors_orphanedWhenTextRemoved(t *testing.T) {
	old := "a\nremove me\nc\n"
	newMd := "a\nc\n"
	anchors := []AnchorInput{{
		AnchorID: "a1", ThreadID: "t1",
		StartLine: 2, EndLine: 2,
		QuotedText: "remove me",
		ReviewState: "resolved",
	}}

	got := RemapAnchors(old, newMd, anchors)
	if got[0].AnchorStatus != AnchorStatusOrphaned {
		t.Fatalf("status = %q", got[0].AnchorStatus)
	}
	if got[0].ReviewState != "resolved" {
		t.Fatalf("review state lost")
	}
}

func TestRemapAnchors_multilineQuote(t *testing.T) {
	old := "intro\nline one\nline two\noutro\n"
	newMd := "intro\nline one\nline two\noutro\n"
	quoted := "line one\nline two"
	anchors := []AnchorInput{{
		AnchorID: "a1", ThreadID: "t1",
		StartLine: 2, EndLine: 3,
		QuotedText: quoted,
		ReviewState: "open",
	}}

	got := RemapAnchors(old, newMd, anchors)
	if got[0].AnchorStatus != AnchorStatusActive {
		t.Fatalf("status = %q", got[0].AnchorStatus)
	}
	if got[0].StartLine != 2 || got[0].EndLine != 3 {
		t.Fatalf("lines = %d-%d", got[0].StartLine, got[0].EndLine)
	}
}
