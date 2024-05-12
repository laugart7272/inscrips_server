package util

import (
	"github.com/rs/zerolog/log"

	"github.com/golang-jwt/jwt"
)

func ParseToken(tokenString string) (claims *Claims, err error) {
	config, err := LoadConfig(".")
	if err != nil {
		log.Fatal().Err(err).Msg("cannot load configuration:")
	}

	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.TokenSymetricKey), nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)

	if !ok {
		return nil, err
	}

	return claims, nil
}
