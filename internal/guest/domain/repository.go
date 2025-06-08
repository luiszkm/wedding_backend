// file: internal/guest/domain/repository.go
package domain

import (
	"context"

	"github.com/google/uuid"
)

type GroupRepository interface {
	Save(ctx context.Context, group *GrupoDeConvidados) error
	FindByAccessKey(ctx context.Context, accessKey string) (*GrupoDeConvidados, error)
	Update(ctx context.Context, userID uuid.UUID, group *GrupoDeConvidados) error        // <-- userID adicionado
	FindByID(ctx context.Context, userID, groupID uuid.UUID) (*GrupoDeConvidados, error) // <-- userID adicionado
	UpdateRSVP(ctx context.Context, group *GrupoDeConvidados) error
}
