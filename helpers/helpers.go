package helpers

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

func ErrResponse(statusCode int, message string, code string) ErrorResponse {
	return ErrorResponse{
		Code:    code,
		Message: message,
	}
}

func SuccessResponse(data interface{}) ApiResponse {
	return ApiResponse{
		Code:    200,
		Message: "success",
		Data:    data,
	}
}

type jwtClaims struct {
	UserID string `json:"_id" bson:"_id"`
	jwt.StandardClaims
}

func GenerateAccessToken(email string) (string, error) {

	expirationTime := time.Now().Add(5 * time.Minute)

	claims := &jwt.RegisteredClaims{
		Subject:   email,
		ExpiresAt: jwt.NewNumericDate(expirationTime),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	secretKey := []byte("secret")
	signedToken, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

func VerifyToken(tokenString string) (*jwtClaims, error) {
	fmt.Println("Received Token:", tokenString)

	claims := &jwtClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte("secret"), nil
	})

	if err != nil {
		return nil, errors.New("invalid token")
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}
