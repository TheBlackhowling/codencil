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
	ownerID := testOwner(t, s, ctx, "owner-review")

	doc, err := s.CreateDocument(ctx, store.CreateDocumentInput{
		OrgID:         "org-review",
		Title:         "Review doc",
		DraftMarkdown: "# Title\n\nQuoted text here.",
		OwnerUserID:   ownerID,
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

func TestPublishRemapsAnchorsToV2(t *testing.T) {
	s := testStore(t)
	ctx := context.Background()
	ownerID := testOwner(t, s, ctx, "owner-remap")

	doc, err := s.CreateDocument(ctx, store.CreateDocumentInput{
		OrgID:         "org-remap",
		Title:         "Remap",
		DraftMarkdown: "line one\nanchor text\nline three",
		OwnerUserID:   ownerID,
	})
	if err != nil {
		t.Fatalf("create: %v", err)
	}

	if _, err := s.PublishDocument(ctx, doc.ID, "author"); err != nil {
		t.Fatalf("publish v1: %v", err)
	}

	_, err = s.CreateAnchorWithThread(ctx, store.CreateAnchorInput{
		DocumentID: doc.ID,
		Version:    1,
		StartLine:  2,
		EndLine:    2,
		QuotedText: "anchor text",
		AuthorID:   "reviewer",
		Body:       "comment",
	})
	if err != nil {
		t.Fatalf("anchor: %v", err)
	}

	newDraft := "line one\ninserted\nanchor text\nline three"
	if _, err := s.UpdateDocumentDraft(ctx, doc.ID, store.UpdateDocumentDraftInput{
		DraftMarkdown: &newDraft,
	}); err != nil {
		t.Fatalf("update draft: %v", err)
	}

	if _, err := s.PublishDocument(ctx, doc.ID, "author"); err != nil {
		t.Fatalf("publish v2: %v", err)
	}

	v2Anchors, err := s.ListAnchorsByVersion(ctx, doc.ID, 2)
	if err != nil {
		t.Fatalf("list v2 anchors: %v", err)
	}
	if len(v2Anchors) != 1 {
		t.Fatalf("v2 anchor count = %d", len(v2Anchors))
	}
	anchor := v2Anchors[0].Anchor
	if anchor.AnchorStatus != models.AnchorStatusShifted {
		t.Fatalf("status = %q", anchor.AnchorStatus)
	}
	if anchor.StartLine != 3 || anchor.EndLine != 3 {
		t.Fatalf("lines = %d-%d", anchor.StartLine, anchor.EndLine)
	}
}
