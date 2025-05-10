package handler

import (
	"net/http"

	"firebase.google.com/go/v4/auth"
	"github.com/gin-gonic/gin"
)

// UserHandler управляет пользователями
type UserHandler struct {
	authClient *auth.Client
}

func NewUserHandler(authClient *auth.Client) *UserHandler {
	return &UserHandler{authClient: authClient}
}

// ChangeRole — только для админа: меняет custom claim "role"
func (h *UserHandler) ChangeRole(c *gin.Context) {
	uid := c.Param("id")
	var req struct {
		Role string `json:"role"` // ожидаем "admin", "manager" или "user"
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// ставим новый claim
	err := h.authClient.SetCustomUserClaims(c.Request.Context(), uid, map[string]interface{}{
		"role": req.Role,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}
