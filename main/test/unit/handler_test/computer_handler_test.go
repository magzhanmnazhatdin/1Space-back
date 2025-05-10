package handler_test

import (
	"context"
	"encoding/json"
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
	compUC := new(mocks.MockComputerUC)
	clubUC := new(mocks.MockClubUC)
	h := handler.NewComputerHandler(compUC, clubUC)

	comps := []*entities.Computer{
		{ID: "1", PCNumber: 1},
		{ID: "2", PCNumber: 2},
	}
	compUC.On("GetAll", context.Background()).Return(comps, nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	h.GetAllComputers(c)

	assert.Equal(t, http.StatusOK, w.Code)
	var result []entities.Computer
	json.Unmarshal(w.Body.Bytes(), &result)
	assert.Len(t, result, 2)
}

func TestGetClubComputers(t *testing.T) {
	compUC := new(mocks.MockComputerUC)
	clubUC := new(mocks.MockClubUC)
	h := handler.NewComputerHandler(compUC, clubUC)

	comps := []*entities.Computer{{ID: "1", PCNumber: 1}}
	compUC.On("GetByClub", context.Background(), "club1").Return(comps, nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{{Key: "id", Value: "club1"}}

	h.GetClubComputers(c)

	assert.Equal(t, http.StatusOK, w.Code)
	var result []entities.Computer
	json.Unmarshal(w.Body.Bytes(), &result)
	assert.Len(t, result, 1)
}
