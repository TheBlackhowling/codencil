package auth

import (
	"fmt"
	"net/http"

	"github.com/TheBlackHowling/codencil/apps/api/internal/models"
)

type oidcVerifier struct {
	// populated in P4.3 when AUTH_MODE=oidc
}

func newOIDCVerifier(cfg Config) (*oidcVerifier, error) {
	if cfg.OIDCIssuer == "" || cfg.OIDCClientID == "" {
		return nil, fmt.Errorf("OIDC_ISSUER and OIDC_CLIENT_ID required when AUTH_MODE=oidc")
	}
	return nil, fmt.Errorf("oidc verifier not configured — complete P4.3")
}

func (m *Middleware) authenticateOIDC(r *http.Request) (*models.User, int, string) {
	if m.oidc == nil {
		return nil, http.StatusNotImplemented, "oidc auth not available"
	}
	return nil, http.StatusNotImplemented, "oidc auth not implemented"
}
