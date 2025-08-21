package infrastructure

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/luiszkm/wedding_backend/internal/event/domain"
)

func (r *PostgresEventoRepository) FindByUserID(ctx context.Context, userID uuid.UUID) ([]*domain.Evento, error) {
	sql := `SELECT id, id_usuario, nome, data, tipo, url_slug, id_template, id_template_arquivo, paleta_cores FROM eventos WHERE id_usuario = $1 ORDER BY data DESC`
	rows, err := r.db.Query(ctx, sql, userID)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar eventos por usuário: %w", err)
	}
	defer rows.Close()

	var eventos []*domain.Evento
	for rows.Next() {
		var id, idUsuario uuid.UUID
		var nome, tipo, urlSlug, idTemplate string
		var idTemplateArquivo *string
		var data time.Time
		var paletaJSON []byte

		err := rows.Scan(&id, &idUsuario, &nome, &data, &tipo, &urlSlug, &idTemplate, &idTemplateArquivo, &paletaJSON)
		if err != nil {
			return nil, fmt.Errorf("erro ao scanear evento: %w", err)
		}

		// Converter JSON para PaletaCores
		var paletaCores domain.PaletaCores
		if len(paletaJSON) > 0 {
			if err := json.Unmarshal(paletaJSON, &paletaCores); err != nil {
				return nil, fmt.Errorf("erro ao converter paleta de cores: %w", err)
			}
		}

		eventos = append(eventos, domain.HydrateEvento(id, idUsuario, nome, data, domain.TipoEvento(tipo), urlSlug, idTemplate, idTemplateArquivo, paletaCores))
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("erro ao iterar sobre eventos: %w", err)
	}

	return eventos, nil
}