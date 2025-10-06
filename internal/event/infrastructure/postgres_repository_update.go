// file: internal/event/infrastructure/postgres_repository_update.go
package infrastructure

import (
	"context"
	"fmt"

	"github.com/luiszkm/wedding_backend/internal/event/domain"
)

// Update atualiza um evento existente
func (r *PostgresEventoRepository) Update(ctx context.Context, evento *domain.Evento) error {
	sql := `
        UPDATE eventos
        SET nome = $2, data = $3, tipo = $4, url_slug = $5
        WHERE id = $1
    `
	result, err := r.db.Exec(ctx, sql,
		evento.ID(),
		evento.Nome(),
		evento.Data(),
		evento.Tipo(),
		evento.UrlSlug(),
	)
	if err != nil {
		return fmt.Errorf("falha ao atualizar evento: %w", err)
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return domain.ErrEventoNaoEncontrado
	}

	return nil
}