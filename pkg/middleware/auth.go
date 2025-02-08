package middleware

import (
	"context"
	"go-flashcards-server/pkg/config"
	"log"
	"net/http"
	"strconv"

	"github.com/dgrijalva/jwt-go"
)


func AuthMiddleware(next http.Handler) http.Handler {
	log.Println("AuthMiddleware()")
	jwtKey := config.JWTSecretKey
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("cookies: ", r.Cookies())
		cookie, err := r.Cookie("token")
		if err != nil {
			if err == http.ErrNoCookie {
				http.Error(w, "JWT cookie missing", http.StatusUnauthorized)
				return
			}
			http.Error(w, "Error reading cookie", http.StatusBadRequest)
			return
		}

		tokenString := cookie.Value
		claims := &jwt.MapClaims{}
		_, err = jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})
		if err != nil {
			if ve, ok := err.(*jwt.ValidationError); ok && ve.Errors == jwt.ValidationErrorExpired {
				http.Error(w, "Token has expired", http.StatusUnauthorized)
			} else {
				http.Error(w, "Invalid or malformed token", http.StatusUnauthorized)
			}
			return
		}

		sub, ok := (*claims)["sub"].(string)
		if !ok {
			http.Error(w, "User ID not found", http.StatusUnauthorized)
			return
		}

		userID, err := strconv.Atoi(sub)
		if err != nil {
			http.Error(w, "Invalid User ID", http.StatusUnauthorized)
		}

		ctx := context.WithValue(r.Context(), "userID", userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
