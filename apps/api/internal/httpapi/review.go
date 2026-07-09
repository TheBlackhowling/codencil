package httpapi

import (
	"net/http"
	"strconv"

	"github.com/TheBlackHowling/codencil/apps/api/internal/models"
	"github.com/TheBlackHowling/codencil/apps/api/internal/store"
	"github.com/go-chi/chi/v5"
)

// ReviewHandler serves anchor and comment endpoints.
type ReviewHandler struct {
	store *store.Store
}

func NewReviewHandler(s *store.Store) *ReviewHandler {
	return &ReviewHandler{store: s}
}

func (h *ReviewHandler) Register(r chi.Router) {
	r.Get("/documents/{id}/versions/{version}/anchors", h.listAnchors)
	r.Post("/documents/{id}/versions/{version}/anchors", h.createAnchor)
	r.Post("/threads/{id}/comments", h.addComment)
	r.Post("/anchors/{id}/resolve", h.resolveAnchor)
	r.Post("/anchors/{id}/reopen", h.reopenAnchor)
}

type commentResponse struct {
	ID        string `json:"id"`
	ThreadID  string `json:"thread_id"`
	AuthorID  string `json:"author_id"`
	Body      string `json:"body"`
	CreatedAt string `json:"created_at"`
}

func toCommentResponse(c *models.Comment) commentResponse {
	return commentResponse{
		ID:        c.ID,
		ThreadID:  c.ThreadID,
		AuthorID:  c.AuthorID,
		Body:      c.Body,
		CreatedAt: c.CreatedAt.UTC().Format(timeRFC3339),
	}
}

type anchorResponse struct {
	ID           string  `json:"id"`
	AnchorID     string  `json:"anchor_id"`
	DocumentID   string  `json:"document_id"`
	Version      int     `json:"version"`
	ThreadID     string  `json:"thread_id"`
	StartLine    int     `json:"start_line"`
	EndLine      int     `json:"end_line"`
	QuotedText   string  `json:"quoted_text"`
	AnchorStatus string  `json:"anchor_status"`
	ReviewState  string  `json:"review_state"`
	ResolvedBy   *string `json:"resolved_by,omitempty"`
	ResolvedAt   *string `json:"resolved_at,omitempty"`
	CreatedAt    string  `json:"created_at"`
	Comments     []commentResponse `json:"comments"`
}

func toAnchorResponse(item store.AnchorWithThread) anchorResponse {
	var resolvedAt *string
	if item.Anchor.ResolvedAt != nil {
		s := item.Anchor.ResolvedAt.UTC().Format(timeRFC3339)
		resolvedAt = &s
	}
	comments := make([]commentResponse, 0, len(item.Comments))
	for _, c := range item.Comments {
		comments = append(comments, toCommentResponse(c))
	}
	return anchorResponse{
		ID:           item.Anchor.ID,
		AnchorID:     item.Anchor.AnchorID,
		DocumentID:   item.Anchor.DocumentID,
		Version:      item.Anchor.Version,
		ThreadID:     item.Anchor.ThreadID,
		StartLine:    item.Anchor.StartLine,
		EndLine:      item.Anchor.EndLine,
		QuotedText:   item.Anchor.QuotedText,
		AnchorStatus: item.Anchor.AnchorStatus,
		ReviewState:  item.Anchor.ReviewState,
		ResolvedBy:   item.Anchor.ResolvedBy,
		ResolvedAt:   resolvedAt,
		CreatedAt:    item.Anchor.CreatedAt.UTC().Format(timeRFC3339),
		Comments:     comments,
	}
}

func (h *ReviewHandler) listAnchors(w http.ResponseWriter, r *http.Request) {
	documentID := chi.URLParam(r, "id")
	version, err := parseVersionParam(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid version")
		return
	}

	items, err := h.store.ListAnchorsByVersion(r.Context(), documentID, version)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "list anchors failed")
		return
	}

	out := make([]anchorResponse, 0, len(items))
	for _, item := range items {
		out = append(out, toAnchorResponse(item))
	}
	writeJSON(w, http.StatusOK, out)
}

