package authenticate

import (
	"gin-template-app/pkg/net/http/ctx"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

const secret = "gbbmn"

// Claims authenticate customize claims
type Claims struct {
	ctx.UserInfo
	jwt.StandardClaims
}

// GenerateToken authenticate generate token
func GenerateToken(user ctx.UserInfo) (string, error) {
	// Create a new token object, specifying signing method and the claims
	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		user,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour).Unix(),
			Issuer:    "blog",
		},
	})
	// Sign and get the complete encoded token as a string using the secret
	sercret := []byte(secret)
	return tokenClaims.SignedString(sercret)
}

// ParseToken authenticate parse token
func ParseToken(token string) (ctx.UserInfo, error) {
	tokenClaims, err := jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})

	if tokenClaims != nil {
		if claims, ok := tokenClaims.Claims.(*Claims); ok && tokenClaims.Valid {
			return claims.UserInfo, nil
		}
	}
	return ctx.UserInfo{}, err
}
