package handler

import (
	"github.com/gin-gonic/gin"
	"main/internal/application/usecase"
	"main/internal/domain/entities"
	"net/http"
)

// ComputerHandler handles HTTP requests for computers.
type ComputerHandler struct {
	uc usecase.ComputerUseCase
}

// NewComputerHandler creates a new ComputerHandler with injected use case.
func NewComputerHandler(uc usecase.ComputerUseCase) *ComputerHandler {
	return &ComputerHandler{uc: uc}
}

func (h *ComputerHandler) GetAllComputers(c *gin.Context) {
	comps, err := h.uc.GetAll(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, comps)
}

func (h *ComputerHandler) GetClubComputers(c *gin.Context) {
	clubID := c.Param("id")
	comps, err := h.uc.GetByClub(c.Request.Context(), clubID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, comps)
}

func (h *ComputerHandler) CreateComputerList(c *gin.Context) {
	clubID := c.Param("id")
	var list []entities.Computer
	if err := c.ShouldBindJSON(&list); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	for i := range list {
		list[i].ClubID = clubID
		if err := h.uc.Create(c.Request.Context(), &list[i]); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}
	c.JSON(http.StatusCreated, list)
}
