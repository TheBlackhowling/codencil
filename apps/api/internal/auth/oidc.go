package auth

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/TheBlackHowling/codencil/apps/api/internal/models"
	"github.com/TheBlackHowling/codencil/apps/api/internal/store"
	"github.com/coreos/go-oidc/v3/oidc"
)

type oidcVerifier struct {
	verifier *oidc.IDTokenVerifier
}

func newOIDCVerifier(cfg Config) (*oidcVerifier, error) {
	if cfg.OIDCIssuer == "" || cfg.OIDCClientID == "" {
		return nil, fmt.Errorf("OIDC_ISSUER and OIDC_CLIENT_ID required when AUTH_MODE=oidc")
	}
	provider, err := oidc.NewProvider(context.Background(), cfg.OIDCIssuer)
	if err != nil {
		return nil, fmt.Errorf("oidc provider: %w", err)
	}
	return &oidcVerifier{
		verifier: provider.Verifier(&oidc.Config{ClientID: cfg.OIDCClientID}),
	}, nil
}

func (m *Middleware) authenticateOIDC(r *http.Request) (*models.User, error) {
	if r.Header.Get(devUserHeader) != "" {
		return nil, authError{status: http.StatusUnauthorized, message: "dev auth header not allowed when AUTH_MODE=oidc"}
	}
	authHeader := r.Header.Get("Authorization")
	if !strings.HasPrefix(authHeader, "Bearer ") {
		return nil, authError{status: http.StatusUnauthorized, message: "missing bearer token"}
	}
	rawToken := strings.TrimSpace(strings.TrimPrefix(authHeader, "Bearer "))
	idToken, err := m.oidc.verifier.Verify(r.Context(), rawToken)
	if err != nil {
		return nil, authError{status: http.StatusUnauthorized, message: "invalid token"}
	}

	var claims struct {
		Sub   string `json:"sub"`
		Email string `json:"email"`
		Name  string `json:"name"`
	}
	if err := idToken.Claims(&claims); err != nil {
		return nil, authError{status: http.StatusUnauthorized, message: "invalid token claims"}
	}
	displayName := claims.Name
	if displayName == "" {
		displayName = claims.Email
	}
	user, err := m.store.GetOrCreateUser(r.Context(), store.GetOrCreateUserInput{
		OrgID:       m.config.OrgID,
		ExternalID:  claims.Sub,
		DisplayName: displayName,
	})
	if err != nil {
		return nil, authError{status: http.StatusInternalServerError, message: "auth failed"}
	}
	return user, nil
}

type authError struct {
	status  int
	message string
}

func (e authError) Error() string {
	return e.message
}
