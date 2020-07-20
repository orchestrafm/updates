package identity

import (
	"os"

	"github.com/Nerzal/gocloak"
)

func GetAccount(uuid string) (*gocloak.User, error) {
	return idp.GetUserByID(token.AccessToken,
		os.Getenv("IDP_REALM"),
		uuid)
}

func RefreshToken(ref string) (*gocloak.JWT, error) {
	return idp.RefreshToken(
		ref,
		os.Getenv("OIDC_CLIENT_ID"),
		os.Getenv("OIDC_CLIENT_SECRET"),
		os.Getenv("IDP_REALM"))
}
