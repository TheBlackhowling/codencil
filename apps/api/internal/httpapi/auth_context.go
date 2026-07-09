package httpapi

import (
	"net/http"

	"github.com/TheBlackHowling/codencil/apps/api/internal/auth"
	"github.com/TheBlackHowling/codencil/apps/api/internal/models"
)

func userFromRequest(r *http.Request) (*models.User, bool) {
	return auth.UserFromContext(r.Context())
}

func userExternalID(r *http.Request) string {
	if u, ok := userFromRequest(r); ok {
		return u.ExternalID
	}
	return auth.DefaultDevUserID
}

func requireUser(w http.ResponseWriter, r *http.Request) (*models.User, bool) {
	u, ok := userFromRequest(r)
	if !ok {
		writeError(w, http.StatusUnauthorized, "authentication required")
		return nil, false
	}
	return u, true
}
