package auth

import (
	"encoding/json"
	"net/http"

	"github.com/TheBlackHowling/codencil/apps/api/internal/models"
	"github.com/TheBlackHowling/codencil/apps/api/internal/store"
)

const devUserHeader = "X-Dev-User-Id"

// Middleware resolves the caller to a persisted user row.
type Middleware struct {
	store  *store.Store
	config Config
	oidc   *oidcVerifier
}

func NewMiddleware(s *store.Store, cfg Config) (*Middleware, error) {
	m := &Middleware{store: s, config: cfg}
	if cfg.Mode == ModeOIDC {
		v, err := newOIDCVerifier(cfg)
		if err != nil {
			return nil, err
		}
		m.oidc = v
	}
	return m, nil
}

func (m *Middleware) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, status, msg := m.authenticate(r)
		if status != 0 {
			writeAuthError(w, status, msg)
			return
		}
		next.ServeHTTP(w, r.WithContext(WithUser(r.Context(), user)))
	})
}

func (m *Middleware) authenticate(r *http.Request) (*models.User, int, string) {
	switch m.config.Mode {
	case ModeDev:
		return m.authenticateDev(r)
	case ModeOIDC:
		return m.authenticateOIDC(r)
	default:
		return nil, http.StatusInternalServerError, "invalid AUTH_MODE"
	}
}

func (m *Middleware) authenticateDev(r *http.Request) (*models.User, int, string) {
	externalID := r.Header.Get(devUserHeader)
	if externalID == "" {
		externalID = DefaultDevUserID
	}
	user, err := m.store.GetOrCreateUser(r.Context(), store.GetOrCreateUserInput{
		OrgID:      m.config.OrgID,
		ExternalID: externalID,
	})
	if err != nil {
		return nil, http.StatusInternalServerError, "auth failed"
	}
	return user, 0, ""
}

func writeAuthError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(map[string]string{"error": message})
}
