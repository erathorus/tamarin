package auth

import (
	"gitlab.com/horo-go/auth0"
	"gitlab.com/lattetalk/lattetalk/config"
)

var Auth0 *auth0.Auth0

const TokenCookie = "AID"

func init() {
	// no need to validate audience and issuer for now
	Auth0 = auth0.New(auth0.Config{
		SigningMethod:  auth0.HS256,
		Secret:         []byte(config.Config.Auth0.Secret),
		JWKSURI:        config.Config.Auth0.JWKSURI,
		TokenExtractor: auth0.ExtractTokenFromCookie(TokenCookie),
	})
}
