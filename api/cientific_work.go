package api

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	db "github.com/laugart7272/inscrips/db/sqlc"
	"github.com/laugart7272/inscrips/util"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

//////////////////////////////////////////////////////////// Create CientificWork

type createCientificWorkRequest struct {
	Title            string `json:"title" binding:"required"`
	AuthorID         int32  `json:"author_id" binding:"required,min=1"`
	InscriptionID    int32  `json:"inscription_id"`
	Resume           string `json:"resume" binding:"required"`
	File             string `json:"file"`
	PresentationType string `json:"presentation_type"`
	ExpositionType   string `json:"exposition_type"`
}

func (server *Server) createCientificWork(ctx *gin.Context) {
	var req createCientificWorkRequest
	// Bind JSON Arguments
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	//Authorization
	payload, err := server.authorizeUser(ctx, []string{util.AdminRole, util.UserRole}) //authPayload
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	//Find Active Event
	event, err := server.store.ActiveEvent(ctx)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusBadRequest, gin.H{"status": "event not found: %s"})
			return
		}
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "failed to found active event"})
		return
	}

	///// Inscription
	inscription, inscription_found, err := server.FindCreateInscription(ctx, int32(payload.User_ID), int32(event.ID))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": err.Error()})
		return
	}

	arg := db.CreateCientificWorkParams{
		Title:            req.Title,
		AuthorID:         req.AuthorID,
		InscriptionID:    int32(inscription.ID),
		Resume:           req.Resume,
		File:             req.File,
		PresentationType: req.PresentationType,
		ExpositionType:   req.ExpositionType,
	}

	if req.Resume == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "Must have Resume"})
		return
	}

	cientificWork, err := server.store.CreateCientificWork(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	//If inscription not found then create a payment
	// log.Info().Msgf("Inscription Found %v", inscription_found)
	if !inscription_found {
		// Create Payment
		//Find User
		user, err := server.store.GetUser(ctx, int64(req.AuthorID))
		if err != nil {
			if err == sql.ErrNoRows {
				ctx.JSON(http.StatusNotFound, errorResponse(err))
				return
			}
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}

		if payload.Role != util.AdminRole {
			if payload.Email != user.Email {
				ctx.JSON(http.StatusBadRequest, gin.H{"status": "cannot update other user's info"})
				return
			}
		}

		//Find Payment Amount
		arg_amount := db.GetPaymentAmountParams{
			EventID:          int32(event.ID),
			UserType:         user.UserType,
			PresentationType: arg.PresentationType,
		}

		amount, err := server.store.GetPaymentAmount(ctx, arg_amount)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}

		arg_payment := db.CreatePaymentParams{
			UserID:           arg.AuthorID,
			EventID:          int32(event.ID),
			UserType:         user.UserType,
			PresentationType: arg.PresentationType,
			Amount:           amount, //From Payment Rate
			InscriptionID:    int32(inscription.ID),
		}

		_, err = server.store.CreatePayment(ctx, arg_payment)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
	}

	ctx.JSON(http.StatusOK, cientificWork)
}

func (server *Server) FindCreateInscription(ctx *gin.Context, user_id int32, event_id int32) (*db.Inscription, bool, error) {
	found := false
	//Find if Exists
	arg_findinscription := db.FindInscriptionParams{
		UserID:  user_id,
		EventID: event_id,
	}

	inscription, err := server.store.FindInscription(ctx, arg_findinscription)
	if err != nil {
		if err == sql.ErrNoRows {
			// return nil, status.Errorf(codes.Internal, "inscription not foud: %s", err)
			//Creating Inscription
			arg_inscription := db.CreateInscriptionParams{
				UserID:  user_id,
				EventID: event_id,
			}

			inscription, err = server.store.CreateInscription(ctx, arg_inscription)
			if err != nil {
				return nil, false, status.Errorf(codes.Internal, "failed to create Inscription: %s - %v", err, arg_inscription)
			}
		}
	} else {
		//Inscription Found !
		found = true
	}

	return &inscription, found, nil
}

////////////////////////////////////////////////////////////// Get CientificWork

type getCientificWorkRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

func (server *Server) getCientificWork(ctx *gin.Context) {
	var req getCientificWorkRequest

	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
	}

	CientificWork, err := server.store.GetCientificWork(ctx, req.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, CientificWork)
}

//////////////////////////////////////////////////////////// Get CientificWorks

type listCientificWorkRequest struct {
	PageID   int32 `form:"page_id" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=5,max=10"`
}

func (server *Server) listCientificWorks(ctx *gin.Context) {
	var req listCientificWorkRequest

	//Authorization
	_, err := server.authorizeUser(ctx, []string{util.AdminRole, util.ReViewerRole}) //authPayload
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
	}

	arg := db.ListCientificWorksParams{
		Limit:  req.PageSize,
		Offset: (req.PageID - 1) * req.PageSize,
	}

	CientificWorks, err := server.store.ListCientificWorks(ctx, arg)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, CientificWorks)
}

////////////////////////////////////////////////////////

type listCientificWorkRequestByAuthorId struct {
	AuthorID int32 `form:"author_id" binding:"required,min=1"`
	PageID   int32 `form:"page_id" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=5,max=10"`
}

func (server *Server) listCientificsWorksByAuthorId(ctx *gin.Context) {
	var req listCientificWorkRequestByAuthorId
	//Authorization
	_, err := server.authorizeUser(ctx, []string{util.AdminRole, util.ReViewerRole, util.UserRole}) //authPayload
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
	}

	arg := db.GetCientificWorkByAuthorIdParams{
		AuthorID: int32(req.AuthorID),
		Limit:    req.PageSize,
		Offset:   ((req.PageID - 1) * req.PageSize),
	}

	cientificWorks, err := server.store.GetCientificWorkByAuthorId(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusNotFound, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, cientificWorks)
}

////////////////////////////////////////////////////////////////// Update CientificWork

type updateCientificWorkRequest struct {
	ID               int64  `json:"id" binding:"required"`
	Title            string `json:"title" binding:"required"`
	AuthorID         int32  `json:"author_id" binding:"required,min=1"`
	Resume           string `json:"resume" binding:"required"`
	File             string `json:"file"`
	PresentationType string `json:"presentation_type"`
	ExpositionType   string `json:"exposition_type"`
	InscriptionID    int64  `json:"inscription_id" binding:"required"`
}

func (server *Server) updateCientificWork(ctx *gin.Context) {
	var req updateCientificWorkRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
	}

	if req.Resume == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "Must have Resume"})
		return
	}

	arg := db.UpdateCientificWorkParams{
		ID:               req.ID,
		Title:            req.Title,
		AuthorID:         req.AuthorID,
		Resume:           req.Resume,
		File:             req.File,
		PresentationType: req.PresentationType,
		ExpositionType:   req.ExpositionType,
		InscriptionID:    int32(req.InscriptionID),
	}

	cientificWork, err := server.store.UpdateCientificWork(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, cientificWork)
}

//////////////////////////////////////////////////////////// Delete CientificWork

type deleteCientificWorkRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

func (server *Server) deleteCientificWork(ctx *gin.Context) {
	var req deleteCientificWorkRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
	}

	err := server.store.DeleteCientificWork(ctx, req.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, "Delete CientificWork Ok")
}