type createAnchorRequest struct {
	StartLine  int    `json:"start_line"`
	EndLine    int    `json:"end_line"`
	QuotedText string `json:"quoted_text"`
	AuthorID   string `json:"author_id"`
	Body       string `json:"body"`
}

func (h *ReviewHandler) createAnchor(w http.ResponseWriter, r *http.Request) {
	documentID := chi.URLParam(r, "id")
	version, err := parseVersionParam(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid version")
		return
	}

	var req createAnchorRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}
	if req.StartLine < 1 || req.EndLine < req.StartLine || req.QuotedText == "" || req.Body == "" {
		writeError(w, http.StatusBadRequest, "start_line, end_line, quoted_text, and body are required")
		return
	}

	authorID := userExternalID(r)

	created, err := h.store.CreateAnchorWithThread(r.Context(), store.CreateAnchorInput{
		DocumentID: documentID,
		Version:    version,
		StartLine:  req.StartLine,
		EndLine:    req.EndLine,
		QuotedText: req.QuotedText,
		AuthorID:   authorID,
		Body:       req.Body,
	})
	if err != nil {
		writeError(w, http.StatusNotFound, "create anchor failed")
		return
	}
	writeJSON(w, http.StatusCreated, toAnchorResponse(*created))
}

type addCommentRequest struct {
	AuthorID string `json:"author_id"`
	Body     string `json:"body"`
}

func (h *ReviewHandler) addComment(w http.ResponseWriter, r *http.Request) {
	threadID := chi.URLParam(r, "id")
	var req addCommentRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}
	if req.Body == "" {
		writeError(w, http.StatusBadRequest, "body is required")
		return
	}

	authorID := userExternalID(r)

	comment, err := h.store.AddComment(r.Context(), threadID, authorID, req.Body)
	if err != nil {
		writeError(w, http.StatusNotFound, "add comment failed")
		return
	}
	writeJSON(w, http.StatusCreated, toCommentResponse(comment))
}

type resolveAnchorRequest struct {
	ResolvedBy string `json:"resolved_by"`
}

func (h *ReviewHandler) resolveAnchor(w http.ResponseWriter, r *http.Request) {
	anchorID := chi.URLParam(r, "id")
	resolvedBy := userExternalID(r)

	anchor, err := h.store.ResolveAnchor(r.Context(), anchorID, resolvedBy)
	if err != nil {
		writeError(w, http.StatusNotFound, "resolve anchor failed")
		return
	}
	writeJSON(w, http.StatusOK, anchorResponse{
		ID:           anchor.ID,
		AnchorID:     anchor.AnchorID,
		DocumentID:   anchor.DocumentID,
		Version:      anchor.Version,
		ThreadID:     anchor.ThreadID,
		StartLine:    anchor.StartLine,
		EndLine:      anchor.EndLine,
		QuotedText:   anchor.QuotedText,
		AnchorStatus: anchor.AnchorStatus,
		ReviewState:  anchor.ReviewState,
		ResolvedBy:   anchor.ResolvedBy,
		CreatedAt:    anchor.CreatedAt.UTC().Format(timeRFC3339),
	})
}

func (h *ReviewHandler) reopenAnchor(w http.ResponseWriter, r *http.Request) {
	anchorID := chi.URLParam(r, "id")
	anchor, err := h.store.ReopenAnchor(r.Context(), anchorID)
	if err != nil {
		writeError(w, http.StatusNotFound, "reopen anchor failed")
		return
	}
	writeJSON(w, http.StatusOK, anchorResponse{
		ID:           anchor.ID,
		AnchorID:     anchor.AnchorID,
		DocumentID:   anchor.DocumentID,
		Version:      anchor.Version,
		ThreadID:     anchor.ThreadID,
		StartLine:    anchor.StartLine,
		EndLine:      anchor.EndLine,
		QuotedText:   anchor.QuotedText,
		AnchorStatus: anchor.AnchorStatus,
		ReviewState:  anchor.ReviewState,
		CreatedAt:    anchor.CreatedAt.UTC().Format(timeRFC3339),
	})
}

func parseVersionParam(r *http.Request) (int, error) {
	versionStr := chi.URLParam(r, "version")
	version, err := strconv.Atoi(versionStr)
	if err != nil || version < 1 {
		return 0, err
	}
	return version, nil
}
