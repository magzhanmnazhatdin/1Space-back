package handler

import (
	"github.com/gin-gonic/gin"
	"main/internal/application/usecase"
	"main/internal/domain/entities"
	"net/http"
)

// ClubHandler handles HTTP requests for clubs.
type ClubHandler struct {
	uc usecase.ClubUseCase
}

// NewClubHandler creates a new ClubHandler with injected use case.
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
	var in entities.Club
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.uc.Create(c.Request.Context(), &in); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, in)
}

func (h *ClubHandler) UpdateClub(c *gin.Context) {
	id := c.Param("id")
	var in entities.Club
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	in.ID = id
	if err := h.uc.Update(c.Request.Context(), &in); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *ClubHandler) DeleteClub(c *gin.Context) {
	id := c.Param("id")
	if err := h.uc.Delete(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}
