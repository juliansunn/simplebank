package gapi

import (
	"fmt"

	db "github.com/juliansunn/simple_bank/db/sqlc"
	"github.com/juliansunn/simple_bank/pb"
	"github.com/juliansunn/simple_bank/token"
	"github.com/juliansunn/simple_bank/util"
)

// Server serves gRPC requests for our banking service.
type Server struct {
	pb.UnimplementedSimpleBankServer
	config     util.Config
	store      db.Store
	tokenMaker token.Maker
}

// NewServer creates a new gRPC server.
func NewServer(config util.Config, store db.Store) (*Server, error) {
	// tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	// replace the following line with the line above to use JWT instead of Paseto
	tokenMaker, err := token.NewJWTMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}
	server := &Server{config: config, store: store, tokenMaker: tokenMaker}

	return server, nil
}
