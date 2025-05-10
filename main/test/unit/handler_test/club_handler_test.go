package handler_test

import (
	"bytes"
	"encoding/json"
	"errors"
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

func TestGetAllClubs(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockUC := new(mocks.MockClubUC)
	h := handler.NewClubHandler(mockUC)

	clubs := []*entities.Club{{ID: "1", Name: "Club1"}, {ID: "2", Name: "Club2"}}
	mockUC.On("GetAll", mock.Anything).Return(clubs, nil).Once()

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	req, _ := http.NewRequest("GET", "/clubs", nil)
	c.Request = req

	h.GetAllClubs(c)

	assert.Equal(t, http.StatusOK, w.Code)
	mockUC.AssertExpectations(t)
}

func TestGetClubByID_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockUC := new(mocks.MockClubUC)
	h := handler.NewClubHandler(mockUC)

	club := &entities.Club{ID: "1", Name: "Club1"}
	mockUC.On("GetByID", mock.Anything, "1").Return(club, nil).Once()

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	req, _ := http.NewRequest("GET", "/clubs/1", nil)
	c.Request = req
	c.Params = gin.Params{{Key: "id", Value: "1"}}

	h.GetClubByID(c)

	assert.Equal(t, http.StatusOK, w.Code)
	mockUC.AssertExpectations(t)
}

func TestGetClubByID_NotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockUC := new(mocks.MockClubUC)
	h := handler.NewClubHandler(mockUC)

	// возвращаем явный *entities.Club = nil, чтобы mock не падал
	var nilClub *entities.Club = nil
	mockUC.On("GetByID", mock.Anything, "999").Return(nilClub, errors.New("not found")).Once()

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	req, _ := http.NewRequest("GET", "/clubs/999", nil)
	c.Request = req
	c.Params = gin.Params{{Key: "id", Value: "999"}}

	h.GetClubByID(c)

	assert.Equal(t, http.StatusNotFound, w.Code)
	mockUC.AssertExpectations(t)
}

func TestCreateClub_InvalidRole(t *testing.T) {
	mockUC := new(mocks.MockClubUC)
	h := handler.NewClubHandler(mockUC)

	club := entities.Club{Name: "Unauthorized"}
	payload, _ := json.Marshal(club)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("POST", "/clubs", bytes.NewReader(payload))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Set("uid", "user1")
	c.Set("role", "user")

	h.CreateClub(c)

	assert.Equal(t, http.StatusForbidden, w.Code)
}

func TestDeleteClub_Success_Admin(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockUC := new(mocks.MockClubUC)
	h := handler.NewClubHandler(mockUC)

	// для admin роли GetByID не вызываем, сразу Delete
	mockUC.On("Delete", mock.Anything, "club123").Return(nil).Once()

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	req, _ := http.NewRequest("DELETE", "/manager/clubs/club123", nil)
	c.Request = req
	c.Params = gin.Params{{Key: "id", Value: "club123"}}
	c.Set("uid", "adminUser")
	c.Set("role", "admin")

	h.DeleteClub(c)

	assert.Equal(t, http.StatusNoContent, w.Code)
	mockUC.AssertExpectations(t)
}
