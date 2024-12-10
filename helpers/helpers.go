package helpers

import (
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
