package handler_test

import (
	"github.com/stretchr/testify/mock"
	"main/internal/domain/entities"
	"main/internal/interfaces/http/handler"
	"main/test/mocks"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestGetAllComputers(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockCompUC := new(mocks.MockComputerUC)
	mockClubUC := new(mocks.MockClubUC)
	h := handler.NewComputerHandler(mockCompUC, mockClubUC)

	comps := []*entities.Computer{{ID: "c1", ClubID: "club1"}}
	mockCompUC.On("GetAll", mock.Anything).Return(comps, nil).Once()

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	req, _ := http.NewRequest("GET", "/computers", nil)
	c.Request = req

	h.GetAllComputers(c)

	assert.Equal(t, http.StatusOK, w.Code)
	mockCompUC.AssertExpectations(t)
}

func TestGetClubComputers(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockCompUC := new(mocks.MockComputerUC)
	mockClubUC := new(mocks.MockClubUC)
	h := handler.NewComputerHandler(mockCompUC, mockClubUC)

	comps := []*entities.Computer{{ID: "c1", ClubID: "club1"}}
	mockCompUC.On("GetByClub", mock.Anything, "club1").Return(comps, nil).Once()

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	req, _ := http.NewRequest("GET", "/clubs/club1/computers", nil)
	c.Request = req
	c.Params = gin.Params{{Key: "id", Value: "club1"}}

	h.GetClubComputers(c)

	assert.Equal(t, http.StatusOK, w.Code)
	mockCompUC.AssertExpectations(t)
}
