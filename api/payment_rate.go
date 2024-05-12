package api

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	db "github.com/laugart7272/inscrips/db/sqlc"
)

////////////////////////////////////////////////////////////// Get PaymentAmount

type GetPaymentAmountRequest struct {
	EventID          int32  `json:"event_id"`
	UserType         string `json:"user_type"`
	PresentationType string `json:"presentation_type"`
}

func (server *Server) getPaymentAmount(ctx *gin.Context) {
	var req GetPaymentAmountRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
	}

	arg := db.GetPaymentAmountParams{
		EventID:          req.EventID,
		UserType:         req.UserType,
		PresentationType: req.PresentationType,
	}

	PaymentAmount, err := server.store.GetPaymentAmount(ctx, arg)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, PaymentAmount)
}

func (server *Server) listDetailsPaymentRates(ctx *gin.Context) {
	var req listPaymentRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
	}

	arg := db.ListDetailsPaymentRatesParams{
		Limit:  req.PageSize,
		Offset: (req.PageID - 1) * req.PageSize,
	}

	Payments, err := server.store.ListDetailsPaymentRates(ctx, arg)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, Payments)
}
