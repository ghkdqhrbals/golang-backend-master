package api

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/ghkdqhrbals/golang-backend-master/token"
	"github.com/gin-gonic/gin"
)

const (
	authortizeationHeaderKey = "authorization"
	authorizationTypeBearer  = "bearer"
	authorizationPayloadKey  = "authorization_payload"
)

func authMiddleware(tokenMaker token.Maker) gin.HandlerFunc {
	// Header 에서 Token 가져옴
	return func(ctx *gin.Context) {
		authortizeationHeader := ctx.GetHeader(authortizeationHeaderKey)
		if len(authortizeationHeader) == 0 {
			err := errors.New("authorization header is not provided")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		fields := strings.Fields(authortizeationHeader)
		if len(fields) < 2 {
			err := errors.New("Invalid authorization header format")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		authorizationType := strings.ToLower(fields[0])
		if authorizationType != authorizationTypeBearer {
			err := fmt.Errorf("unsupported authorization type %s", authorizationType)
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		accessToken := fields[1]
		payload, err := tokenMaker.VerifyToken(accessToken)

		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		ctx.Set(authorizationPayloadKey, payload) // ctx[authorizationPayloadKey] = payload
		ctx.Next()

	}
}
