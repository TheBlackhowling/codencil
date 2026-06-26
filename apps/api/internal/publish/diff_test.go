package publish

import "testing"

func TestDiffLines_unchanged(t *testing.T) {
	changes := DiffLines("a\nb", "a\nb")
	if len(changes) != 2 {
		t.Fatalf("len = %d", len(changes))
	}
	for _, c := range changes {
		if c.Kind != LineUnchanged {
			t.Fatalf("unexpected kind %v", c.Kind)
		}
	}
}

func TestDiffLines_addedLine(t *testing.T) {
	changes := DiffLines("a", "a\nb")
	var added int
	for _, c := range changes {
		if c.Kind == LineAdded {
			added++
		}
	}
	if added != 1 {
		t.Fatalf("added = %d", added)
	}
}

func TestDiffLines_removedLine(t *testing.T) {
	changes := DiffLines("a\nb", "a")
	var removed int
	for _, c := range changes {
		if c.Kind == LineRemoved {
			removed++
		}
	}
	if removed != 1 {
		t.Fatalf("removed = %d", removed)
	}
}
