package http_test

import (
	"bytes"
	"encoding/json"
	interfaceHttp "main/internal/interfaces/http"
	"main/internal/interfaces/http/handler"
	"main/test/mocks"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func setupRouter() *gin.Engine {
	clubUC := new(mocks.MockClubUC)
	clubH := handler.NewClubHandler(clubUC)

	compUC := new(mocks.MockComputerUC)
	compH := handler.NewComputerHandler(compUC, clubUC)

	bookUC := new(mocks.MockBookingUC)
	bookH := handler.NewBookingHandler(bookUC, clubUC)

	router := interfaceHttp.NewRouter(clubH, compH, bookH, nil, nil, nil, nil)
	return router
}

func TestPublicGetClubs(t *testing.T) {
	router := setupRouter()

	req, _ := http.NewRequest("GET", "/clubs", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.True(t, resp.Code == http.StatusOK || resp.Code == http.StatusUnauthorized)
}

func TestProtectedPostBooking_Unauthorized(t *testing.T) {
	router := setupRouter()

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
