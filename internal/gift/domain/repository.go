// file: internal/gift/domain/repository.go
package domain

import (
	"context"

	"github.com/google/uuid"
)

type PresenteRepository interface {
	Save(ctx context.Context, presente *Presente) error
	ListarDisponiveisPorCasamento(ctx context.Context, casamentoID uuid.UUID) ([]*Presente, error)
	FindByIDs(ctx context.Context, presenteIDs []uuid.UUID) ([]*Presente, error)
	Update(ctx context.Context, presente *Presente) error
	SaveWithCotas(ctx context.Context, presente *Presente) error
	UpdateCotasStatus(ctx context.Context, cotaIDs []uuid.UUID, status string, selecaoID *uuid.UUID) error
}
