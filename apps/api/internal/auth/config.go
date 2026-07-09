package auth

import "os"

const (
	ModeDev  = "dev"
	ModeOIDC = "oidc"

	DefaultDevUserID = "dev-user"
	DefaultOrgID     = "dev"
)

// Config holds auth mode settings from environment.
type Config struct {
	Mode         string
	OrgID        string
	OIDCIssuer   string
	OIDCClientID string
}

func LoadConfig() Config {
	mode := os.Getenv("AUTH_MODE")
	if mode == "" {
		mode = ModeDev
	}
	orgID := os.Getenv("ORG_ID")
	if orgID == "" {
		orgID = DefaultOrgID
	}
	return Config{
		Mode:         mode,
		OrgID:        orgID,
		OIDCIssuer:   os.Getenv("OIDC_ISSUER"),
		OIDCClientID: os.Getenv("OIDC_CLIENT_ID"),
	}
}
