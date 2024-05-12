package api

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	db "github.com/laugart7272/inscrips/db/sqlc"
)

//////////////////////////////////////////////////////////// Create Inscription

type createInscriptionRequest struct {
	UserID  int32 `json:"user_id" binding:"required,min=1"`
	EventID int32 `json:"event_id" binding:"required,min=1"`
}

func (server *Server) createInscription(ctx *gin.Context) {
	var req createInscriptionRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
	}

	arg := db.CreateInscriptionParams{
		UserID:  req.UserID,
		EventID: req.EventID,
	}

	Inscription, err := server.store.CreateInscription(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, Inscription)
}

////////////////////////////////////////////////////////////// Get Inscription

type getInscriptionRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

func (server *Server) getInscription(ctx *gin.Context) {
	var req getInscriptionRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
	}

	Inscription, err := server.store.GetInscription(ctx, req.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, Inscription)
}

//////////////////////////////////////////////////////////// Get Inscriptions

type listInscriptionRequest struct {
	PageID   int32 `form:"page_id" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=5,max=10"`
}

func (server *Server) listInscriptions(ctx *gin.Context) {
	var req listInscriptionRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
	}

	arg := db.ListInscriptionsParams{
		Limit:  req.PageSize,
		Offset: (req.PageID - 1) * req.PageSize,
	}

	inscriptions, err := server.store.ListInscriptions(ctx, arg)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, inscriptions)
}

/////////////////////////////////////////////////// Details Inscriptions

func (server *Server) listDetailsInscriptions(ctx *gin.Context) {
	var req listInscriptionRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
	}

	arg := db.ListDetailsInscriptionsParams{
		Limit:  req.PageSize,
		Offset: (req.PageID - 1) * req.PageSize,
	}

	inscriptions, err := server.store.ListDetailsInscriptions(ctx, arg)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, inscriptions)
}

////////////////////////////////////////////////////////////////// Update Inscription

type updateInscriptionRequest struct {
	ID      int64 `json:"id" binding:"required"`
	UserID  int32 `json:"user_id" binding:"required,min=1"`
	EventID int32 `json:"event_id" binding:"required,min=1"`
}

func (server *Server) updateInscription(ctx *gin.Context) {
	var req updateInscriptionRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
	}

	arg := db.UpdateInscriptionParams{
		ID:      req.ID,
		UserID:  req.UserID,
		EventID: req.EventID,
	}

	_, err := server.store.UpdateInscription(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, "Inscription update Ok")
}

//////////////////////////////////////////////////////////// Delete Inscription

type deleteInscriptionRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

func (server *Server) deleteInscription(ctx *gin.Context) {
	var req deleteInscriptionRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
	}

	err := server.store.DeleteInscription(ctx, req.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, "Delete Inscription Ok")
}
