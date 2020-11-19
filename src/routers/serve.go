package routers

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"

	"github.com/spidernest-go/logger"
	"github.com/spidernest-go/mux"
	"github.com/spidernest-go/mux/middleware"
)

var r *echo.Echo

const (
	ErrGeneric   = `{"errno": "404", "message": "Bad Request"}`
	rsaPublicKey = `-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAl3Mw0lnzWr+KrhhP1/jnKHblCM/DqIhvUHgsOYZWrE3+fEvHjc6wrUT9RtC3eRZfRxtdyxa9CPuSnPEt/Jmu2YPVRWxOVUJfUxgZQg0OPXurMy0h6O1Yal4s9yNq0+OmCSIE3DFVNTs5hlYNI7TNkjPp/UJx8Xc+J+g/gUPrIVQo+XWNGoKv+udiQhi9LrYZuQOy9MZPKgUKSfJwmwWRBb7CZmvWSwprQ3/619+2vf1gS/K3vqenlZfCRFadPuxebmQ595LKAn0tgnw2R0c4aAU/G1LsJsBFfY0kvhE/asFvNSoAoJA3jnQMYmMekqgVdVNV2FLrLWve5520EjTeMQIDAQAB
-----END PUBLIC KEY-----`
)

func ListenAndServe() {
	// decode pem block into rsa public key
	block, _ := pem.Decode([]byte(rsaPublicKey))
	if block == nil {
		logger.Fatal().
			Msg("PEM RSA public key block was invalid and failed to decode.")
	}
	pkey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		logger.Fatal().
			Err(err).
			Msg("Decoded PEM block failed to parse.")
	}
	rsaKey, ok := pkey.(*rsa.PublicKey)
	if !ok {
		logger.Fatal().
			Msgf("got unexpected key type: %T", pkey)
	}

	// Start serving API routes
	r = echo.New()
	r.BodyLimit(32 * 1024 * 1024) // 32 MB

	v0 := r.Group("/api/v0")
	v0AuthReq := v0.Group("", middleware.JWTWithConfig(middleware.JWTConfig{
		SigningMethod: "RS256",
		SigningKey:    rsaKey,
	}))
	v0AuthReq.POST("/patch", pushUpdate)
	v0.GET("/patch", listUpdates)

	r.Start(":5004")
}
