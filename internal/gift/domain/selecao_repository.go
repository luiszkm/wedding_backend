// file: internal/gift/domain/selecao_repository.go
package domain

import (
	"context"

	"github.com/google/uuid"
)

type SelecaoRepository interface {
	// SalvarSelecao deve ser uma operação transacional.
	// O mapa quantidades mapeia ID do presente -> quantidade desejada
	SalvarSelecao(ctx context.Context, chaveDeAcesso string, quantidades map[uuid.UUID]int) (*Selecao, error)
}
