package infrastructure

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/luiszkm/wedding_backend/internal/event/domain"
)

func (r *PostgresEventoRepository) FindByUserID(ctx context.Context, userID uuid.UUID) ([]*domain.Evento, error) {
	sql := `SELECT id, id_usuario, nome, data, tipo, url_slug FROM eventos WHERE id_usuario = $1 ORDER BY data DESC`
	rows, err := r.db.Query(ctx, sql, userID)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar eventos por usuário: %w", err)
	}
	defer rows.Close()

	var eventos []*domain.Evento
	for rows.Next() {
		var id, idUsuario uuid.UUID
		var nome, tipo, urlSlug string
		var data time.Time

		err := rows.Scan(&id, &idUsuario, &nome, &data, &tipo, &urlSlug)
		if err != nil {
			return nil, fmt.Errorf("erro ao scanear evento: %w", err)
		}

		eventos = append(eventos, domain.HydrateEvento(id, idUsuario, nome, data, domain.TipoEvento(tipo), urlSlug))
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("erro ao iterar sobre eventos: %w", err)
	}

	return eventos, nil
}