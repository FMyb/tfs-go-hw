package utils

import (
	"context"
	"log"
	"net/http"

	"github.com/golang-jwt/jwt"

	"github.com/FMyb/tfs-go-hw/lection05/homework/domain"
)

type key int

const (
	KeyUserID key = iota
)

type Token struct {
	jwt.StandardClaims
	UserID uint
}

var superSecretKey = "SUPERSECRETJWTTOKEN"

func GetToken(user domain.User) string {
	tk := &Token{UserID: user.ID}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, tk)
	tokenString, err := token.SignedString([]byte(superSecretKey))
	if err != nil {
		log.Printf("Error: %s", err)
		return ""
	}
	return tokenString
}

func ParseToken(tokenHeader string) *Token {
	tk := &Token{}
	token, err := jwt.ParseWithClaims(tokenHeader, tk, func(token *jwt.Token) (interface{}, error) {
		return []byte(superSecretKey), nil
	})

	if err != nil {
		log.Printf("Unvalid token")
		return nil
	}
	if !token.Valid {
		log.Printf("Unvalid token")
		return nil
	}
	return tk
}

func JwtAuthentication(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		tokenHeader := r.Header.Get("Authorization")
		if tokenHeader == "" {
			w.WriteHeader(http.StatusForbidden)
			return
		}
		tk := ParseToken(tokenHeader)
		if tk == nil {
			w.WriteHeader(http.StatusForbidden)
			return
		}
		ctx := context.WithValue(r.Context(), KeyUserID, tk.UserID)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}
