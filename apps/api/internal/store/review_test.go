package store_test

import (
	"context"
	"testing"

	"github.com/TheBlackHowling/codencil/apps/api/internal/models"
	"github.com/TheBlackHowling/codencil/apps/api/internal/store"
)

func TestReviewAnchorFlow(t *testing.T) {
	s := testStore(t)
	ctx := context.Background()

	doc, err := s.CreateDocument(ctx, store.CreateDocumentInput{
		OrgID:         "org-review",
		Title:         "Review doc",
		DraftMarkdown: "# Title\n\nQuoted text here.",
	})
	if err != nil {
		t.Fatalf("create doc: %v", err)
	}

	if _, err := s.PublishDocument(ctx, doc.ID, "author"); err != nil {
		t.Fatalf("publish: %v", err)
	}

	created, err := s.CreateAnchorWithThread(ctx, store.CreateAnchorInput{
		DocumentID: doc.ID,
		Version:    1,
		StartLine:  3,
		EndLine:    3,
		QuotedText: "Quoted text here.",
		AuthorID:   "reviewer-1",
		Body:       "Please clarify this.",
	})
	if err != nil {
		t.Fatalf("create anchor: %v", err)
	}
	if created.Anchor.ReviewState != models.ReviewStateOpen {
		t.Fatalf("review state = %q", created.Anchor.ReviewState)
	}
	if len(created.Comments) != 1 {
		t.Fatalf("comments = %d", len(created.Comments))
	}

	listed, err := s.ListAnchorsByVersion(ctx, doc.ID, 1)
	if err != nil {
		t.Fatalf("list anchors: %v", err)
	}
	if len(listed) != 1 {
		t.Fatalf("listed = %d", len(listed))
	}

	reply, err := s.AddComment(ctx, created.Anchor.ThreadID, "author", "Will update.")
	if err != nil {
		t.Fatalf("add comment: %v", err)
	}
	if reply.Body != "Will update." {
		t.Fatalf("reply body = %q", reply.Body)
	}

	resolved, err := s.ResolveAnchor(ctx, created.Anchor.ID, "author")
	if err != nil {
		t.Fatalf("resolve: %v", err)
	}
	if resolved.ReviewState != models.ReviewStateResolved {
		t.Fatalf("resolved state = %q", resolved.ReviewState)
	}

	reopened, err := s.ReopenAnchor(ctx, created.Anchor.ID)
	if err != nil {
		t.Fatalf("reopen: %v", err)
	}
	if reopened.ReviewState != models.ReviewStateOpen {
		t.Fatalf("reopened state = %q", reopened.ReviewState)
	}
}
