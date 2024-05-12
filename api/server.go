package api

import (
	"fmt"
	"net/http"

	"github.com/laugart7272/inscrips/token"
	"github.com/laugart7272/inscrips/worker"
	cors "github.com/rs/cors/wrapper/gin"
	"github.com/spf13/viper"

	"github.com/gin-gonic/gin"
	db "github.com/laugart7272/inscrips/db/sqlc"
	"github.com/laugart7272/inscrips/util"
)

// Serves HTTP requests for our banking services
type Server struct {
	config          util.Config
	store           *db.Store
	router          *gin.Engine
	tokenMaker      token.Maker
	taskDistributor worker.TaskDistributor
}

// Create a new HTTP server and setup routing
func NewServer(config util.Config, store *db.Store, taskDistributor worker.TaskDistributor) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %v", err)
	}

	server := &Server{
		config:          config,
		store:           store,
		tokenMaker:      tokenMaker,
		taskDistributor: taskDistributor,
	}

	server.setupRouter()

	return server, nil
}

func (server *Server) setupRouter() {
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	//CORS
	crs := cors.New(cors.Options{
		AllowedOrigins:   viper.GetStringSlice("CORSHTTP"),
		AllowedMethods:   []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete, http.MethodPatch, http.MethodHead},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
	})

	router.Use(crs)

	router.Static("/assets", "./assets")

	router.POST("/register", server.createUser)
	router.POST("/login", server.loginUser)
	router.POST("/logout", server.logoutUser)
	router.POST("/tokens/renew_access", server.renewAccessToken)

	//Email
	router.GET("/verify_email", server.verifyEmail)

	authRoutes := router.Group("/").Use(authMiddleware(server.tokenMaker))

	//Users
	authRoutes.GET("/user/get/:id", server.getUser)
	authRoutes.GET("/users/list", server.listUsers)
	authRoutes.PUT("/user/update", server.updateUser)
	authRoutes.DELETE("/user/delete/:id", server.deleteUser)

	//Cientific Work
	authRoutes.POST("/cientific_work/register", server.createCientificWork)
	authRoutes.GET("/cientific_work/get/:id", server.getCientificWork)
	authRoutes.GET("/cientifics_works/by_author/list", server.listCientificsWorksByAuthorId)
	authRoutes.GET("/cientifics_works/list", server.listCientificWorks)
	authRoutes.PUT("/cientific_work/update", server.updateCientificWork)
	authRoutes.DELETE("/cientific_work/delete/:id", server.deleteCientificWork)

	//Inscriptions
	authRoutes.POST("/inscriptions", server.createInscription)
	authRoutes.GET("/inscription/:id", server.getInscription)
	authRoutes.GET("/inscriptions/", server.listInscriptions)
	authRoutes.GET("/inscriptions/details/", server.listDetailsInscriptions)
	authRoutes.PUT("/inscription/update", server.updateInscription)
	authRoutes.DELETE("/inscription/:id", server.deleteInscription)

	//Payment
	authRoutes.POST("/payment/register", server.createPayment)
	authRoutes.GET("/payment/get/:id", server.getPayment)
	authRoutes.GET("/payments/list", server.listPayments)
	authRoutes.GET("/payments/details/list", server.listDetailsPayments)
	authRoutes.PUT("/payment/update", server.updatePayment)
	authRoutes.DELETE("/payment/delete/:id", server.deletePayment)

	//Payment Amount
	authRoutes.GET("/payment_amount/get", server.getPaymentAmount)

	//Payment Rates
	authRoutes.GET("/payment_rates/details/list", server.listDetailsPaymentRates)

	server.router = router
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
