// file: internal/guest/domain/repository.go
package domain

import (
	"context"
)

// GroupRepository é a interface de persistência para o agregado GrupoDeConvidados.
// O repositório lida com o agregado como um todo.
type GroupRepository interface {
	Save(ctx context.Context, group *GrupoDeConvidados) error
	FindByAccessKey(ctx context.Context, accessKey string) (*GrupoDeConvidados, error)

	// Outros métodos como FindByID(id uuid.UUID) ou FindByAccessKey(key string) viriam aqui.
}
