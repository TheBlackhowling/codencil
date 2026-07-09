package store_test

import (
	"context"
	"testing"

	"github.com/TheBlackHowling/codencil/apps/api/internal/models"
	"github.com/TheBlackHowling/codencil/apps/api/internal/store"
)

func TestDocumentRoleEnforcement(t *testing.T) {
	s := testStore(t)
	ctx := context.Background()

	ownerID := testOwner(t, s, ctx, "owner-role")
	viewerID := testOwner(t, s, ctx, "viewer-role")

	doc, err := s.CreateDocument(ctx, store.CreateDocumentInput{
		OrgID:         "org-role",
		Title:         "Roles",
		DraftMarkdown: "# doc",
		OwnerUserID:   ownerID,
	})
	if err != nil {
		t.Fatalf("create: %v", err)
	}

	if err := s.AddDocumentMember(ctx, doc.ID, viewerID, models.RoleViewer); err != nil {
		t.Fatalf("add viewer: %v", err)
	}

	if err := s.RequireDocumentRole(ctx, doc.ID, viewerID, models.RoleOwner); err == nil {
		t.Fatal("viewer should not have owner role")
	}
	if err := s.RequireDocumentRole(ctx, doc.ID, viewerID, models.RoleViewer); err != nil {
		t.Fatalf("viewer read: %v", err)
	}
}
