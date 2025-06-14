// file: internal/event/infrastructure/postgres_repository.go
package infrastructure

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/luiszkm/wedding_backend/internal/event/domain"
)

type PostgresEventoRepository struct {
	db *pgxpool.Pool
}

func NewPostgresEventoRepository(db *pgxpool.Pool) domain.EventoRepository {
	return &PostgresEventoRepository{db: db}
}

func (r *PostgresEventoRepository) Save(ctx context.Context, evento *domain.Evento) error {
	sql := `
        INSERT INTO eventos (id, id_usuario, nome, data, tipo, url_slug) 
        VALUES ($1, $2, $3, $4, $5, $6)
    `
	_, err := r.db.Exec(ctx, sql,
		evento.ID(),
		evento.IDUsuario(),
		evento.Nome(),
		evento.Data(),
		evento.Tipo(),
		evento.UrlSlug(),
	)
	if err != nil {
		// Aqui poderíamos verificar erros de constraint, como slug duplicado
		return fmt.Errorf("falha ao salvar evento: %w", err)
	}
	return nil
}

func (r *PostgresEventoRepository) FindBySlug(ctx context.Context, slug string) (*domain.Evento, error) {
	sql := `SELECT id, id_usuario, nome, data, tipo, url_slug FROM eventos WHERE url_slug = $1`
	row := r.db.QueryRow(ctx, sql, slug)

	var id, idUsuario uuid.UUID
	var nome, tipo, urlSlug string
	var data time.Time
	err := row.Scan(&id, &idUsuario, &nome, &data, &tipo, &urlSlug)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrEventoNaoEncontrado
		}
		return nil, fmt.Errorf("falha ao buscar evento por slug: %w", err)
	}

	// Usariamos uma função Hydrate aqui, similar aos outros contextos
	return domain.HydrateEvento(id, idUsuario, nome, data, domain.TipoEvento(tipo), urlSlug), nil

}
func (r *PostgresEventoRepository) FindByID(ctx context.Context, userID, eventID uuid.UUID) (*domain.Evento, error) {
	sql := `SELECT id, id_usuario, nome, data, tipo, url_slug FROM eventos WHERE id = $1 AND id_usuario = $2`
	row := r.db.QueryRow(ctx, sql, eventID, userID)

	var id, idUsuario uuid.UUID
	var nome, tipo, urlSlug string
	var data time.Time

	err := row.Scan(&id, &idUsuario, &nome, &data, &tipo, &urlSlug)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrEventoNaoEncontrado
		}
		return nil, fmt.Errorf("falha ao buscar evento por id: %w", err)
	}
	// Usariamos uma função Hydrate aqui, similar aos outros contextos
	return domain.HydrateEvento(id, idUsuario, nome, data, domain.TipoEvento(tipo), urlSlug), nil
}
