package mocks

import (
	"context"
	"net/http"

	"firebase.google.com/go/v4/auth"
	"github.com/gin-gonic/gin"
)

type FakeAuthClient struct{}

func NewFakeAuthClient() *FakeAuthClient {
	return &FakeAuthClient{}
}

func (f *FakeAuthClient) VerifyIDToken(ctx context.Context, idToken string) (*auth.Token, error) {
	return &auth.Token{
		UID:    "test-user",
		Claims: map[string]interface{}{"role": "user"},
	}, nil
}

// Stub AuthHandler
type StubAuthHandler struct{}

func NewStubAuthHandler() *StubAuthHandler {
	return &StubAuthHandler{}
}

func (h *StubAuthHandler) Auth(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "stub login",
		"uid":     "test-user",
	})
}

// Stub PaymentHandler
type StubPaymentHandler struct{}

func NewStubPaymentHandler() *StubPaymentHandler {
	return &StubPaymentHandler{}
}

func (h *StubPaymentHandler) CreateIntent(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"clientSecret": "test-secret"})
}

func (h *StubPaymentHandler) Webhook(c *gin.Context) {
	c.Status(http.StatusOK)
}
