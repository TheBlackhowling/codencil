package store

import (
	"context"
	"fmt"
	"time"

	"github.com/TheBlackHowling/codencil/apps/api/internal/models"
	"github.com/TheBlackHowling/codencil/apps/api/internal/publish"
	"github.com/TheBlackHowling/typrow"
)

// AnchorWithThread is a version anchor with its thread comments.
type AnchorWithThread struct {
	Anchor   *models.VersionAnchor
	Comments []*models.Comment
}

type CreateAnchorInput struct {
	DocumentID string
	Version    int
	StartLine  int
	EndLine    int
	QuotedText string
	AuthorID   string
	Body       string
}

func (s *Store) CreateAnchorWithThread(ctx context.Context, in CreateAnchorInput) (*AnchorWithThread, error) {
	var result *AnchorWithThread
	err := s.db.WithTx(ctx, func(tx *typrow.Tx) error {
		v := &models.DocumentVersion{DocumentID: in.DocumentID, Version: in.Version}
		if err := typrow.LoadByComposite(ctx, tx, v, "document_version"); err != nil {
			return fmt.Errorf("version not found: %w", err)
		}

		thread := &models.CommentThread{DocumentID: in.DocumentID}
		thread, err := typrow.InsertAndLoad[*models.CommentThread](ctx, tx, thread)
		if err != nil {
			return err
		}

		anchorRow, err := typrow.QueryFirst[*models.VersionAnchor](ctx, tx, `
			INSERT INTO version_anchors (
				anchor_id, document_id, version, thread_id,
				start_line, end_line, quoted_text,
				anchor_status, review_state
			)
			VALUES (gen_random_uuid(), $1, $2, $3, $4, $5, $6, $7, $8)
			RETURNING id, anchor_id, document_id, version, thread_id, start_line, end_line,
				quoted_text, anchor_status, review_state, resolved_by, resolved_at, created_at`,
			in.DocumentID, in.Version, thread.ID,
			in.StartLine, in.EndLine, in.QuotedText,
			models.AnchorStatusActive, models.ReviewStateOpen,
		)
		if err != nil {
			return err
		}

		comment, err := typrow.InsertAndLoad[*models.Comment](ctx, tx, &models.Comment{
			ThreadID: thread.ID,
			AuthorID: in.AuthorID,
			Body:     in.Body,
		})
		if err != nil {
			return err
		}

		result = &AnchorWithThread{Anchor: anchorRow, Comments: []*models.Comment{comment}}
		return nil
	}, nil)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (s *Store) ListAnchorsByVersion(ctx context.Context, documentID string, version int) ([]AnchorWithThread, error) {
	anchors, err := typrow.QueryAll[*models.VersionAnchor](ctx, s.db, `
		SELECT id, anchor_id, document_id, version, thread_id, start_line, end_line,
			quoted_text, anchor_status, review_state, resolved_by, resolved_at, created_at
		FROM version_anchors
		WHERE document_id = $1 AND version = $2
		ORDER BY start_line ASC, created_at ASC`,
		documentID, version,
	)
	if err != nil {
		return nil, err
	}

	out := make([]AnchorWithThread, 0, len(anchors))
	for _, anchor := range anchors {
		comments, err := typrow.QueryAll[*models.Comment](ctx, s.db, `
			SELECT id, thread_id, author_id, body, created_at
			FROM comments
			WHERE thread_id = $1
			ORDER BY created_at ASC`,
			anchor.ThreadID,
		)
		if err != nil {
			return nil, err
		}
		out = append(out, AnchorWithThread{Anchor: anchor, Comments: comments})
	}
	return out, nil
}

func (s *Store) GetAnchor(ctx context.Context, anchorRowID string) (*models.VersionAnchor, error) {
	anchor := &models.VersionAnchor{ID: anchorRowID}
	if err := typrow.Load(ctx, s.db, anchor); err != nil {
		return nil, err
	}
	return anchor, nil
}

func (s *Store) AddComment(ctx context.Context, threadID, authorID, body string) (*models.Comment, error) {
	thread := &models.CommentThread{ID: threadID}
	if err := typrow.Load(ctx, s.db, thread); err != nil {
		return nil, fmt.Errorf("thread not found: %w", err)
	}
	return typrow.InsertAndLoad[*models.Comment](ctx, s.db, &models.Comment{
		ThreadID: threadID,
		AuthorID: authorID,
		Body:     body,
	})
}

func (s *Store) ResolveAnchor(ctx context.Context, anchorRowID, resolvedBy string) (*models.VersionAnchor, error) {
	anchor, err := s.GetAnchor(ctx, anchorRowID)
	if err != nil {
		return nil, err
	}
	now := time.Now().UTC()
	anchor.ReviewState = models.ReviewStateResolved
	anchor.ResolvedBy = &resolvedBy
	anchor.ResolvedAt = &now
	if err := typrow.Update(ctx, s.db, anchor); err != nil {
		return nil, err
	}
	return s.GetAnchor(ctx, anchorRowID)
}

func (s *Store) ReopenAnchor(ctx context.Context, anchorRowID string) (*models.VersionAnchor, error) {
	anchor, err := s.GetAnchor(ctx, anchorRowID)
	if err != nil {
		return nil, err
	}
	anchor.ReviewState = models.ReviewStateOpen
	anchor.ResolvedBy = nil
	anchor.ResolvedAt = nil
	if err := typrow.Update(ctx, s.db, anchor); err != nil {
		return nil, err
	}
	return s.GetAnchor(ctx, anchorRowID)
}

func (s *Store) remapAnchorsToVersion(ctx context.Context, tx *typrow.Tx, previous, next *models.DocumentVersion) error {
	prevAnchors, err := typrow.QueryAll[*models.VersionAnchor](ctx, tx, `
		SELECT id, anchor_id, document_id, version, thread_id, start_line, end_line,
			quoted_text, anchor_status, review_state, resolved_by, resolved_at, created_at
		FROM version_anchors
		WHERE document_id = $1 AND version = $2`,
		previous.DocumentID, previous.Version,
	)
	if err != nil {
		return err
	}
	if len(prevAnchors) == 0 {
		return nil
	}

	inputs := make([]publish.AnchorInput, 0, len(prevAnchors))
	for _, anchor := range prevAnchors {
		inputs = append(inputs, publish.AnchorInput{
			AnchorID:    anchor.AnchorID,
			ThreadID:    anchor.ThreadID,
			StartLine:   anchor.StartLine,
			EndLine:     anchor.EndLine,
			QuotedText:  anchor.QuotedText,
			ReviewState: anchor.ReviewState,
		})
	}

	remapped := publish.RemapAnchors(previous.Markdown, next.Markdown, inputs)
	for i, item := range remapped {
		source := prevAnchors[i]
		_, err := typrow.QueryFirst[*models.VersionAnchor](ctx, tx, `
			INSERT INTO version_anchors (
				anchor_id, document_id, version, thread_id,
				start_line, end_line, quoted_text,
				anchor_status, review_state, resolved_by, resolved_at
			)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
			RETURNING id`,
			item.AnchorID, next.DocumentID, next.Version, item.ThreadID,
			item.StartLine, item.EndLine, item.QuotedText,
			item.AnchorStatus, item.ReviewState, source.ResolvedBy, source.ResolvedAt,
		)
		if err != nil {
			return err
		}
	}
	return nil
}
