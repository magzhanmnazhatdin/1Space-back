// handlers.go
package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/gin-gonic/gin"
)

func getClubComputers(c *gin.Context) {
	clubID := c.Param("id")

	// Получаем список компьютеров для клуба
	docs, err := client.Collection("computers").
		Where("club_id", "==", clubID).
		Documents(context.Background()).
		GetAll()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var computers []Computer
	for _, doc := range docs {
		var computer Computer
		doc.DataTo(&computer)
		computer.ID = doc.Ref.ID
		computers = append(computers, computer)
	}

	c.JSON(http.StatusOK, computers)
}

func createBooking(c *gin.Context) {
	uid := c.MustGet("uid").(string)

	var booking struct {
		ClubID    string    `json:"club_id"`
		PCNumber  int       `json:"pc_number"`
		StartTime time.Time `json:"start_time"`
		Hours     int       `json:"hours"`
	}

	if err := c.ShouldBindJSON(&booking); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Получаем информацию о клубе
	clubDoc, err := client.Collection("clubs").Doc(booking.ClubID).Get(context.Background())
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Клуб не найден"})
		return
	}

	var club ComputerClub
	clubDoc.DataTo(&club)

	// Проверяем доступность компьютера
	computerDocs, err := client.Collection("computers").
		Where("club_id", "==", booking.ClubID).
		Where("number", "==", booking.PCNumber).
		Limit(1).
		Documents(context.Background()).
		GetAll()

	if err != nil || len(computerDocs) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Компьютер не найден"})
		return
	}

	var computer Computer
	computerDocs[0].DataTo(&computer)

	if !computer.IsAvailable {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Компьютер уже занят"})
		return
	}

	// Проверяем нет ли пересечений по времени
	bookingsQuery := client.Collection("bookings").
		Where("club_id", "==", booking.ClubID).
		Where("pc_number", "==", booking.PCNumber).
		Where("status", "==", "active").
		Where("end_time", ">", time.Now())

	existingBookings, err := bookingsQuery.Documents(context.Background()).GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if len(existingBookings) > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Компьютер уже забронирован на это время"})
		return
	}

	// Создаем бронирование
	endTime := booking.StartTime.Add(time.Duration(booking.Hours) * time.Hour)
	totalPrice := club.PricePerHour * float64(booking.Hours)

	newBooking := Booking{
		ID:         firestore.NewDocID().ID,
		ClubID:     booking.ClubID,
		UserID:     uid,
		PCNumber:   booking.PCNumber,
		StartTime:  booking.StartTime,
		EndTime:    endTime,
		TotalPrice: totalPrice,
		Status:     "active",
		CreatedAt:  time.Now(),
	}

	_, err = client.Collection("bookings").Doc(newBooking.ID).Set(context.Background(), newBooking)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Обновляем статус компьютера
	_, err = client.Collection("computers").Doc(computer.ID).Update(context.Background(), []firestore.Update{
		{Path: "is_available", Value: false},
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, newBooking)
}

// handlers.go
func getUserBookings(c *gin.Context) {
	uid := c.MustGet("uid").(string)

	// Получаем активные бронирования пользователя
	docs, err := client.Collection("bookings").
		Where("user_id", "==", uid).
		Where("status", "==", "active").
		Where("end_time", ">", time.Now()).
		OrderBy("start_time", firestore.Asc).
		Documents(context.Background()).
		GetAll()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var bookings []Booking
	for _, doc := range docs {
		var booking Booking
		doc.DataTo(&booking)
		booking.ID = doc.Ref.ID

		// Получаем информацию о клубе
		clubDoc, err := client.Collection("clubs").Doc(booking.ClubID).Get(context.Background())
		if err == nil {
			var club ComputerClub
			clubDoc.DataTo(&club)
			booking.ClubName = club.Name // Добавляем имя клуба в ответ
		}

		bookings = append(bookings, booking)
	}

	c.JSON(http.StatusOK, bookings)
}

// handlers.go
func cancelBooking(c *gin.Context) {
	uid := c.MustGet("uid").(string)
	bookingID := c.Param("id")

	// Получаем бронирование
	bookingRef := client.Collection("bookings").Doc(bookingID)
	bookingDoc, err := bookingRef.Get(context.Background())
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Бронирование не найдено"})
		return
	}

	var booking Booking
	bookingDoc.DataTo(&booking)

	// Проверяем, что бронирование принадлежит пользователю
	if booking.UserID != uid {
		c.JSON(http.StatusForbidden, gin.H{"error": "Нельзя отменить чужое бронирование"})
		return
	}

	// Проверяем, что бронирование еще активно
	if booking.Status != "active" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Бронирование уже отменено или завершено"})
		return
	}

	// Проверяем, что не слишком поздно отменять (минимум 1 час до начала)
	if time.Until(booking.StartTime) < time.Hour {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Можно отменить только за час до начала"})
		return
	}

	// Обновляем статус бронирования
	_, err = bookingRef.Update(context.Background(), []firestore.Update{
		{Path: "status", Value: "cancelled"},
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Освобождаем компьютер
	computerDocs, err := client.Collection("computers").
		Where("club_id", "==", booking.ClubID).
		Where("number", "==", booking.PCNumber).
		Limit(1).
		Documents(context.Background()).
		GetAll()

	if err != nil {
		log.Printf("Ошибка поиска компьютера: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при освобождении компьютера"})
		return
	}

	if len(computerDocs) > 0 {
		_, err = computerDocs[0].Ref.Update(context.Background(), []firestore.Update{
			{Path: "is_available", Value: true},
		})

		if err != nil {
			log.Printf("Ошибка обновления статуса компьютера: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при обновлении статуса компьютера"})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "Бронирование успешно отменено"})
}
