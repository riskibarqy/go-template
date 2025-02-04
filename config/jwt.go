package config

import (
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/riskibarqy/go-template/models"
)

// Define a struct for the JWT claims (you can customize this as needed)
type Claims struct {
	ID int `json:"id"`
	jwt.StandardClaims
}

func GenerateJWTToken(user *models.User) (string, error) {
	// Set the claims
	claims := Claims{
		ID: user.ID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 72).Unix(), // 72-hour expiration
			Issuer:    "account-app",
		},
	}

	// Create the token with the specified claims and sign it with the secret key
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Generate the signed token string
	signedToken, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		return "", err
	}

	return signedToken, nil
}
