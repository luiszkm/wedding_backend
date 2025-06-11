// file: internal/billing/domain/gateway.go
package domain

import "context"

// PaymentGateway define a interface para um serviço externo de pagamentos.
// Note que os métodos usam nossos próprios objetos de domínio, não objetos da Stripe.
type PaymentGateway interface {
	CriarSessaoCheckout(ctx context.Context, assinatura *Assinatura, plano *Plano) (checkoutURL string, err error)
}
