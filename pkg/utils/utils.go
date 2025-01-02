package utils

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
)

var SecretKey = []byte("your-secret-key")

func CreateJWT(username string) (string, error) {
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": username,                         // Subject (user identifier)
		"iss": "social-media",                   // Issuer
		"aud": "coder",                          // Audience (user role)
		"exp": time.Now().Add(time.Hour).Unix(), // Expiration time
		"iat": time.Now().Unix(),                // Issued at
	})

	fmt.Printf("Token claims added: %+v", claims)

	tokenString, err := claims.SignedString(SecretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func VerifyJWT(tokenString string) (*jwt.Token, error) {
	// Parse the token with the secret key
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return SecretKey, nil
	})

	// Check for verification errors
	if err != nil {
		return nil, err
	}

	// Check if the token is valid
	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("no claims")
	}

	if iss := claims["iss"]; iss != "social-media" {
		return nil, errors.New("invalid issuer")
	}

	if float64(time.Now().Unix()) > claims["exp"].(float64) {
		return nil, errors.New("token expired")
	}

	// Return the verified token
	return token, nil
}
