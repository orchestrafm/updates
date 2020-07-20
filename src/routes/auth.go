package routes

import (
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/mitchellh/mapstructure"
	"github.com/spidernest-go/logger"
	"github.com/spidernest-go/mux"
)

type jwtExtendedClaims struct {
	RealmAccess struct {
		Roles []string `json:"roles,omitempty"`
	} `json:"realm_access,omitempty"`
	Scope             string `json:"scope"`
	PreferredUsername string `json:"preferred_username,omitempty"`
	Audience          string `json:"aud,omitempty"`
	ExpiresAt         int64  `json:"exp,omitempty"`
	Id                string `json:"jti,omitempty"`
	IssuedAt          int64  `json:"iat,omitempty"`
	Issuer            string `json:"iss,omitempty"`
	NotBefore         int64  `json:"nbf,omitempty"`
	Subject           string `json:"sub,omitempty"`
	//jwt.StandardClaims
}

func decodeToClaims(src, dst interface{}) error {
	//BUGFIX: Introduced in b7fdf80, I'd like to get rid of this function
	//BUGFIX: commit is related to tracks microservice
	dec, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		TagName: "json",
		Result:  dst,
	})
	if err != nil {
		logger.Error().
			Err(err).
			Msg("mapstructure decoder could not be initialized")
		return err
	}
	err = dec.Decode(src)
	if err != nil {
		logger.Error().
			Err(err).
			Msg("structure could not be decoded")
		return err
	}
	return nil
}

func FullAuthCheck(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	//claims := user.Claims.(*jwtExtendedClaims)
	claims := new(jwtExtendedClaims)
	err := decodeToClaims(user.Claims, claims)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, &struct {
			Message string
		}{
			Message: "Authorization Token could not be parsed."})
	}

	auth := strings.Contains(claims.Scope, "track:write") || strings.Contains(claims.Scope, "track:admin")
	admin := strings.Contains(claims.Scope, "track:admin")
	if !admin || !auth {
		return c.JSON(http.StatusUnauthorized, &struct {
			Message string
		}{
			Message: "Insufficient Permissions."})
	}
	return nil
}

func AuthorizationCheck(c echo.Context) (bool, bool) {
	user := c.Get("user").(*jwt.Token)
	//claims := user.Claims.(*jwtExtendedClaims)
	claims := new(jwtExtendedClaims)
	err := decodeToClaims(user.Claims, claims)
	if err != nil {
		return false, false
	}

	auth := strings.Contains(claims.Scope, "track:write") || strings.Contains(claims.Scope, "track:admin")
	admin := strings.Contains(claims.Scope, "track:admin")

	return admin, auth
}

func SelfAuthCheck(c echo.Context) *jwtExtendedClaims {
	user := c.Get("user").(*jwt.Token)
	claims := new(jwtExtendedClaims)
	err := decodeToClaims(user.Claims, claims)
	if err != nil {
		return nil
	}
	return claims
}
