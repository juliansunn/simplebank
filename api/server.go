package api

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	db "github.com/juliansunn/simple_bank/db/sqlc"
	"github.com/juliansunn/simple_bank/token"
	"github.com/juliansunn/simple_bank/util"
)

type Server struct {
	config     util.Config
	store      db.Store
	tokenMaker token.Maker
	router     *gin.Engine
}

// NewServer creates a new Server instance with the given store.
//
// Parameters:
// - store: a pointer to a db.Store object.
//
// Returns:
// - a pointer to the newly created Server object.
func NewServer(config util.Config, store db.Store) (*Server, error) {
	// tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	// replace the following line with the line above to use JWT instead of Paseto
	tokenMaker, err := token.NewJWTMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}
	server := &Server{config: config, store: store, tokenMaker: tokenMaker}

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("currency", validCurrency)
	}
	server.setupRouter()
	return server, nil
}

func (server *Server) setupRouter() {
	router := gin.Default()
	router.POST("/users", server.createUser)
	router.POST("/users/login", server.loginUser)

	// middleware
	authRoutes := router.Group("/").Use(authMiddleware(server.tokenMaker))
	authRoutes.POST("/accounts", server.createAccount)
	authRoutes.GET("/accounts/:id", server.getAccount)
	authRoutes.GET("/accounts", server.listAccounts)
	authRoutes.POST("/transfers", server.createTransfer)
	server.router = router
}

// Start starts the server on the specified address.
//
// Parameters:
// - address: the address to listen on.
//
// Returns:
// - error: an error if the server failed to start.
func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

// errorResponse generates an error response in the form of a gin.H map.
//
// The function takes an error as a parameter and returns a gin.H map
// with a single key "error" and the error message as its value.
func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
