package store_test

import (
	"context"
	"os"
	"testing"

	"github.com/TheBlackHowling/codencil/apps/api/internal/db"
	"github.com/TheBlackHowling/codencil/apps/api/internal/store"
)

func testStore(t *testing.T) *store.Store {
	t.Helper()

	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		t.Skip("DATABASE_URL not set")
	}

	database, err := db.Open(dsn)
	if err != nil {
		t.Fatalf("open db: %v", err)
	}
	t.Cleanup(func() { _ = database.Close() })

	return store.New(database)
}

func TestDocumentCRUDAndPublish(t *testing.T) {
	s := testStore(t)
	ctx := context.Background()

	created, err := s.CreateDocument(ctx, store.CreateDocumentInput{
		OrgID:         "org-test",
		Title:         "Hello",
		DraftMarkdown: "# Draft v1",
	})
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	if created.ID == "" {
		t.Fatal("expected document id")
	}

	got, err := s.GetDocument(ctx, created.ID)
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	if got.Title != "Hello" {
		t.Fatalf("title = %q", got.Title)
	}

	newTitle := "Updated"
	newBody := "# Draft v2"
	updated, err := s.UpdateDocumentDraft(ctx, created.ID, store.UpdateDocumentDraftInput{
		Title:         &newTitle,
		DraftMarkdown: &newBody,
	})
	if err != nil {
		t.Fatalf("update: %v", err)
	}
	if updated.DraftMarkdown != newBody {
		t.Fatalf("draft = %q", updated.DraftMarkdown)
	}

	v1, err := s.PublishDocument(ctx, created.ID, "author-1")
	if err != nil {
		t.Fatalf("publish: %v", err)
	}
	if v1.Version != 1 || v1.Markdown != newBody {
		t.Fatalf("version: %+v", v1)
	}

	v1Again, err := s.GetDocumentVersion(ctx, created.ID, 1)
	if err != nil {
		t.Fatalf("get version: %v", err)
	}
	if v1Again.Markdown != newBody {
		t.Fatalf("snapshot markdown = %q", v1Again.Markdown)
	}
}
