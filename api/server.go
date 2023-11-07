package api

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	db "github.com/juliansunn/simple_bank/db/sqlc"
)

type Server struct {
	store  db.Store
	router *gin.Engine
}

// NewServer creates a new Server instance with the given store.
//
// Parameters:
// - store: a pointer to a db.Store object.
//
// Returns:
// - a pointer to the newly created Server object.
func NewServer(store db.Store) *Server {
	server := &Server{store: store}
	router := gin.Default()

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("currency", validCurrency)
	}

	router.POST("/users", server.createUser)
	router.POST("/accounts", server.createAccount)
	router.GET("/accounts/:id", server.getAccount)
	router.GET("/accounts", server.listAccounts)

	router.POST("/transfers", server.createTransfer)

	server.router = router
	return server
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
