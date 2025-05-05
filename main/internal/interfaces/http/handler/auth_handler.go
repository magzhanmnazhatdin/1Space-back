package handler

import (
	"context"
	"net/http"
	"strings"

	"firebase.google.com/go/v4/auth"
	"github.com/gin-gonic/gin"
)

// AuthHandler handles authentication requests.
type AuthHandler struct {
	authClient *auth.Client
}

// NewAuthHandler creates a new AuthHandler with injected Firebase Auth client.
func NewAuthHandler(client *auth.Client) *AuthHandler {
	return &AuthHandler{authClient: client}
}

func (h *AuthHandler) Auth(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "No Authorization header"})
		return
	}

	// remove "Bearer " prefix (if any) and trim spaces
	token := strings.TrimSpace(strings.TrimPrefix(authHeader, "Bearer "))
	if token == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid Authorization header"})
		return
	}

	// verify ID token
	decodedToken, err := h.authClient.VerifyIDToken(context.Background(), token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Successful login",
		"uid":     decodedToken.UID,
	})
}
