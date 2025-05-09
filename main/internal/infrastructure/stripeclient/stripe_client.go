package stripeclient

import (
	"github.com/stripe/stripe-go/v72"
	"github.com/stripe/stripe-go/v72/paymentintent"
)

// Init sets Stripe API key
func Init(key string) {
	stripe.Key = key
}

// CreatePaymentIntent creates a new PaymentIntent
func CreatePaymentIntent(amount int64, currency string, metadata map[string]string) (*stripe.PaymentIntent, error) {
	params := &stripe.PaymentIntentParams{
		Params: stripe.Params{
			Metadata: metadata,
		},
		Amount:   stripe.Int64(amount),
		Currency: stripe.String(currency),
	}
	return paymentintent.New(params)
}
