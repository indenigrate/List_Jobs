package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	jwt "github.com/golang-jwt/jwt/v5"
)

// jwt auth
func validateJWT(tokenString string) (*jwt.Token, error) {
	secret := os.Getenv("JWT_SECRET")
	// fmt.Printf("secret is %+v\n", secret)
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return []byte(secret), nil
	})
}

func withJWTAuth(handlerfunc http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("jwt-token")
		// fmt.Printf("%+v\n", tokenString)
		token, err := validateJWT(tokenString)
		// fmt.Printf("%+v\n", token)

		if err != nil {
			// fmt.Printf("%+v\n", err)
			respondWithError(w, http.StatusForbidden, "api error: invalid token")
			return
		}
		//writing the token to the environment context
		ctx := context.WithValue(r.Context(), "token", token)
		handlerfunc(w, r.WithContext(ctx))
	}
}

func createJWT(userRole string) (string, error) {
	// Create the Claims
	secret := os.Getenv("JWT_SECRET")
	claims := &jwt.MapClaims{
		"expiresAt": jwt.NewNumericDate(time.Now().Add(24 * time.Hour)), // Set expiration to 24 hours
		// "expiresAt": jwt.NewNumericDate(time.Unix(1516239022, 0)),
		"role": userRole,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}
