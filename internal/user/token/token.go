package token

import (
	"time"

	"github.com/golang-jwt/jwt"
)

type AppInfo struct {
	Email string
	AppID string
}

type Token struct {
	*jwt.StandardClaims
	AppInfo
}

func GenerateToken(email, appID string) (jwt.Token, error) {
	t := jwt.New(jwt.GetSigningMethod("RS256"))
	t.Claims = &Token{
		&jwt.StandardClaims{
			ExpiresAt: time.Now().Add(1 * time.Hour).Unix(),
		},
		AppInfo{
			Email: email,
			AppID: appID,
		},
	}
	return *t, nil
}
