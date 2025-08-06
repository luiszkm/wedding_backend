// file: internal/event/infrastructure/postgres_repository_update.go
package infrastructure

import (
	"context"
	"fmt"

	"github.com/luiszkm/wedding_backend/internal/event/domain"
)

// Update atualiza um evento existente
func (r *PostgresEventoRepository) Update(ctx context.Context, evento *domain.Evento) error {
	paletaJSON, err := evento.PaletaCoresJSON()
	if err != nil {
		return fmt.Errorf("erro ao converter paleta de cores: %w", err)
	}
	
	sql := `
        UPDATE eventos 
        SET nome = $2, data = $3, tipo = $4, url_slug = $5, id_template = $6, id_template_arquivo = $7, paleta_cores = $8
        WHERE id = $1
    `
	result, err := r.db.Exec(ctx, sql,
		evento.ID(),
		evento.Nome(),
		evento.Data(),
		evento.Tipo(),
		evento.UrlSlug(),
		evento.IDTemplate(),
		evento.IDTemplateArquivo(),
		paletaJSON,
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