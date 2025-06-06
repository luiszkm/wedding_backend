// file: internal/guest/domain/repository.go
package domain

import (
	"context"

	"github.com/google/uuid"
)

// GroupRepository é a interface de persistência para o agregado GrupoDeConvidados.
// O repositório lida com o agregado como um todo.
type GroupRepository interface {
	Save(ctx context.Context, group *GrupoDeConvidados) error
	FindByAccessKey(ctx context.Context, accessKey string) (*GrupoDeConvidados, error)
	Update(ctx context.Context, group *GrupoDeConvidados) error
	FindByID(ctx context.Context, id uuid.UUID) (*GrupoDeConvidados, error)
	// Outros métodos como FindByID(id uuid.UUID) ou FindByAccessKey(key string) viriam aqui.
}
