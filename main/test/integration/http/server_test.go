package http_test

import (
	"bytes"
	"encoding/json"
	"main/internal/domain/entities"
	interfaceHttp "main/internal/interfaces/http"
	"main/internal/interfaces/http/handler"
	"main/test/mocks"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func setupRouter() *gin.Engine {
	clubUC := new(mocks.MockClubUC)
	clubUC.On("GetAll", mock.Anything).Return([]*entities.Club{}, nil)
	clubUC.On("GetByID", mock.Anything, "testclub").Return(&entities.Club{ID: "testclub", PricePerHour: 10}, nil)
	clubH := handler.NewClubHandler(clubUC)

	compUC := new(mocks.MockComputerUC)
	compH := handler.NewComputerHandler(compUC, clubUC)

	bookUC := new(mocks.MockBookingUC)
	bookUC.On("Create", mock.Anything, mock.Anything).Return(nil)
	bookH := handler.NewBookingHandler(bookUC, clubUC)

	// router без auth/payment/user
	r := interfaceHttp.NewRouter(clubH, compH, bookH, nil, nil, nil, nil)

	// заменим middleware на фейковый, всегда возвращающий 401
	r.Use(mocks.FakeAuthMiddlewareAlways401())

	return r
}

func TestPublicGetClubs(t *testing.T) {
	router := setupRouter()

	req, _ := http.NewRequest("GET", "/clubs", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
}

func TestProtectedPostBooking_Unauthorized(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.Default()
	router.POST("/bookings", mocks.FakeAuthMiddlewareAlways401(), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "should not reach"})
	})

	body := map[string]interface{}{
		"club_id":    "testclub",
		"pc_number":  5,
		"start_time": "2025-05-11T18:00:00Z",
		"end_time":   "2025-05-11T20:00:00Z",
	}
	jsonBody, _ := json.Marshal(body)

	req, _ := http.NewRequest("POST", "/bookings", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}
