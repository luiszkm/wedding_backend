// file: internal/gift/domain/repository.go
package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// PresenteComSelecao representa um presente com informações de confirmação
// Presentes fracionados podem aparecer múltiplas vezes (uma para cada palavra mágica)
type PresenteComSelecao struct {
	Presente        *Presente
	ChaveDeAcesso   *string    // null se não confirmado
	QuantidadeCotas int        // quantas cotas essa palavra mágica pegou (1 para integrais)
	DataSelecao     *time.Time // null se não confirmado
}

type PresenteRepository interface {
	Save(ctx context.Context, presente *Presente) error
	ListarDisponiveisPorCasamento(ctx context.Context, casamentoID uuid.UUID) ([]*Presente, error)
	ListarTodosPorEvento(ctx context.Context, eventoID uuid.UUID) ([]*PresenteComSelecao, error)
	FindByIDs(ctx context.Context, presenteIDs []uuid.UUID) ([]*Presente, error)
	Update(ctx context.Context, presente *Presente) error
	SaveWithCotas(ctx context.Context, presente *Presente) error
	UpdateCotasStatus(ctx context.Context, cotaIDs []uuid.UUID, status string, selecaoID *uuid.UUID) error
}
