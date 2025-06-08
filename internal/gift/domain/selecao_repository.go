// file: internal/gift/domain/selecao_repository.go
package domain

import (
	"context"

	"github.com/google/uuid"
)

type SelecaoRepository interface {
	// SalvarSelecao deve ser uma operação transacional.
	SalvarSelecao(ctx context.Context, chaveDeAcesso string, idsDosPresentes []uuid.UUID) (*Selecao, error)
}
