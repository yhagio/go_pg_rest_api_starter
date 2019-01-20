package auth

import (
	"time"

	"go_rest_pg_starter/models"

	jwt "github.com/dgrijalva/jwt-go"
)

func IssueJWT(user *models.User, signingKey string) (string, error) {
	var signKey = []byte(signingKey)

	token := jwt.New(jwt.SigningMethodHS256)

	// Create a map to store the custom claims
	claims := token.Claims.(jwt.MapClaims)

	// Set token claims
	claims["role"] = "standard_user"
	claims["logged_in_user_id"] = user.ID
	claims["exp"] = time.Now().Add(time.Hour * time.Duration(1)).Unix()
	claims["iat"] = time.Now().Unix()

	// Sign the token with the secret key
	tokenString, err := token.SignedString(signKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
