package domain

import (
	"context"

	"github.com/google/uuid"
)

type ComunicadoRepository interface {
	Criar(ctx context.Context, comunicado *Comunicado) error
	BuscarPorEvento(ctx context.Context, idEvento uuid.UUID) ([]*Comunicado, error)
	BuscarPorID(ctx context.Context, id uuid.UUID) (*Comunicado, error)
	Editar(ctx context.Context, comunicado *Comunicado) error
	Deletar(ctx context.Context, id uuid.UUID) error
}
