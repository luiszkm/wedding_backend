// file: internal/event/domain/repository.go
package domain

import (
	"context"

	"github.com/google/uuid"
)

type EventoRepository interface {
	Save(ctx context.Context, evento *Evento) error
	FindBySlug(ctx context.Context, slug string) (*Evento, error)
	FindByID(ctx context.Context, userID, eventID uuid.UUID) (*Evento, error)
}
