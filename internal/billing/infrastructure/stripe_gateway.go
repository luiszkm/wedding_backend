// file: internal/billing/infrastructure/stripe_gateway.go
package infrastructure

import (
	"context"
	"fmt"
	"time"

	"github.com/luiszkm/wedding_backend/internal/billing/domain" // Ajuste o path
	"github.com/stripe/stripe-go/sub"
	"github.com/stripe/stripe-go/v82"
	"github.com/stripe/stripe-go/v82/client"
)

type StripeGateway struct {
	client *client.API
}

func NewStripeGateway(secretKey string) domain.PaymentGateway {
	sc := &client.API{}
	sc.Init(secretKey, nil)
	return &StripeGateway{client: sc}
}
func (sg *StripeGateway) CriarSessaoCheckout(ctx context.Context, assinatura *domain.Assinatura, plano *domain.Plano) (string, error) {
	// A lógica que estava no serviço agora vive aqui.
	params := &stripe.CheckoutSessionParams{
		Mode: stripe.String(string(stripe.CheckoutSessionModeSubscription)),
		LineItems: []*stripe.CheckoutSessionLineItemParams{
			{
				Price:    stripe.String(plano.IDStripePrice()), // O Plano agora precisa ter o ID de preço da Stripe
				Quantity: stripe.Int64(1),
			},
		},
		ClientReferenceID: stripe.String(assinatura.ID().String()),
		SuccessURL:        stripe.String("http://localhost:8080/sucesso?session_id={CHECKOUT_SESSION_ID}"),
		CancelURL:         stripe.String("http://localhost:8080/cancelado"),
	}

	sess, err := sg.client.CheckoutSessions.New(params)
	if err != nil {
		return "", fmt.Errorf("infra: falha ao criar a sessão de checkout da stripe: %w", err)
	}

	return sess.URL, nil

}

func (sg *StripeGateway) GetSubscriptionDetails(ctx context.Context, stripeSubscriptionID string) (*domain.SubscriptionDetails, error) {
	s, err := sub.Get(stripeSubscriptionID, nil)
	if err != nil {
		return nil, fmt.Errorf("infra: falha ao buscar detalhes da assinatura na stripe: %w", err)
	}

	// Mapeia os dados da Stripe para a nossa struct de domínio.
	details := &domain.SubscriptionDetails{
		ID:                 s.ID,
		CurrentPeriodStart: time.Unix(s.CurrentPeriodStart, 0),
		CurrentPeriodEnd:   time.Unix(s.CurrentPeriodEnd, 0),
	}

	return details, nil
}
