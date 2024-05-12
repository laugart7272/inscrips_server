package api

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/laugart7272/inscrips/token"
	"github.com/laugart7272/inscrips/util"
)

// const (
// 	authorizationHeader = "authorization"
// 	authorizationBearer = "bearer"
// )

func (server *Server) authorizeUser(ctx *gin.Context, accesibleRoles []string) (*token.Payload, error) {

	cookie, err := ctx.Cookie("token")

	if err != nil {
		ctx.JSON(401, gin.H{"error": "unauthorized"})
		return nil, fmt.Errorf("unauthorized cookie %v", err)
	}

	claims, err := util.ParseToken(cookie)
	if err != nil {
		return nil, fmt.Errorf("error parse token %v", err)
	}

	accessToken := claims.Token
	if len(accessToken) == 0 {
		return nil, fmt.Errorf("missing authorization header")
	}

	payload, err := server.tokenMaker.VerifyToken(accessToken)
	if err != nil {
		return nil, fmt.Errorf("invalid access token: %s", err)
	}

	// log.Info().Msgf("Inscription Author Role %s", payload.Role)
	if !hasPermission(payload.Role, accesibleRoles) {
		return nil, fmt.Errorf("permission denied")
	}

	return payload, nil
}

func hasPermission(userRole string, accesibleRoles []string) bool {
	for _, role := range accesibleRoles {
		if userRole == role {
			return true
		}
	}

	return false
}
