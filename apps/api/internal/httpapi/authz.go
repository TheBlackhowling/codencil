package httpapi

import (
	"errors"
	"net/http"

	"github.com/TheBlackHowling/codencil/apps/api/internal/models"
	"github.com/TheBlackHowling/codencil/apps/api/internal/store"
	"github.com/go-chi/chi/v5"
)

var (
	roleRead     = []string{models.RoleOwner, models.RoleReviewer, models.RoleViewer}
	roleReview   = []string{models.RoleOwner, models.RoleReviewer}
	roleOwner    = []string{models.RoleOwner}
)

func (h *DocumentHandler) checkDocumentRole(w http.ResponseWriter, r *http.Request, documentID string, allowed ...string) bool {
	user, ok := requireUser(w, r)
	if !ok {
		return false
	}
	if err := h.store.RequireDocumentRole(r.Context(), documentID, user.ID, allowed...); err != nil {
		writeError(w, http.StatusForbidden, "forbidden")
		return false
	}
	return true
}

func (h *ReviewHandler) checkDocumentRole(w http.ResponseWriter, r *http.Request, documentID string, allowed ...string) bool {
	user, ok := requireUser(w, r)
	if !ok {
		return false
	}
	if err := h.store.RequireDocumentRole(r.Context(), documentID, user.ID, allowed...); err != nil {
		writeError(w, http.StatusForbidden, "forbidden")
		return false
	}
	return true
}

func (h *ReviewHandler) checkThreadRole(w http.ResponseWriter, r *http.Request, threadID string, allowed ...string) bool {
	docID, err := h.store.DocumentIDForThread(r.Context(), threadID)
	if err != nil {
		writeError(w, http.StatusNotFound, "thread not found")
		return false
	}
	return h.checkDocumentRole(w, r, docID, allowed...)
}

func (h *ReviewHandler) checkAnchorRole(w http.ResponseWriter, r *http.Request, anchorID string, allowed ...string) bool {
	docID, err := h.store.DocumentIDForAnchor(r.Context(), anchorID)
	if err != nil {
		writeError(w, http.StatusNotFound, "anchor not found")
		return false
	}
	return h.checkDocumentRole(w, r, docID, allowed...)
}

func documentIDFromRequest(r *http.Request) string {
	return chi.URLParam(r, "id")
}

func isForbidden(err error) bool {
	return errors.Is(err, store.ErrForbidden)
}
