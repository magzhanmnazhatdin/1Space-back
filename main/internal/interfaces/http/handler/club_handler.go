package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"main/internal/application/usecase"
	"main/internal/domain/entities"
)

type ClubHandler struct {
	uc usecase.ClubUseCase
}

func NewClubHandler(uc usecase.ClubUseCase) *ClubHandler {
	return &ClubHandler{uc: uc}
}

func (h *ClubHandler) GetAllClubs(c *gin.Context) {
	clubs, err := h.uc.GetAll(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, clubs)
}

func (h *ClubHandler) GetClubByID(c *gin.Context) {
	id := c.Param("id")
	club, err := h.uc.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	c.JSON(http.StatusOK, club)
}

func (h *ClubHandler) CreateClub(c *gin.Context) {
	uid := c.GetString("uid")
	role := c.GetString("role")

	var req struct {
		Name         string  `json:"name" binding:"required"`
		Address      string  `json:"address" binding:"required"`
		PricePerHour float64 `json:"price_per_hour" binding:"required"`
		ManagerID    string  `json:"manager_id,omitempty"` // теперь опционально
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	club := entities.Club{
		Name:         req.Name,
		Address:      req.Address,
		PricePerHour: req.PricePerHour,
	}

	switch role {
	case "manager":
		club.ManagerID = uid

	case "admin":
		if req.ManagerID != "" {
			club.ManagerID = req.ManagerID
		} else {
			club.ManagerID = uid
		}

	default:
		c.JSON(http.StatusForbidden, gin.H{"error": "access denied"})
		return
	}

	if err := h.uc.Create(c.Request.Context(), &club); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, club)
}

func (h *ClubHandler) UpdateClub(c *gin.Context) {
	id := c.Param("id")
	orig, err := h.uc.GetByID(c.Request.Context(), id) // CHANGED: сохраняем старый ManagerID
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}

	var in entities.Club
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	in.ID = id
	in.ManagerID = orig.ManagerID // CHANGED: не даем менять владельца

	if err := h.uc.Update(c.Request.Context(), &in); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *ClubHandler) DeleteClub(c *gin.Context) {
	uid := c.GetString("uid")
	role := c.GetString("role")
	id := c.Param("id")

	if role == "manager" {
		club, err := h.uc.GetByID(c.Request.Context(), id)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
			return
		}
		if club.ManagerID != uid {
			c.JSON(http.StatusForbidden, gin.H{"error": "access denied"})
			return
		}
	}

	if role != "manager" && role != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "access denied"})
		return
	}

	if err := h.uc.Delete(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}
