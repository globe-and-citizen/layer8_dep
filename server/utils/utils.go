package utils

import (
	"crypto/sha1"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/xdg-go/pbkdf2"
)

type JWK struct {
	Key_ops []string `json:"use,omitempty"` // ["sign", "verify", "encrypt", "decrypt", "wrapKey", "unwrapKey", "deriveKey", "deriveBits"]
	Kty     string   `json:"kty,omitempty"` // "EC", "RSA"
	Kid     string   `json:"kid,omitempty"` // Key ID
	Crv     string   `json:"crv,omitempty"` // "P-256"
	X       string   `json:"x,omitempty"`   // x coordinate as base64 URL encoded string.
	Y       string   `json:"y,omitempty"`   // y coordinate as base64 URL encoded string.
	D       string   `json:"d,omitempty"`   // d coordinate as base64 URL encoded string. Private keys only.
}

func SaltAndHashPassword(password string, salt string) string {
	dk := pbkdf2.Key([]byte(password), []byte(salt), 4096, 32, sha1.New)
	return hex.EncodeToString(dk[:])
}

func GenerateStandardToken(secretKey string) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := &jwt.StandardClaims{
		ExpiresAt: time.Now().Add(time.Hour * 24 * 7).Unix(),
	}

	token.Claims = claims

	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", fmt.Errorf("could not generate standard token: %s", err)
	}

	return tokenString, nil
}

func B64ToJWK(userPubJWK string) (*JWK, error) {
	userPubJWK_BS, err := base64.StdEncoding.DecodeString(userPubJWK)
	if err != nil {
		return nil, fmt.Errorf("Failure to decode userPubJWK", err.Error())
	}
	userPubJWKConverted := &JWK{}
	err = json.Unmarshal(userPubJWK_BS, userPubJWKConverted)
	if err != nil {
		return nil, fmt.Errorf("Failure to unmarshal userPubJWK: ", err.Error())
	}
	return userPubJWKConverted, nil
}

func VerifyStandardToken(tokenString string, secretKey string) (*jwt.StandardClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*jwt.StandardClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}

// func VerifyUserToken(secret, token string) (int64, error) {
// 	claims := &jwt.StandardClaims{}
// 	_, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
// 		return []byte(secret), nil
// 	})
// 	if err != nil {
// 		return 0, fmt.Errorf("could not verify user token: %s", err)
// 	}
// 	if claims.ExpiresAt < jwt.TimeFunc().Unix() {
// 		return 0, fmt.Errorf("could not verify user token: token expired")
// 	}
// 	id, err := strconv.ParseInt(claims.Subject, 10, 64)
// 	if err != nil {
// 		return 0, fmt.Errorf("could not verify user token: %s", err)
// 	}
// 	return id, nil
// }

// func GenerateUserToken(secret string, userID int64) (string, error) {
// 	claims := &jwt.StandardClaims{
// 		Subject:   fmt.Sprintf("%d", userID),
// 		ExpiresAt: jwt.TimeFunc().Add(time.Hour * 24 * 7).Unix(), // expires in 7 days
// 	}
// 	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
// 	tokenString, err := token.SignedString([]byte(secret))
// 	if err != nil {
// 		return "", fmt.Errorf("could not generate user token: %s", err)
// 	}
// 	return tokenString, nil
// }

// func ValidateRequiredFields(s interface{}) error {
// 	err := validator.Validate(s)
// 	if err != nil {
// 		log.Println(err.(validator.ErrorMap))
// 		return errors.New("missing fields")
// 	}
// 	return nil
// }

// // ValidateEmail validates if the email is valid.
// func ValidateEmail(email string) error {
// 	regex := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
// 	if !regex.MatchString(email) {
// 		return errors.New("invalid email")
// 	}
// 	return nil
// }
