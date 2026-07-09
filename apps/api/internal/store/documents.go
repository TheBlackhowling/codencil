package store

import (
	"context"
	"fmt"

	"github.com/TheBlackHowling/codencil/apps/api/internal/models"
	"github.com/TheBlackHowling/typrow"
)

// Store provides document persistence via TypRow.
type Store struct {
	db *typrow.DB
}

type versionNextRow struct {
	typrow.Model
	NextVersion int `db:"next_version"`
}

func New(db *typrow.DB) *Store {
	return &Store{db: db}
}

type CreateDocumentInput struct {
	OrgID         string
	Title         string
	DraftMarkdown string
	OwnerUserID   string
}

func (s *Store) CreateDocument(ctx context.Context, in CreateDocumentInput) (*models.Document, error) {
	var doc *models.Document
	err := s.db.WithTx(ctx, func(tx *typrow.Tx) error {
		created, err := typrow.InsertAndLoad[*models.Document](ctx, tx, &models.Document{
			OrgID:         in.OrgID,
			Title:         in.Title,
			DraftMarkdown: in.DraftMarkdown,
		})
		if err != nil {
			return err
		}
		if in.OwnerUserID != "" {
			_, err = typrow.QueryFirst[*models.DocumentMember](ctx, tx, `
				INSERT INTO document_members (document_id, user_id, role)
				VALUES ($1, $2, $3)
				RETURNING document_id, user_id, role, created_at`,
				created.ID, in.OwnerUserID, models.RoleOwner,
			)
			if err != nil {
				return err
			}
		}
		doc = created
		return nil
	}, nil)
	if err != nil {
		return nil, err
	}
	return doc, nil
}

func (s *Store) GetDocument(ctx context.Context, id string) (*models.Document, error) {
	doc := &models.Document{ID: id}
	if err := typrow.Load(ctx, s.db, doc); err != nil {
		return nil, err
	}
	return doc, nil
}

type UpdateDocumentDraftInput struct {
	Title         *string
	DraftMarkdown *string
}

func (s *Store) UpdateDocumentDraft(ctx context.Context, id string, in UpdateDocumentDraftInput) (*models.Document, error) {
	doc, err := s.GetDocument(ctx, id)
	if err != nil {
		return nil, err
	}
	if in.Title != nil {
		doc.Title = *in.Title
	}
	if in.DraftMarkdown != nil {
		doc.DraftMarkdown = *in.DraftMarkdown
	}
	if err := typrow.Update(ctx, s.db, doc); err != nil {
		return nil, err
	}
	return s.GetDocument(ctx, id)
}

func (s *Store) PublishDocument(ctx context.Context, documentID, publishedBy string) (*models.DocumentVersion, error) {
	var published *models.DocumentVersion
	err := s.db.WithTx(ctx, func(tx *typrow.Tx) error {
		doc := &models.Document{ID: documentID}
		if err := typrow.Load(ctx, tx, doc); err != nil {
			return err
		}

		var previous *models.DocumentVersion
		prevRow, err := typrow.QueryFirst[*versionNextRow](ctx, tx,
			`SELECT COALESCE(MAX(version), 0) + 1 AS next_version FROM document_versions WHERE document_id = $1`,
			documentID,
		)
		if err != nil {
			return fmt.Errorf("next version: %w", err)
		}
		nextVersion := prevRow.NextVersion
		if nextVersion > 1 {
			previous = &models.DocumentVersion{DocumentID: documentID, Version: nextVersion - 1}
			if err := typrow.LoadByComposite(ctx, tx, previous, "document_version"); err != nil {
				return fmt.Errorf("load previous version: %w", err)
			}
		}

		inserted, err := typrow.QueryFirst[*models.DocumentVersion](ctx, tx, `
			INSERT INTO document_versions (document_id, version, markdown, published_by)
			VALUES ($1, $2, $3, $4)
			RETURNING document_id, version, markdown, published_at, published_by`,
			documentID, nextVersion, doc.DraftMarkdown, publishedBy,
		)
		if err != nil {
			return err
		}
		published = inserted

		if previous != nil {
			if err := s.remapAnchorsToVersion(ctx, tx, previous, inserted); err != nil {
				return err
			}
		}
		return nil
	}, nil)
	if err != nil {
		return nil, err
	}
	return published, nil
}

func (s *Store) ListDocumentVersions(ctx context.Context, documentID string) ([]*models.DocumentVersion, error) {
	doc := &models.Document{ID: documentID}
	if err := typrow.Load(ctx, s.db, doc); err != nil {
		return nil, err
	}
	return typrow.QueryAll[*models.DocumentVersion](ctx, s.db, `
		SELECT document_id, version, markdown, published_at, published_by
		FROM document_versions
		WHERE document_id = $1
		ORDER BY version DESC`,
		documentID,
	)
}

func (s *Store) GetDocumentVersion(ctx context.Context, documentID string, version int) (*models.DocumentVersion, error) {
	v := &models.DocumentVersion{DocumentID: documentID, Version: version}
	if err := typrow.LoadByComposite(ctx, s.db, v, "document_version"); err != nil {
		return nil, err
	}
	return v, nil
}
