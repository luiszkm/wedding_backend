// file: internal/iam/domain/repository.go
package domain

import "context"

type UsuarioRepository interface {
	Save(ctx context.Context, usuario *Usuario) error
	FindByEmail(ctx context.Context, email string) (*Usuario, error)
}
