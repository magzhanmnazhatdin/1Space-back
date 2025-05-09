package handler

import (
	"main/internal/config"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/stripe/stripe-go/v72/webhook"
	"main/internal/application/usecase"
)

type PaymentHandler struct {
	uc usecase.PaymentUseCase
}

func NewPaymentHandler(uc usecase.PaymentUseCase) *PaymentHandler {
	return &PaymentHandler{uc: uc}
}

func (h *PaymentHandler) CreateIntent(c *gin.Context) {
	var req struct {
		Amount   int64  `json:"amount"`
		Currency string `json:"currency"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	clientSecret, err := h.uc.CreateIntent(c.Request.Context(), req.Amount, req.Currency)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"clientSecret": clientSecret})
}

func (h *PaymentHandler) Webhook(c *gin.Context) {
	payload, err := c.GetRawData()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	sig := c.GetHeader("Stripe-Signature")
	event, err := webhook.ConstructEvent(payload, sig, config.Cfg.Stripe.WebhookSecret)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid signature"})
		return
	}
	switch event.Type {
	case "payment_intent.succeeded":
		// handle success
	case "payment_intent.payment_failed":
		// handle failure
	}
	c.Status(http.StatusOK)
}
