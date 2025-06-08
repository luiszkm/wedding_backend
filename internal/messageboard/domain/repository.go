// file: internal/messageboard/domain/repository.go
package domain

import (
	"context"

	"github.com/google/uuid"
)

type RecadoRepository interface {
	Save(ctx context.Context, recado *Recado) error
	ListarPorEvento(ctx context.Context, casamentoID uuid.UUID) ([]*Recado, error)
	FindByID(ctx context.Context, userID, recadoID uuid.UUID) (*Recado, error)
	Update(ctx context.Context, userID uuid.UUID, recado *Recado) error
	ListarAprovadosPorCasamento(ctx context.Context, casamentoID uuid.UUID) ([]*Recado, error)
}
