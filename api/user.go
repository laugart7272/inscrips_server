package api

import (
	"database/sql"
	"fmt"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/hibiken/asynq"
	db "github.com/laugart7272/inscrips/db/sqlc"
	"github.com/laugart7272/inscrips/util"
	"github.com/laugart7272/inscrips/worker"
	"github.com/lib/pq"
	"github.com/rs/zerolog/log"
)

//////////////////////////////////////////////////////////// Create User

type createUserRequest struct {
	Name     string `json:"name" binding:"required,alphanum"`
	LastName string `json:"last_name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Phone    string `json:"phone"`
	Password string `json:"password" binding:"required,min=6"`
	UserType string `json:"user_type"`
}

type userResponse struct {
	ID                uint64    `json:"id"`
	Name              string    `json:"name"`
	LastName          string    `json:"last_name"`
	Email             string    `json:"email"`
	Phone             string    `json:"phone"`
	UserType          string    `json:"user_type"`
	PasswordChangedAt time.Time `json:"password_changed_at"`
	Role              string    `json:"role"`
	AvatarPath        string    `json:"avatar_path"`
}

func newUserResponse(user db.User) userResponse {
	return userResponse{
		ID:                uint64(user.ID),
		Name:              user.Name,
		LastName:          user.LastName,
		Email:             user.Email,
		Phone:             user.Phone,
		UserType:          user.UserType,
		PasswordChangedAt: user.PasswordChangedAt,
		Role:              user.Role,
		AvatarPath:        user.AvatarPath,
	}
}

func (server *Server) createUser(ctx *gin.Context) {
	var req createUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	if req.Email == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "Must have Email"})
		return
	}

	if req.Password == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "Must have Password"})
		return
	}

	hashedPassword, err := util.HashPassword(req.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	arg := db.CreateUserTxParams{
		CreateUserParams: db.CreateUserParams{
			Name:           req.Name,
			LastName:       req.LastName,
			Email:          req.Email,
			Phone:          req.Phone,
			HashedPassword: hashedPassword,
			UserType:       req.UserType,
		},
		AfterCreate: func(user db.User) error {
			tasPayload := &worker.PayloadSendVerifyEmail{
				Email: user.Email,
			}

			opts := []asynq.Option{
				asynq.MaxRetry(10),
				asynq.ProcessIn(10 * time.Second),
				asynq.Queue(worker.QueueCritical),
			}

			return server.taskDistributor.DistributeTaskSendVerifyEmail(ctx, tasPayload, opts...)
		},
	}

	txResult, err := server.store.CreateUserTx(ctx, arg)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "unique_violation":
				ctx.JSON(http.StatusForbidden, errorResponse(err))
				return
			}
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// if user type is participant the payment

	if req.UserType == "Participante" {
		err = server.doPayment(ctx, txResult.User)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
	}

	resp := newUserResponse(txResult.User)

	ctx.JSON(http.StatusOK, resp)
}

// ////////////////////////////////// Create Payment for Participant
func (server *Server) doPayment(ctx *gin.Context, user db.User) error {
	//Find Active Event
	event, err := server.store.ActiveEvent(ctx)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusBadRequest, gin.H{"status": "event not found: %s"})
			return nil
		}
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "failed to found active event"})
		return nil
	}

	//Find Payment Amount
	arg_amount := db.GetPaymentAmountParams{
		EventID:          int32(event.ID),
		UserType:         user.UserType,
		PresentationType: "Presencial", //arg.PresentationType,
	}

	// log.Info().Msgf("Arg %v", arg_amount)

	amount, err := server.store.GetPaymentAmount(ctx, arg_amount)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return nil
	}

	// log.Info().Msgf("Amount %v", amount)

	//Creating Inscription for Participant
	arg_inscription := db.CreateInscriptionParams{
		UserID:  int32(user.ID),
		EventID: int32(event.ID),
	}

	inscription, err := server.store.CreateInscription(ctx, arg_inscription)
	if err != nil {
		return nil
	}

	arg_payment := db.CreatePaymentParams{
		UserID:           int32(user.ID),
		EventID:          int32(event.ID),
		UserType:         user.UserType,
		PresentationType: "Precencial",          //arg.PresentationType,
		Amount:           amount,                //From Payment Rate
		InscriptionID:    int32(inscription.ID), //int32(inscription.ID),
	}

	payment, err := server.store.CreatePayment(ctx, arg_payment)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return nil
	} else {
		log.Info().Msgf("Payment %v", payment)
	}
	return nil
}

////////////////////////////////////////////////////////////// Get User

type getUserRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

func (server *Server) getUser(ctx *gin.Context) {
	//Authorization
	_, err := server.authorizeUser(ctx, []string{util.AdminRole, util.ReViewerRole, util.UserRole}) //authPayload
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	var req getUserRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	user, err := server.store.GetUser(ctx, req.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	resp := newUserResponse(user)

	ctx.JSON(http.StatusOK, resp)
}

//////////////////////////////////////////////////////////// Get Users

type listUserRequest struct {
	PageID   int32 `form:"page_id" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=5,max=10"`
}

func (server *Server) listUsers(ctx *gin.Context) {
	//Authorization
	_, err := server.authorizeUser(ctx, []string{util.AdminRole, util.ReViewerRole}) //authPayload
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	var req listUserRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.ListUsersParams{
		Limit:  req.PageSize,
		Offset: (req.PageID - 1) * req.PageSize,
	}

	users, err := server.store.ListUsers(ctx, arg)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, users)
}

////////////////////////////////////////////////////////////////// Update User

