// file: internal/billing/domain/repository.go
package domain

import (
	"context"

	"github.com/google/uuid"
)

type PlanoRepository interface {
	ListAll(ctx context.Context) ([]*Plano, error)
	FindByID(ctx context.Context, id uuid.UUID) (*Plano, error) // <-- NOVO

}

type AssinaturaRepository interface {
	Save(ctx context.Context, assinatura *Assinatura) error
	FindByID(ctx context.Context, id uuid.UUID) (*Assinatura, error) // <-- NOVO
	FindByStripeSubscriptionID(ctx context.Context, stripeID string) (*Assinatura, error)

	Update(ctx context.Context, assinatura *Assinatura) error
}
