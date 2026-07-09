package httpapi

import (
	"net/http"
	"strconv"

	"github.com/TheBlackHowling/codencil/apps/api/internal/models"
	"github.com/TheBlackHowling/codencil/apps/api/internal/store"
	"github.com/go-chi/chi/v5"
)

const defaultOrgID = "dev"

// DocumentHandler serves document CRUD and publish endpoints.
type DocumentHandler struct {
	store *store.Store
}

func NewDocumentHandler(s *store.Store) *DocumentHandler {
	return &DocumentHandler{store: s}
}

func (h *DocumentHandler) Register(r chi.Router) {
	r.Post("/documents", h.createDocument)
	r.Get("/documents/{id}", h.getDocument)
	r.Patch("/documents/{id}", h.patchDocument)
	r.Post("/documents/{id}/publish", h.publishDocument)
	r.Get("/documents/{id}/versions", h.listDocumentVersions)
	r.Get("/documents/{id}/versions/{version}", h.getDocumentVersion)
}

type createDocumentRequest struct {
	Title         string `json:"title"`
	DraftMarkdown string `json:"draft_markdown"`
}

type documentResponse struct {
	ID            string `json:"id"`
	OrgID         string `json:"org_id"`
	Title         string `json:"title"`
	DraftMarkdown string `json:"draft_markdown"`
	CreatedAt     string `json:"created_at"`
	UpdatedAt     string `json:"updated_at"`
}

func toDocumentResponse(d *models.Document) documentResponse {
	return documentResponse{
		ID:            d.ID,
		OrgID:         d.OrgID,
		Title:         d.Title,
		DraftMarkdown: d.DraftMarkdown,
		CreatedAt:     d.CreatedAt.UTC().Format(timeRFC3339),
		UpdatedAt:     d.UpdatedAt.UTC().Format(timeRFC3339),
	}
}

const timeRFC3339 = "2006-01-02T15:04:05Z07:00"

func (h *DocumentHandler) createDocument(w http.ResponseWriter, r *http.Request) {
	user, ok := requireUser(w, r)
	if !ok {
		return
	}
	var req createDocumentRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}

	doc, err := h.store.CreateDocument(r.Context(), store.CreateDocumentInput{
		OrgID:         defaultOrgID,
		Title:         req.Title,
		DraftMarkdown: req.DraftMarkdown,
		OwnerUserID:   user.ID,
	})
	if err != nil {
		writeError(w, http.StatusInternalServerError, "create document failed")
		return
	}
	writeJSON(w, http.StatusCreated, toDocumentResponse(doc))
}

func (h *DocumentHandler) getDocument(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if !h.checkDocumentRole(w, r, id, roleRead...) {
		return
	}
	doc, err := h.store.GetDocument(r.Context(), id)
	if err != nil {
		writeError(w, http.StatusNotFound, "document not found")
		return
	}
	writeJSON(w, http.StatusOK, toDocumentResponse(doc))
}

type patchDocumentRequest struct {
	Title         *string `json:"title"`
	DraftMarkdown *string `json:"draft_markdown"`
}

func (h *DocumentHandler) patchDocument(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if !h.checkDocumentRole(w, r, id, roleOwner...) {
		return
	}
	var req patchDocumentRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}

	doc, err := h.store.UpdateDocumentDraft(r.Context(), id, store.UpdateDocumentDraftInput{
		Title:         req.Title,
		DraftMarkdown: req.DraftMarkdown,
	})
	if err != nil {
		writeError(w, http.StatusNotFound, "document not found")
		return
	}
	writeJSON(w, http.StatusOK, toDocumentResponse(doc))
}

type publishDocumentRequest struct {
	PublishedBy string `json:"published_by"`
}

type documentVersionResponse struct {
	DocumentID  string `json:"document_id"`
	Version     int    `json:"version"`
	Markdown    string `json:"markdown"`
	PublishedAt string `json:"published_at"`
	PublishedBy string `json:"published_by"`
}

func toVersionResponse(v *models.DocumentVersion) documentVersionResponse {
	return documentVersionResponse{
		DocumentID:  v.DocumentID,
		Version:     v.Version,
		Markdown:    v.Markdown,
		PublishedAt: v.PublishedAt.UTC().Format(timeRFC3339),
		PublishedBy: v.PublishedBy,
	}
}

func (h *DocumentHandler) publishDocument(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if !h.checkDocumentRole(w, r, id, roleOwner...) {
		return
	}
	publishedBy := userExternalID(r)

	version, err := h.store.PublishDocument(r.Context(), id, publishedBy)
	if err != nil {
		writeError(w, http.StatusNotFound, "document not found")
		return
	}
	writeJSON(w, http.StatusCreated, toVersionResponse(version))
}

type versionSummaryResponse struct {
	Version     int    `json:"version"`
	PublishedAt string `json:"published_at"`
	PublishedBy string `json:"published_by"`
}

func (h *DocumentHandler) listDocumentVersions(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if !h.checkDocumentRole(w, r, id, roleRead...) {
		return
	}
	versions, err := h.store.ListDocumentVersions(r.Context(), id)
	if err != nil {
		writeError(w, http.StatusNotFound, "document not found")
		return
	}
	out := make([]versionSummaryResponse, 0, len(versions))
	for _, v := range versions {
		out = append(out, versionSummaryResponse{
			Version:     v.Version,
			PublishedAt: v.PublishedAt.UTC().Format(timeRFC3339),
			PublishedBy: v.PublishedBy,
		})
	}
	writeJSON(w, http.StatusOK, out)
}

func (h *DocumentHandler) getDocumentVersion(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if !h.checkDocumentRole(w, r, id, roleRead...) {
		return
	}
	versionStr := chi.URLParam(r, "version")
	version, err := strconv.Atoi(versionStr)
	if err != nil || version < 1 {
		writeError(w, http.StatusBadRequest, "invalid version")
		return
	}

	v, err := h.store.GetDocumentVersion(r.Context(), id, version)
	if err != nil {
		writeError(w, http.StatusNotFound, "version not found")
		return
	}
	writeJSON(w, http.StatusOK, toVersionResponse(v))
}
