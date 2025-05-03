package handler

import (
	"firebase.google.com/go/v4/auth"
	"github.com/gin-gonic/gin"
	"net/http"
)

// AuthHandler handles authentication requests.
type AuthHandler struct {
	client *auth.Client
}

// NewAuthHandler creates a new AuthHandler with injected Firebase Auth client.
func NewAuthHandler(client *auth.Client) *AuthHandler {
	return &AuthHandler{client: client}
}

func (h *AuthHandler) Auth(c *gin.Context) {
	var req struct {
		IDToken string `json:"id_token"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	token, err := h.client.VerifyIDToken(c.Request.Context(), req.IDToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	c.Header("X-User-ID", token.UID)
	c.JSON(http.StatusOK, gin.H{"uid": token.UID})
}
