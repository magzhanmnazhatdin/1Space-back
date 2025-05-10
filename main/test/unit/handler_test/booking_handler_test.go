package handler_test

import (
	"bytes"
	"encoding/json"
	"github.com/stretchr/testify/mock"
	"main/internal/domain/entities"
	"main/internal/interfaces/http/handler"
	"main/test/mocks"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestGetUserBookings(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockUC := new(mocks.MockBookingUC)
	h := handler.NewBookingHandler(mockUC, nil)

	bookings := []*entities.Booking{{ID: "1"}, {ID: "2"}}
	mockUC.On("GetByUser", mock.Anything, "user123").Return(bookings, nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set("uid", "user123")

	h.GetUserBookings(c)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestCreateBooking_Success(t *testing.T) {
	mockBooking := new(mocks.MockBookingUC)
	mockClub := new(mocks.MockClubUC)
	h := handler.NewBookingHandler(mockBooking, mockClub)

	club := &entities.Club{ID: "club1", PricePerHour: 10.0}
	mockClub.On("GetByID", mock.Anything, "club1").Return(club, nil)
	mockBooking.On("Create", mock.Anything, mock.AnythingOfType("*entities.Booking")).Return(nil)

	body := map[string]interface{}{
		"club_id":    "club1",
		"pc_number":  1,
		"start_time": time.Now().Format(time.RFC3339),
		"end_time":   time.Now().Add(time.Hour).Format(time.RFC3339),
	}
	jsonBody, _ := json.Marshal(body)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("POST", "/bookings", bytes.NewReader(jsonBody))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Set("uid", "user123")

	h.CreateBooking(c)

	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestCancelBooking_Success(t *testing.T) {
	mockUC := new(mocks.MockBookingUC)
	h := handler.NewBookingHandler(mockUC, nil)

	mockUC.On("Cancel", mock.Anything, "booking123").Return(nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{{Key: "id", Value: "booking123"}}

	h.CancelBooking(c)

	assert.Equal(t, http.StatusNoContent, w.Code)
}
