package util

import "github.com/golang-jwt/jwt"

type Claims struct {
	Role  string `json:"role"`
	Token string `json:"accesstoken"`
	jwt.StandardClaims
}
