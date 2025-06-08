// file: internal/gift/domain/repository.go
package domain

import (
	"context"

	"github.com/google/uuid"
)

type PresenteRepository interface {
	Save(ctx context.Context, presente *Presente) error
	ListarDisponiveisPorCasamento(ctx context.Context, casamentoID uuid.UUID) ([]*Presente, error)
}
