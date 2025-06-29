// file: internal/gallery/domain/repository.go
package domain

import (
	"context"

	"github.com/google/uuid"
)

type FotoRepository interface {
	// Salva múltiplas fotos e seus rótulos em uma única transação.
	SalvarMultiplas(ctx context.Context, fotos []*Foto) error
	ListarPublicasPorCasamento(ctx context.Context, casamentoID uuid.UUID, filtroRotulo Rotulo) ([]*Foto, error)
	FindByID(ctx context.Context, userID, fotoID uuid.UUID) (*Foto, error)
	Update(ctx context.Context, userID uuid.UUID, foto *Foto) error
	Delete(ctx context.Context, userID uuid.UUID, fotoID uuid.UUID) error
}
