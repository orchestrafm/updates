package identity

import (
	"crypto/rand"
	"encoding/base64"
	"os"

	oidc "github.com/coreos/go-oidc"
	"github.com/spidernest-go/logger"
	"golang.org/x/net/context"
	"golang.org/x/oauth2"
)

var (
	clientID     = os.Getenv("OIDC_CLIENT_ID")
	clientSecret = os.Getenv("OIDC_CLIENT_SECRET")

	Context              context.Context
	NonceEnabledVerifier *oidc.IDTokenVerifier
	OAuth2               oauth2.Config
	Nonce                string
)

func EnableOIDC() {
	// generate a random, cryptographically secure, 32-bit nonce
	buf := make([]byte, 32)
	_, err := rand.Read(buf)
	if err != nil {
		logger.Fatal().
			Err(err).
			Msg("Buffer for generating secure 32-bit nonce failed.")
	}
	Nonce = base64.StdEncoding.EncodeToString(buf)
	Context = context.Background()

	provider, err := oidc.NewProvider(Context, os.Getenv("OIDC_URL"))
	if err != nil {
		logger.Fatal().
			Err(err).
			Msg("OpenID Connect provider's configuration was not located at the specified URL.")
	}

	oidcConfig := &oidc.Config{
		ClientID: clientID,
	}

	NonceEnabledVerifier = provider.Verifier(oidcConfig)

	OAuth2 = oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Endpoint:     provider.Endpoint(),
		RedirectURL:  "http://localhost:5004/api/v0/oidc/callback",
		Scopes:       []string{oidc.ScopeOpenID, "profile", "email", "track:write", "track:admin", "board:write", "board:admin", "score:write", "score:admin"},
	}

	logger.Info().
		Msg("OpenID Connect module loaded successfully.")
}
