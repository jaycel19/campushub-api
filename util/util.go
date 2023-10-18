package util

import (
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// TODO: Create a access token payload

func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed to has password: %w", err)
	}

	return string(hashedPassword), nil
}

func CheckPassword(password string, hashedPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

func GenerateAccessToken(username string) (string, uuid.UUID, error) {
	// TODO: add a token payload
	tokenID, err := uuid.NewRandom()
	if err != nil {
		return "", uuid.UUID{}, err
	}
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	exp := time.Now().Add(time.Minute * 15).Unix()
	claims["id"] = tokenID
	claims["username"] = username
	claims["exp"] = exp
	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET")))
	return tokenString, tokenID, nil
}

func GenerateRefreshToken(username string, id uuid.UUID) (string, error) {
	// TODO: add a refreshtoken payload
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	exp := time.Now().Add(time.Hour * 24 * 30).Unix()
	claims["id"] = id
	claims["username"] = username
	claims["exp"] = exp
	fmt.Println("token expiration", exp)
	tokenString, _ := token.SignedString([]byte(os.Getenv("SECRET")))
	return tokenString, nil
}

// VerifyToken verifies the refresh token and returns the ID (user ID or session ID) associated with it.
func VerifyToken(refreshToken string) (string, error) {
	// Define a function to validate the token
	tokenValidationFunction := func(token *jwt.Token) (interface{}, error) {
		// Check the signing method to ensure it's the one you expect
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			fmt.Print("Signing error")
			return nil, fmt.Errorf("Unexpected signing method")
		}

		// Return the secret key used to sign the token
		return []byte(os.Getenv("SECRET")), nil
	}

	// Parse and validate the token
	token, err := jwt.Parse(refreshToken, tokenValidationFunction)
	if err != nil {
		fmt.Print(err)
		return "", err
	}

	// Check if the token is valid and not expired
	if !token.Valid {
		fmt.Print("Valid error")
		return "", fmt.Errorf("Token is not valid")
	}

	// Extract the claims from the token
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		fmt.Print("Token claims error")
		return "", fmt.Errorf("Token claims are not valid")
	}

	// Extract the ID from the claims
	id, idExists := claims["id"].(string)
	if !idExists {
		fmt.Print("Id error")
		return "", fmt.Errorf("Token claims do not contain ID")
	}

	return id, nil
}
