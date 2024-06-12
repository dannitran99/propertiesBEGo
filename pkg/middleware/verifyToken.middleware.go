package middleware

import (
	"context"
	"net/http"
	"os"
	"strings"

	"github.com/dgrijalva/jwt-go"
)

func VerifyJWT(next http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        tokenString := r.Header.Get("Authorization")
        if tokenString == "" {
            http.Error(w, "Unauthorized", http.StatusUnauthorized)
            return
        }
		jwtKey := os.Getenv("JWT_SECRET_KEY")
		splitToken := strings.Split(tokenString, "Bearer ")
		tokenString = splitToken[1]
		claims := jwt.MapClaims{}
        token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
            return []byte(jwtKey), nil
        })
        if err != nil || !token.Valid {
            http.Error(w, "Unauthorized", http.StatusUnauthorized)
            return
        }
		ctx := context.WithValue(r.Context(), "username", claims["userId"].(string))
        // Token is valid, proceed to the next handler
        next.ServeHTTP(w, r.WithContext(ctx))
    }
}