type updateUserRequest struct {
	ID         int64                 `form:"id" binding:"required"`
	Name       string                `form:"name" binding:"required,alphanum"`
	LastName   string                `form:"last_name" binding:"required"`
	Email      string                `form:"email" binding:"required,email"`
	Phone      string                `form:"phone"`
	Password   string                `form:"password" binding:"required,min=6"`
	UserType   string                `form:"user_type"`
	Role       string                `form:"role" binding:"required,oneof=admin user reviewer"`
	AvatarPath string                `form:"avatar_path"`
	Avatar     *multipart.FileHeader `form:"avatar"`
}

func (server *Server) updateUser(ctx *gin.Context) {
	var req updateUserRequest
	if err := ctx.ShouldBind(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	if req.Email == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "Must have Email"})
		return
	}

	if req.Password == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "Must have Password"})
		return
	}

	hashedPassword, err := util.HashPassword(req.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	avatar_path := "/assets/images/user/user.png"

	//Avatar
	if req.Avatar != nil {
		ext := filepath.Ext(req.Avatar.Filename)
		avatar_path = "assets/images/user/user_" + fmt.Sprint(req.ID) + ext
		err = ctx.SaveUploadedFile(req.Avatar, avatar_path)
		if err != nil {
			ctx.JSON(http.StatusForbidden, errorResponse(err))
			return
		}
	}

	arg := db.UpdateUserParams{
		ID:             req.ID,
		Name:           req.Name,
		LastName:       req.LastName,
		Email:          req.Email,
		Phone:          req.Phone,
		HashedPassword: hashedPassword,
		UserType:       req.UserType,
		Role:           req.Role,
		AvatarPath:     avatar_path,
	}

	user, err := server.store.UpdateUser(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	//Modify User Role Claims
	cookie, _ := ctx.Cookie("token")
	claims, _ := util.ParseToken(cookie)

	claims = &util.Claims{
		Role:  user.Role,
		Token: claims.Token,
		StandardClaims: jwt.StandardClaims{
			Subject:   user.Email,
			ExpiresAt: claims.ExpiresAt,
		},
	}

	// Build a token for Cookie
	new_token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	var jwtKey = []byte(server.config.TokenSymetricKey)
	tokenString, err := new_token.SignedString(jwtKey)

	if err != nil {
		ctx.JSON(500, gin.H{"error": "could not generate token"})
		return
	}

	//Build Cookie for authorization
	ctx.SetCookie("token", tokenString, int(claims.ExpiresAt), "", "", false, true)

	ctx.JSON(http.StatusOK, user)
}

//////////////////////////////////////////////////////////// Delete User

type deleteUserRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

func (server *Server) deleteUser(ctx *gin.Context) {
	var req deleteUserRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
	}

	err := server.store.DeleteUser(ctx, req.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, "Delete User Ok")
}

////////////////////////////////////////////////////////// Login User

type loginUserRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type loginUserResponse struct {
	SessionID             uuid.UUID    `json:"session_id"`
	AccessToken           string       `json:"accesss_token"`
	AccessTokenExpiresAt  time.Time    `json:"accesss_token_expires_at"`
	RefreshToken          string       `json:"refresh_token"`
	RefreshTokenExpiresAt time.Time    `json:"refresh_token_expires_at"`
	User                  userResponse `json:"user"`
}

func (server *Server) loginUser(ctx *gin.Context) {
	var req loginUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	user, err := server.store.GetUserEmail(ctx, req.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	err = util.CheckPassword(req.Password, user.HashedPassword)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	accessToken, accessPayload, err := server.tokenMaker.CreateToken(
		user.ID,
		user.Email,
		user.Role,
		server.config.AccessTokenDuration,
	)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	refreshToken, refreshPayload, err := server.tokenMaker.CreateToken(
		user.ID,
		user.Name,
		user.Role,
		server.config.RefreshTokenDuration,
	)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	session, err := server.store.CreateSession(ctx, db.CreateSessionParams{
		ID:           refreshPayload.ID,
		UserID:       int32(user.ID),
		Email:        user.Email,
		RefreshToken: refreshToken,
		UserAgent:    ctx.Request.UserAgent(),
		ClientIp:     ctx.ClientIP(),
		IsBloked:     false,
		ExpiredAt:    refreshPayload.ExpiredAt,
	})

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	resp := loginUserResponse{
		SessionID:             session.ID,
		AccessToken:           accessToken,
		AccessTokenExpiresAt:  accessPayload.ExpiredAt,
		RefreshToken:          refreshToken,
		RefreshTokenExpiresAt: refreshPayload.ExpiredAt,
		User:                  newUserResponse(user),
	}

	claims := &util.Claims{
		Role:  user.Role,
		Token: accessToken,
		StandardClaims: jwt.StandardClaims{
			Subject:   user.Email,
			ExpiresAt: accessPayload.ExpiredAt.Unix(),
		},
	}

	// Build a token for Cookie
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	var jwtKey = []byte(server.config.TokenSymetricKey)
	tokenString, err := token.SignedString(jwtKey)

	if err != nil {
		ctx.JSON(500, gin.H{"error": "could not generate token"})
		return
	}

	//Build Cookie for authorization
	ctx.SetCookie("token", tokenString, int(accessPayload.ExpiredAt.Unix()), "", "", false, true)

	ctx.JSON(http.StatusOK, resp)
}

func (server *Server) logoutUser(ctx *gin.Context) {
	ctx.SetCookie("token", "", -1, "", "", false, true)
	ctx.JSON(http.StatusOK, gin.H{"success": "user logged out"})
}
