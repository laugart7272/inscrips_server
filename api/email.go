package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	db "github.com/laugart7272/inscrips/db/sqlc"
	"github.com/laugart7272/inscrips/val"
)

type VerifyEmailTxRequest struct {
	EmailId    int64  `form:"email_id" binding:"required,min=1"`
	SecretCode string `form:"secret_code" binding:"required"`
}

func (server *Server) verifyEmail(ctx *gin.Context) {
	var req VerifyEmailTxRequest

	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// log.Info().Msgf("Req Email %v", req)

	if err := val.ValidateEmailId(req.EmailId); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	if err := val.ValidateSecretCode(req.SecretCode); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	txResult, err := server.store.VerifyEmailTx(ctx, db.VerifyEmailTxParams{
		EmailId:    req.EmailId,
		SecretCode: req.SecretCode,
	})

	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, txResult.VerifyEmail.IsUsed)
}
