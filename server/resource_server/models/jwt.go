package models

import "github.com/golang-jwt/jwt/v4"

type Claims struct {
	UserName string `json:"username"`
	UserID   uint   `json:"user_id"`
	jwt.RegisteredClaims
}
