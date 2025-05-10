package handler

import (
	"github.com/gin-gonic/gin"
	"main/internal/application/usecase"
	"main/internal/domain/entities"
	"net/http"
)

// ComputerHandler handles HTTP requests for computers.
type ComputerHandler struct {
	uc     usecase.ComputerUseCase
	clubUC usecase.ClubUseCase
}

// NewComputerHandler creates a new ComputerHandler with injected use case.
func NewComputerHandler(uc usecase.ComputerUseCase, clubUC usecase.ClubUseCase) *ComputerHandler {
	return &ComputerHandler{uc: uc, clubUC: clubUC}
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
	uidIf, _ := c.Get("uid")
	userID := uidIf.(string)
	roleIf, _ := c.Get("role")
	role := roleIf.(string)

	clubID := c.Param("id")
	// проверяем права: если manager — то только свой клуб
	if role == "manager" {
		club, err := h.clubUC.GetByID(c.Request.Context(), clubID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid club ID"})
			return
		}
		if club.ManagerID != userID {
			c.JSON(http.StatusForbidden, gin.H{"error": "access denied"})
			return
		}
	}

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

// UPDATE
func (h *ComputerHandler) UpdateComputer(c *gin.Context) {
	userID, _ := c.Get("uid")
	role, _ := c.Get("role")
	cid := c.Param("id")

	// fetch existing computer
	comp, err := h.uc.GetByID(c.Request.Context(), cid)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "computer not found"})
		return
	}
	// если manager — проверяем, что он владелец клуба
	if role == "manager" {
		club, _ := h.clubUC.GetByID(c.Request.Context(), comp.ClubID)
		if club.ManagerID != userID {
			c.JSON(http.StatusForbidden, gin.H{"error": "access denied"})
			return
		}
	}

	var in entities.Computer
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	in.ID = cid
	in.ClubID = comp.ClubID // не даём менять привязку к клубу

	if err := h.uc.Update(c.Request.Context(), &in); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, in)
}

// DELETE
func (h *ComputerHandler) DeleteComputer(c *gin.Context) {
	userID, _ := c.Get("uid")
	role, _ := c.Get("role")
	cid := c.Param("id")

	comp, err := h.uc.GetByID(c.Request.Context(), cid)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "computer not found"})
		return
	}
	if role == "manager" {
		club, _ := h.clubUC.GetByID(c.Request.Context(), comp.ClubID)
		if club.ManagerID != userID {
			c.JSON(http.StatusForbidden, gin.H{"error": "access denied"})
			return
		}
	}

	if err := h.uc.Delete(c.Request.Context(), cid); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}
