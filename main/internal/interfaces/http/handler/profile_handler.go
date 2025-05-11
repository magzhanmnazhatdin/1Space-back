// internal/interfaces/http/handler/profile_handler.go
package handler

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"net/http"

	"github.com/gin-gonic/gin"
	"main/internal/application/usecase"
	"main/internal/domain/entities"
)

type ProfileHandler struct {
	uc usecase.ProfileUseCase
}

func NewProfileHandler(uc usecase.ProfileUseCase) *ProfileHandler {
	return &ProfileHandler{uc: uc}
}

// GET /profile
func (h *ProfileHandler) GetProfile(c *gin.Context) {
	uid := c.GetString("uid")

	p, err := h.uc.GetProfile(c.Request.Context(), uid)
	if err != nil {
		// если firestore вернул NotFound — просто создаём пустой профиль
		if st, ok := status.FromError(err); ok && st.Code() == codes.NotFound {
			p = &entities.Profile{
				UserID:      uid,
				ProfileName: "",
				AvatarURL:   "",
				Bio:         "",
			}
		} else {
			// настоящая ошибка
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	c.JSON(http.StatusOK, p)
}

// PUT /profile
func (h *ProfileHandler) UpdateProfile(c *gin.Context) {
	uid := c.GetString("uid")
	var req struct {
		ProfileName string `json:"profile_name" binding:"required"`
		AvatarURL   string `json:"avatar_url,omitempty"`
		Bio         string `json:"bio,omitempty"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	p := &entities.Profile{
		UserID:      uid,
		ProfileName: req.ProfileName,
		AvatarURL:   req.AvatarURL,
		Bio:         req.Bio,
	}
	if err := h.uc.UpdateProfile(c.Request.Context(), p); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}
