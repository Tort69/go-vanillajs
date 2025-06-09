package token

import (
	"Allusion/logger"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt"
)

func ExtractJWTSecret(r *http.Request, logger logger.Logger) (email string, ok bool) {

	tokenStr := r.Header.Get("Authorization")
	if tokenStr == "" {
		return "", false
	}
	tokenStr = strings.TrimPrefix(tokenStr, "Bearer ")

	token, err := jwt.Parse(tokenStr,
		func(t *jwt.Token) (interface{}, error) {
			// Ensure the signing method is HMAC
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return []byte(GetJWTSecret(logger)), nil
		},
	)
	if err != nil || !token.Valid {
		return "", false
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", false
	}

	email, ok = claims["email"].(string)
	if !ok {
		return "", false
	}
	return email, true

}
