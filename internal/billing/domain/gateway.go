// file: internal/billing/domain/gateway.go
package domain

import (
	"context"
	"time"
)

type SubscriptionDetails struct {
	ID                 string
	CurrentPeriodStart time.Time
	CurrentPeriodEnd   time.Time
}

// PaymentGateway define a interface para um serviço externo de pagamentos.
// Note que os métodos usam nossos próprios objetos de domínio, não objetos da Stripe.
type PaymentGateway interface {
	CriarSessaoCheckout(ctx context.Context, assinatura *Assinatura, plano *Plano) (checkoutURL string, err error)
	GetSubscriptionDetails(ctx context.Context, stripeSubscriptionID string) (*SubscriptionDetails, error)
}
