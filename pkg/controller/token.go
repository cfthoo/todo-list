package controller

import (
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
)

var SECRET []byte = make([]byte, 0)

// CreateJWT creates JWT token
func CreateJWT(secret []byte) (string, error) {

	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)

	claims["exp"] = time.Now().Add(time.Hour).Unix()

	tokenStr, err := token.SignedString(secret)
	SECRET = secret
	if err != nil {
		return "", err
	}

	return tokenStr, nil
}

// ValidateJWT validates the token from Token header
func ValidateJWT(next func(w http.ResponseWriter, r *http.Request)) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if r.Header["Authorization"] != nil {
			authHeader := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")
			token, err := jwt.Parse(authHeader, func(t *jwt.Token) (interface{}, error) {
				_, ok := t.Method.(*jwt.SigningMethodHMAC)
				if !ok {
					w.WriteHeader(http.StatusUnauthorized)
					w.Write([]byte("not authorized"))
				}
				return SECRET, nil
			})
			if err != nil {
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte("not authorized: " + err.Error()))
			}

			if token.Valid {
				next(w, r)
			}
		} else {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("not authorized"))
		}
	})
}
