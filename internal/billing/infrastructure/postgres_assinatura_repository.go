// file: internal/billing/infrastructure/postgres_assinatura_repository.go
package infrastructure

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/luiszkm/wedding_backend/internal/billing/domain" // Ajuste o path se necessário
)

type PostgresAssinaturaRepository struct {
	db *pgxpool.Pool
}

func NewPostgresAssinaturaRepository(db *pgxpool.Pool) domain.AssinaturaRepository {
	return &PostgresAssinaturaRepository{db: db}
}

func (r *PostgresAssinaturaRepository) Save(ctx context.Context, assinatura *domain.Assinatura) error {
	sql := `INSERT INTO assinaturas (id, id_usuario, id_plano, status) VALUES ($1, $2, $3, $4)`
	_, err := r.db.Exec(ctx, sql,
		assinatura.ID(),
		assinatura.IDUsuario(),
		assinatura.IDPlano(),
		assinatura.Status(),
	)
	if err != nil {
		return fmt.Errorf("falha ao salvar assinatura: %w", err)
	}
	return nil
}
