package api

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/juliansunn/simple_bank/token"
)

const (
	authorizationHeaderKey  = "authorization"
	authorizationTypeBearer = "bearer"
	authorizationPayloadKey = "authorization_payload"
)

// authMiddleware is a function that returns a Gin middleware function for handling authentication.
//
// It takes a tokenMaker as a parameter, which is responsible for creating and verifying tokens.
//
// The middleware function performs the following steps:
// - Retrieves the authorization header from the request context.
// - Checks if the authorization header is provided. If not, it aborts the request with an unauthorized status and an error response.
// - Parses the authorization header and checks if it has the correct format. If not, it aborts the request with an unauthorized status and an error response.
// - Extracts the authorization type from the parsed header and checks if it is supported. If not, it aborts the request with an unauthorized status and an error response.
// - Extracts the access token from the parsed header.
// - Verifies the access token using the tokenMaker. If the verification fails, it aborts the request with an unauthorized status and an error response.
// - Sets the authorization payload in the request context.
// - Calls the next handler in the chain.
//
// The middleware function does not have any return types.
func authMiddleware(tokenMaker token.Maker) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authorizationHeader := ctx.GetHeader(authorizationHeaderKey)
		if len(authorizationHeader) == 0 {
			err := errors.New("authorization header is not provided")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}
		fields := strings.Fields(authorizationHeader)
		if len(fields) < 2 {
			err := errors.New("invalid authorization header format")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		authorizionType := strings.ToLower(fields[0])
		if authorizionType != authorizationTypeBearer {
			err := fmt.Errorf("unsupported authorization type %s", authorizionType)
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		accessToken := fields[1]
		payload, err := tokenMaker.VerifyToken(accessToken)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return

		}

		ctx.Set(authorizationPayloadKey, payload)
		ctx.Next()

	}
}
