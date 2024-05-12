package api

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	db "github.com/laugart7272/inscrips/db/sqlc"
)

//////////////////////////////////////////////////////////// Create Payment

type createPaymentRequest struct {
	UserID           int32  `json:"user_id" binding:"required,min=1"`
	EventID          int32  `json:"event_id" binding:"required,min=1"`
	UserType         string `json:"user_type_id" binding:"required,min=1"`
	PresentationType string `json:"presentation_type"`
	Amount           int32  `json:"amount" binding:"required"`
}

func (server *Server) createPayment(ctx *gin.Context) {
	var req createPaymentRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
	}

	arg := db.CreatePaymentParams{
		UserID:           req.UserID,
		EventID:          req.EventID,
		UserType:         req.UserType,
		PresentationType: req.PresentationType,
		Amount:           req.Amount,
	}

	Payment, err := server.store.CreatePayment(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, Payment)
}

////////////////////////////////////////////////////////////// Get Payment

type getPaymentRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

func (server *Server) getPayment(ctx *gin.Context) {
	var req getPaymentRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
	}

	Payment, err := server.store.GetPayment(ctx, req.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, Payment)
}

//////////////////////////////////////////////////////////// Get Payments

type listPaymentRequest struct {
	PageID   int32 `form:"page_id" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=5,max=10"`
}

func (server *Server) listPayments(ctx *gin.Context) {
	var req listPaymentRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
	}

	arg := db.ListPaymentsParams{
		Limit:  req.PageSize,
		Offset: (req.PageID - 1) * req.PageSize,
	}

	Payments, err := server.store.ListPayments(ctx, arg)
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

/////////////////////////////////////////////////// List Details Payments

func (server *Server) listDetailsPayments(ctx *gin.Context) {
	var req listPaymentRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
	}

	arg := db.ListDetailsPaymentsParams{
		Limit:  req.PageSize,
		Offset: (req.PageID - 1) * req.PageSize,
	}

	Payments, err := server.store.ListDetailsPayments(ctx, arg)
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

////////////////////////////////////////////////////////////////// Update Payment

type updatePaymentRequest struct {
	ID               int64  `json:"id" binding:"required"`
	UserID           int32  `json:"user_id" binding:"required,min=1"`
	EventID          int32  `json:"event_id" binding:"required,min=1"`
	UserTypeID       string `json:"user_type"`
	PresentationType string `json:"presentation_type"`
	Amount           int32  `json:"amount" binding:"required"`
}

func (server *Server) updatePayment(ctx *gin.Context) {
	var req updatePaymentRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
	}

	arg := db.UpdatePaymentParams{
		ID:               req.ID,
		UserID:           req.UserID,
		EventID:          req.EventID,
		PresentationType: req.PresentationType,
		Amount:           req.Amount,
	}

	_, err := server.store.UpdatePayment(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, "Payment update Ok")
}

//////////////////////////////////////////////////////////// Delete Payment

type deletePaymentRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

func (server *Server) deletePayment(ctx *gin.Context) {
	var req deletePaymentRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
	}

	err := server.store.DeletePayment(ctx, req.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, "Delete Payment Ok")
}
