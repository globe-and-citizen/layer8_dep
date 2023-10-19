package utils

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
)

// AuthCodeClaims represents the claims of an auth code
type AuthCodeClaims struct {
	ClientID    string `json:"cid"`
	UserID      int64  `json:"uid"`
	RedirectURI string `json:"ruri"`
	Scopes      string `json:"scp"`
	ExpiresAt   int64  `json:"exp"`
	jwt.StandardClaims
}

// GenerateRandomString generates random bytes and encodes them to a hex string
func GenerateRandomString(count int) (string, error) {
	buf := make([]byte, count)
	_, err := io.ReadFull(rand.Reader, buf)
	if err != nil {
		return "", fmt.Errorf("could not generate random bytes: %s", err)
	}

	return base64.URLEncoding.EncodeToString(buf), nil
}

// GenerateAuthCode generates a JWT self-encoded auth code
//
//	Args:
//		secret: secret key
//		claims: auth code claims
func GenerateAuthCode(secret string, claims *AuthCodeClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", fmt.Errorf("could not generate auth code: %s", err)
	}
	return tokenString, nil
}

// DecodeAuthCode decodes a JWT self-encoded auth code
//
//	Args:
//		secret: secret key
//		code: auth code
func DecodeAuthCode(secret, code string) (*AuthCodeClaims, error) {
	token, err := jwt.ParseWithClaims(code, &AuthCodeClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil {
		return nil, fmt.Errorf("could not decode auth code: %s", err)
	}
	claims, ok := token.Claims.(*AuthCodeClaims)
	if !ok {
		return nil, fmt.Errorf("could not decode auth code: %s", err)
	}
	if claims.ExpiresAt < jwt.TimeFunc().Unix() {
		return nil, fmt.Errorf("could not decode auth code: token expired")
	}
	return claims, nil
}

// GenerateUserToken generates a JWT token for a user
func GenerateUserToken(secret string, userID int64) (string, error) {
	claims := &jwt.StandardClaims{
		Subject:   fmt.Sprintf("%d", userID),
		ExpiresAt: jwt.TimeFunc().Add(time.Hour * 24 * 7).Unix(), // expires in 7 days
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", fmt.Errorf("could not generate user token: %s", err)
	}
	return tokenString, nil
}

// VerifyUserToken verifies a JWT token for a user and returns the user ID
func VerifyUserToken(secret, token string) (int64, error) {
	claims := &jwt.StandardClaims{}
	_, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil {
		return 0, fmt.Errorf("could not verify user token: %s", err)
	}
	if claims.ExpiresAt < jwt.TimeFunc().Unix() {
		return 0, fmt.Errorf("could not verify user token: token expired")
	}
	id, err := strconv.ParseInt(claims.Subject, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("could not verify user token: %s", err)
	}
	return id, nil
}
