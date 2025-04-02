package middleware

import (
	"daterrr/pkg/auth/tokengen"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func errorResponse(err error) gin.H {
	return gin.H{
		"error": err.Error(),
	}
}

/* Middleware for handling authentication of request (pretty obvi innit?) */
func AuthMiddleware(tok tokengen.Maker) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("authorization")
		// If no auth header was provided with the request
		if len(authHeader) < 0 {
			err := errors.New("No authorization header was provided")
			c.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		fields := strings.Fields(authHeader)
		if len(fields) < 2 {
			err := errors.New("Invalid authorization header format")
			c.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		authType := strings.ToLower(fields[0])
		if authType != "Bearer" {
			err := fmt.Errorf("Unsupported auth type %s", authType)
			c.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		accessToken := fields[1]
		payload, err := tok.VerifyToken(accessToken)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		c.Set("userID", payload.UserID)
		c.Next()

	}

}
