package handler_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
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
	mockUC.On("GetAll", context.Background()).Return(clubs, nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	h.GetAllClubs(c)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestGetClubByID_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockUC := new(mocks.MockClubUC)
	h := handler.NewClubHandler(mockUC)

	club := &entities.Club{ID: "123", Name: "Test Club"}
	mockUC.On("GetByID", context.Background(), "123").Return(club, nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{{Key: "id", Value: "123"}}

	h.GetClubByID(c)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestGetClubByID_NotFound(t *testing.T) {
	mockUC := new(mocks.MockClubUC)
	h := handler.NewClubHandler(mockUC)

	mockUC.On("GetByID", context.Background(), "999").Return(nil, errors.New("not found"))

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{{Key: "id", Value: "999"}}

	h.GetClubByID(c)

	assert.Equal(t, http.StatusNotFound, w.Code)
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
	mockUC := new(mocks.MockClubUC)
	h := handler.NewClubHandler(mockUC)

	mockUC.On("GetByID", context.Background(), "1").Return(&entities.Club{ID: "1", ManagerID: "manager1"}, nil)
	mockUC.On("Delete", context.Background(), "1").Return(nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{{Key: "id", Value: "1"}}
	c.Set("uid", "admin1")
	c.Set("role", "admin")

	h.DeleteClub(c)

	assert.Equal(t, http.StatusNoContent, w.Code)
}
