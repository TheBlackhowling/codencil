package store

import (
	"context"
	"fmt"

	"github.com/TheBlackHowling/codencil/apps/api/internal/models"
	"github.com/TheBlackHowling/typrow"
)

var ErrForbidden = fmt.Errorf("forbidden")

func (s *Store) AddDocumentMember(ctx context.Context, documentID, userID, role string) error {
	_, err := typrow.QueryFirst[*models.DocumentMember](ctx, s.db, `
		INSERT INTO document_members (document_id, user_id, role)
		VALUES ($1, $2, $3)
		ON CONFLICT (document_id, user_id) DO UPDATE SET role = EXCLUDED.role
		RETURNING document_id, user_id, role, created_at`,
		documentID, userID, role,
	)
	return err
}

func (s *Store) GetDocumentRole(ctx context.Context, documentID, userID string) (string, error) {
	member := &models.DocumentMember{DocumentID: documentID, UserID: userID}
	if err := typrow.LoadByComposite(ctx, s.db, member, "document_member"); err != nil {
		return "", err
	}
	return member.Role, nil
}

func (s *Store) RequireDocumentRole(ctx context.Context, documentID, userID string, allowed ...string) error {
	role, err := s.GetDocumentRole(ctx, documentID, userID)
	if err != nil {
		return ErrForbidden
	}
	for _, a := range allowed {
		if role == a {
			return nil
		}
	}
	return ErrForbidden
}

func (s *Store) DocumentIDForThread(ctx context.Context, threadID string) (string, error) {
	row, err := typrow.QueryFirst[*threadDocumentRow](ctx, s.db, `
		SELECT document_id FROM comment_threads WHERE id = $1`, threadID)
	if err != nil {
		return "", err
	}
	return row.DocumentID, nil
}

func (s *Store) DocumentIDForAnchor(ctx context.Context, anchorRowID string) (string, error) {
	row, err := typrow.QueryFirst[*anchorDocumentRow](ctx, s.db, `
		SELECT document_id FROM version_anchors WHERE id = $1`, anchorRowID)
	if err != nil {
		return "", err
	}
	return row.DocumentID, nil
}

type threadDocumentRow struct {
	typrow.Model
	DocumentID string `db:"document_id"`
}

type anchorDocumentRow struct {
	typrow.Model
	DocumentID string `db:"document_id"`
}
