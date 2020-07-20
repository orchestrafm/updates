package identity

import (
	"os"

	"github.com/Nerzal/gocloak"
	"github.com/spidernest-go/logger"
)

var idp gocloak.GoCloak
var token *gocloak.JWT

func Handshake() {
	var err error
	idp = gocloak.NewClient(os.Getenv("IDP_ADDR"))
	token, err = idp.LoginClient(os.Getenv("OIDC_CLIENT_ID"), os.Getenv("OIDC_CLIENT_SECRET"), os.Getenv("IDP_REALM"))
	if err != nil {
		logger.Fatal().
			Err(err).
			Msg("Handshake with IDP server failed.")
	}
}
