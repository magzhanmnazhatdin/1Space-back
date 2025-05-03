package middleware

import (
	"firebase.google.com/go/v4/auth"
	"github.com/gin-gonic/gin"
	"net/http"
)

func AuthMiddleware(authClient *auth.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		idToken := c.GetHeader("Authorization")
		if idToken == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "no token"})
			return
		}
		token, err := authClient.VerifyIDToken(c.Request.Context(), idToken)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}
		// прокинуть UID дальше
		c.Request.Header.Set("X-User-ID", token.UID)
		c.Next()
	}
}
