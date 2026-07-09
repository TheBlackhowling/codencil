package store

import (
	"context"

	"github.com/TheBlackHowling/codencil/apps/api/internal/models"
	"github.com/TheBlackHowling/typrow"
)

type GetOrCreateUserInput struct {
	OrgID       string
	ExternalID  string
	DisplayName string
}

func (s *Store) GetOrCreateUser(ctx context.Context, in GetOrCreateUserInput) (*models.User, error) {
	existing, err := typrow.QueryFirst[*models.User](ctx, s.db, `
		SELECT id, org_id, external_id, display_name, created_at
		FROM users WHERE org_id = $1 AND external_id = $2`,
		in.OrgID, in.ExternalID,
	)
	if err == nil && existing != nil {
		return existing, nil
	}

	displayName := in.DisplayName
	if displayName == "" {
		displayName = in.ExternalID
	}
	return typrow.InsertAndLoad[*models.User](ctx, s.db, &models.User{
		OrgID:       in.OrgID,
		ExternalID:  in.ExternalID,
		DisplayName: displayName,
	})
}
