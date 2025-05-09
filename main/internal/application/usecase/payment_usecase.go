package usecase

import (
	"context"
	"errors"

	"main/internal/infrastructure/stripeclient"
)

type PaymentUseCase interface {
	CreateIntent(ctx context.Context, amount int64, currency string) (string, error)
}

type paymentInteractor struct {
}

func NewPaymentUseCase() PaymentUseCase {
	return &paymentInteractor{}
}

func (u *paymentInteractor) CreateIntent(ctx context.Context, amount int64, currency string) (string, error) {
	if amount <= 0 {
		return "", errors.New("invalid amount")
	}
	pi, err := stripeclient.CreatePaymentIntent(amount, currency, nil)
	if err != nil {
		return "", err
	}
	return pi.ClientSecret, nil
}
