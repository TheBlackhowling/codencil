package store_test

import (
	"context"
	"testing"

	"github.com/TheBlackHowling/codencil/apps/api/internal/store"
)

func TestGetOrCreateUser(t *testing.T) {
	s := testStore(t)
	ctx := context.Background()

	first, err := s.GetOrCreateUser(ctx, store.GetOrCreateUserInput{
		OrgID:      "org-auth",
		ExternalID: "alice",
	})
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	if first.ID == "" {
		t.Fatal("expected user id")
	}

	second, err := s.GetOrCreateUser(ctx, store.GetOrCreateUserInput{
		OrgID:      "org-auth",
		ExternalID: "alice",
	})
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	if second.ID != first.ID {
		t.Fatalf("expected same user id")
	}
}
